package controller

import (
	"fmt"
	"log"

	"github.com/keiko233/V2Board-Bot/common"
	"github.com/keiko233/V2Board-Bot/dao"
	"github.com/keiko233/V2Board-Bot/lib/tgbot"
	"github.com/keiko233/V2Board-Bot/service"
	"github.com/keiko233/V2Board-Bot/utils"
)

func Checkin(ctx *tgbot.Context) error {
	user, notfound, err := dao.GetUserByTelegramID(nil, ctx.Message.Sender.ID)
	if err != nil {
		log.Printf("QueryUser tgid = %d error, %s\n", ctx.Message.Sender.ID, err)
		return err
	}

	if notfound {
		return common.ErrNotBindUser
	}

	if user.PlanId <= 0 {
		return common.ErrNotFoundPlan
	}
	todayNotCheckin, err := service.CheckinTime(ctx.Message.Sender.ID)
	if err != nil {
		log.Printf("CheckinTime err %s\n", err)
		return err
	}
	if !todayNotCheckin {
		return common.ErrAlreadyCheckin
	}

	l, err := service.CheckinUser(ctx.Message.Sender.ID)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("✅ 签到成功\n本次签到获得 %s 流量\n签到次数每日0点刷新，明天再来哦！", utils.ByteSize(l.CheckinTraffic))
	return ctx.Reply(msg)
}
