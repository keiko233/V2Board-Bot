package controller

import (
	"fmt"
	"log"
	"strings"

	"github.com/keiko233/V2Board-Bot/dao"
	"github.com/keiko233/V2Board-Bot/model"
	"github.com/keiko233/V2Board-Bot/service"
	tb "gopkg.in/tucnak/telebot.v2"
)

func BindCmdCtr(m *tb.Message) {
	if m.Chat.ID < 0 {
		model.Bot.Reply(m, "请私聊我命令 /bind <订阅地址>")
		return
	}
	user, notfound, err := dao.GetUserByTelegramID(nil, m.Sender.ID)
	if err != nil {
		log.Printf("QueryUser tgid = %d error, %s\n", m.Sender.ID, err)
		model.Bot.Reply(m, "👀 获取失败，请稍后重试或联系管理员")
		return
	}
	if !notfound {
		model.Bot.Send(m.Chat, fmt.Sprintf("✅ 当前绑定账户: %s\n若需要修改绑定,需要解绑当前账户。", user.Email))
		return
	}

	format := strings.Index(m.Text, "token=")
	if format <= 0 {
		model.Bot.Send(m.Chat, "👀 ️账户绑定格式: /bind <订阅地址>")
		return
	}

	user, err = service.BindUser(m.Text[format:][6:38], m.Sender.ID)
	if err != nil {
		log.Printf("Bind User token=%s and tgid=%d err %s\n", m.Text[6:38], m.Sender.ID, err)
		model.Bot.Send(m.Chat, "❌ 订阅无效,请前往官网复制最新订阅地址!")
		return
	}
	model.Bot.Send(m.Chat, fmt.Sprintf("✅ 账户绑定成功: %s", user.Email))
}

func UnbindCmdCtr(m *tb.Message) {
	if m.Chat.ID < 0 {
		model.Bot.Reply(m, "请私聊我解绑哦~")
		return
	}
	notfound, err := service.UnbindUser(m.Sender.ID)
	if err != nil {
		log.Printf("unbind user by tgid=%d error %s\n", m.Sender.ID, err)
		model.Bot.Reply(m, "❌ 账户解绑失败,请稍后再试")
		return
	}

	if notfound {
		model.Bot.Reply(m, "👀 当前未绑定账户")
		return
	}

	model.Bot.Reply(m, "✅ 账户解绑成功")
}
