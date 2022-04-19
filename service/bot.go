package service

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var Bot *tb.Bot

func Start() {
	var err error
	Bot, err = tb.NewBot(tb.Settings{
		URL:    "https://api.telegram.org",
		Token:  c.Bot.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("Bot 启动失败啦...... \n当前Token [ %s ] \n错误信息:  %s", c.Bot.Token, err)
	}

	setHandle()
	Bot.Start()
}

func setHandle() {
	Bot.Handle("/start", startCmdCtr)
	Bot.Handle("/help", startCmdCtr)
	Bot.Handle("/checkin", checkinCmdCtr)
	Bot.Handle("/account", accountCmdCtr)
	Bot.Handle("/bind", bindCmdCtr)
	Bot.Handle("/unbind", unbindCmdCtr)
}

func startCmdCtr(m *tb.Message) {
	menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	CheckinBtn := menu.Text("👀 每日签到")
	AccountBtn := menu.Text("🚥‍ 账户信息")
	BindBtn := menu.Text("😋 绑定账户")
	UnbindBtn := menu.Text("🤔 解绑账户")

	menu.Reply(
		menu.Row(CheckinBtn, AccountBtn),
		menu.Row(BindBtn, UnbindBtn),
	)

	Bot.Handle(&CheckinBtn, checkinCmdCtr)
	Bot.Handle(&AccountBtn, accountCmdCtr)
	Bot.Handle(&BindBtn, bindCmdCtr)
	Bot.Handle(&UnbindBtn, unbindCmdCtr)

	msg := fmt.Sprintf("%s\n为你提供以下服务:\n\n每日签到 /checkin\n账户信息 /account\n绑定账户 /bind\n解绑账户 /unbind", c.Bot.Name)
	_, _ = Bot.Reply(m, msg, menu)
}

func checkinCmdCtr(m *tb.Message) {
	user := QueryUser(m.Sender.ID)
	if user.Id <= 0 {
		msg := "👀 当前未绑定账户\n请发送 /bind <订阅地址> 绑定账户"
		if _, err := Bot.Reply(m, msg); err != nil {
			log.Printf("未绑定账户 Bot Reply %s\n", err)
		}
		return
	}
	if user.PlanId <= 0 {
		msg := "👀 当前暂无订阅计划,该功能需要订阅后使用～"
		if _, err := Bot.Reply(m, msg); err != nil {
			log.Printf("无订阅计划 Bot Reply %s\n", err)
		}
		return
	}

	cc := CheckinTime(m.Sender.ID)
	if cc == false {
		msg := fmt.Sprintf("✅ 今天已经签到过啦！明天再来哦～")
		if _, err := Bot.Reply(m, msg); err != nil {
			log.Printf("已经签到过 Bot Reply %s\n", err)
		}
		return
	}

	uu := checkinUser(m.Sender.ID)

	msg := fmt.Sprintf("✅ 签到成功\n本次签到获得 %s 流量\n下次签到时间: %s", ByteSize(uu.CheckinTraffic), UnixToStr(uu.NextAt))
	if _, err := Bot.Reply(m, msg); err != nil {
		log.Printf("签到成功Bot Reply %s\n", err)
	}
}

func accountCmdCtr(m *tb.Message) {
	user := QueryUser(m.Sender.ID)
	if user.Id <= 0 {
		msg := "👀 当前未绑定账户\n请私聊发送 /bind <订阅地址> 绑定账户"
		if _, err := Bot.Reply(m, msg); err != nil {
			log.Printf("Bot Reply %s\n", err)
		}
		return
	}
	p := QueryPlan(int(user.PlanId))
	Email := user.Email
	CreatedAt := UnixToStr(user.CreatedAt)
	Balance := user.Balance / 100
	CommissionBalance := user.CommissionBalance / 100
	PlanName := p.Name
	ExpiredAt := UnixToStr(user.ExpiredAt)
	TransferEnable := ByteSize(user.TransferEnable)
	U := ByteSize(user.U)
	D := ByteSize(user.D)
	S := ByteSize(user.TransferEnable - (user.U + user.D))
	if user.PlanId <= 0 {
		msg := fmt.Sprintf("账户信息概况:\n\n当前绑定账户: %s\n注册时间: %s\n账户余额: %d元\n佣金余额: %d元\n\n当前订阅: 当前暂无订阅计划", Email, CreatedAt, Balance, CommissionBalance)
		if _, err := Bot.Reply(m, msg); err != nil {
			log.Printf("Bot Reply %s\n", err)
		}
		return
	}

	msg := fmt.Sprintf("账户信息概况:\n\n当前绑定账户: %s\n注册时间: %s\n账户余额: %d元\n佣金余额: %d元\n\n当前订阅: %s\n到期时间: %s\n订阅流量: %s\n已用上行: %s\n已用下行: %s\n剩余可用: %s", Email, CreatedAt, Balance, CommissionBalance, PlanName, ExpiredAt, TransferEnable, U, D, S)
	if _, err := Bot.Reply(m, msg); err != nil {
		log.Printf("Bot Reply %s\n", err)
	}

}

func bindCmdCtr(m *tb.Message) {
	if m.Chat.ID < 0 {
		// _, _ = Bot.Send(m.Chat, "请私聊我命令 /bind <订阅地址>")
		Bot.Reply(m, "请私聊我命令 /bind <订阅地址>")
		return
	}
	user := QueryUser(m.Sender.ID)
	if user.Id > 0 {
		_, _ = Bot.Send(m.Chat, fmt.Sprintf("✅ 当前绑定账户: %s\n若需要修改绑定,需要解绑当前账户。", user.Email))
		return
	}

	format := strings.Index(m.Text, "token=")
	if format <= 0 {
		_, _ = Bot.Send(m.Chat, "👀 ️账户绑定格式: /bind <订阅地址>")
		return
	}

	b := BindUser(m.Text[format:], m.Sender.ID)
	if b.Id <= 0 {
		_, _ = Bot.Send(m.Chat, "❌ 订阅无效,请前往官网复制最新订阅地址!")
		return
	}

	if b.TelegramId != uint(m.Sender.ID) {
		_, _ = Bot.Send(m.Chat, "❌ 账户绑定失败,请稍后再试")
	}
	_, _ = Bot.Send(m.Chat, fmt.Sprintf("✅ 账户绑定成功: %s", b.Email))
}

func unbindCmdCtr(m *tb.Message) {
	user := unbindUser(m.Sender.ID)
	if user.Id <= 0 {
		_, _ = Bot.Reply(m, "👀 当前未绑定账户")
		return
	}
	if user.TelegramId > 0 {
		_, _ = Bot.Reply(m, "❌ 账户解绑失败,请稍后再试")
		return
	}
	_, _ = Bot.Reply(m, "✅ 账户解绑成功")
}

func UnixToStr(unix int64) string {
	u := time.Unix(unix, 0).Format("2006-01-02 15:04:05")
	return u
}

func ByteSize(size int64) string {
	sizeFloat := float64(size)
	oldSize := sizeFloat
	var n float64 = 0
	for math.Abs(sizeFloat) >= 1024 {
		sizeFloat = sizeFloat / 1024
		n++
	}

	var k string
	if n == 0 {
		k = "B"
	} else if n == 1 {
		k = "KB"
	} else if n == 2 {
		k = "MB"
	} else if n == 3 {
		k = "GB"
	} else if n == 4 {
		k = "TB"
	}

	ns := oldSize / math.Pow(1024, n)

	return fmt.Sprintf("%.2f%s", ns, k)
}
