package controller

import (
	"fmt"
	"log"
	"strconv"

	"github.com/keiko233/V2Board-Bot/common"
	"github.com/keiko233/V2Board-Bot/dao"
	"github.com/keiko233/V2Board-Bot/lib/tgbot"
	"github.com/keiko233/V2Board-Bot/service"
	"github.com/keiko233/V2Board-Bot/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

func Checkin(ctx *tgbot.Context) error {
	id, err := strconv.ParseInt(ctx.Callback.Data, 10, 64)
	if err != nil {
		return err
	}
	if id != ctx.Callback.Sender.ID {
		return ctx.AnswerCallback("再瞎点其他人的按钮, 要报警了!!!")
	}
	user, notfound, err := dao.GetUserByTelegramID(nil, id)
	if err != nil {
		log.Printf("QueryUser tgid = %d error, %s\n", id, err)
		return err
	}

	if notfound {
		return common.ErrNotBindUser
	}

	if user.PlanId <= 0 {
		return common.ErrNotFoundPlan
	}
	todayNotCheckin, err := service.CheckinTime(id)
	if err != nil {
		log.Printf("CheckinTime err %s\n", err)
		return err
	}
	if !todayNotCheckin {
		return common.ErrAlreadyCheckin
	}

	// b := rand.RandInt(model.Config.Bot.MaxByte, model.Config.Bot.MinByte)
	f, err := service.GetFortune(id)
	if err != nil {
		return err
	}
	b, err := service.GetTraffer(f, id)
	l, err := service.CheckinUser(id, b, f)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("✅ 签到成功\n今天的运势是 %s \n本次签到获得 %s 流量\n签到次数每日0点刷新，明天再来哦！", f, utils.ByteSize(l.CheckinTraffic))
	return ctx.Edit(msg)
}

func Fortune(ctx *tgbot.Context) error {
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

	f, err := service.GetFortune(ctx.Message.Sender.ID)
	if err != nil {
		return err
	}
	r := make([][]tb.InlineButton, 0)
	r = append(r, []tb.InlineButton{
		{
			Unique: "checkin",
			Text:   "领取流量",
			Data:   strconv.Itoa(int(ctx.Message.Sender.ID)),
		},
		{
			Unique: "passpool",
			Text:   "还愿",
			Data:   strconv.Itoa(int(ctx.Message.Sender.ID)),
		},
	})

	return ctx.Reply(fmt.Sprintf("今天的运势是%s哦~\n\n要领取流量还是还愿呢?", f),
		&tb.ReplyMarkup{
			InlineKeyboard: r,
		})
}

func PassPool(ctx *tgbot.Context) error {
	id, err := strconv.ParseInt(ctx.Callback.Data, 10, 64)
	if err != nil {
		return err
	}
	if id != ctx.Callback.Sender.ID {
		return ctx.AnswerCallback("再瞎点其他人的按钮, 要报警了!!!")
	}
	user, notfound, err := dao.GetUserByTelegramID(nil, id)
	if err != nil {
		log.Printf("QueryUser tgid = %d error, %s\n", id, err)
		return err
	}

	if notfound {
		return common.ErrNotBindUser
	}

	if user.PlanId <= 0 {
		return common.ErrNotFoundPlan
	}
	todayNotCheckin, err := service.CheckinTime(id)
	if err != nil {
		log.Printf("CheckinTime err %s\n", err)
		return err
	}
	if !todayNotCheckin {
		return common.ErrAlreadyCheckin
	}

	f, err := service.GetFortune(id)
	if err != nil {
		return err
	}
	b, err := service.GetTraffer(f, id)
	if err != nil {
		return err
	}

	n, yes, err := service.PassPool(b, f, id)
	if err != nil {
		return err
	}

	if n != 0 {
		l, err := service.CheckinUser(id, n, f)
		if err != nil {
			return err
		}
		msg := fmt.Sprintf("还愿池爆满啦!\n池子里的 %s 流量现在都是你的啦！！", utils.ByteSize(l.CheckinTraffic))
		return ctx.Edit(msg)
	}
	var s string
	if yes {
		s = "还愿"
	} else {
		s = "还怨"
	}
	_, err = service.CheckinUser(id, 0, f)
	if err != nil {
		return err
	}
	return ctx.Edit(fmt.Sprintf("今天的运势是 %s ,成功%s了%s流量哦~\n\n", f, s, utils.ByteSize(b*1024*1024)))
}
