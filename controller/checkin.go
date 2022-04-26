package controller

import (
	"fmt"
	"log"

	"github.com/keiko233/V2Board-Bot/dao"
	"github.com/keiko233/V2Board-Bot/model"
	"github.com/keiko233/V2Board-Bot/service"
	"github.com/keiko233/V2Board-Bot/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

func CheckinCmdCtr(m *tb.Message) {
	user, notfound, err := dao.GetUserByTelegramID(nil, m.Sender.ID)

	if err != nil {
		log.Printf("QueryUser tgid = %d error, %s\n", m.Sender.ID, err)
		model.Bot.Reply(m, "👀 获取失败，请稍后重试或联系管理员")
		return
	}

	if notfound {
		model.Bot.Reply(m, "👀 当前未绑定账户\n请私聊发送 /bind <订阅地址> 绑定账户")
		return
	}

	if user.PlanId <= 0 {
		msg := "👀 当前暂无订阅计划,该功能需要订阅后使用～"
		if _, err := model.Bot.Reply(m, msg); err != nil {
			log.Printf("无订阅计划 Bot Reply %s\n", err)
		}
		return
	}
	todayNotCheckin, err := service.CheckinTime(m.Sender.ID)
	if err != nil {
		log.Printf("CheckinTime err %s\n", err)
		model.Bot.Reply(m, "👀 获取失败，请稍后重试或联系管理员")
		return
	}
	if !todayNotCheckin {
		msg := fmt.Sprintf("✅ 今天已经签到过啦！明天再来哦～")
		if _, err := model.Bot.Reply(m, msg); err != nil {
			log.Printf("已经签到过 Bot Reply %s\n", err)
		}
		return
	}

	l, err := service.CheckinUser(m.Sender.ID)
	if err != nil {
		log.Printf("操作失败 %s\n", err)
		if _, err := model.Bot.Reply(m, "操作失败！请联系管理员！"); err != nil {
			log.Printf("操作失败 Bot Reply %s\n", err)
		}
		return
	}

	msg := fmt.Sprintf("✅ 签到成功\n本次签到获得 %s 流量\n签到次数每日0点刷新，明天再来哦！", utils.ByteSize(l.CheckinTraffic))
	if _, err := model.Bot.Reply(m, msg); err != nil {
		log.Printf("签到成功 Bot Reply %s\n", err)
	}
}
