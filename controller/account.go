package controller

import (
	"fmt"
	"log"

	"github.com/keiko233/V2Board-Bot/common"
	"github.com/keiko233/V2Board-Bot/dao"
	"github.com/keiko233/V2Board-Bot/lib/tgbot"
	"github.com/keiko233/V2Board-Bot/utils"
)

func Account(ctx *tgbot.Context) error {
	m := ctx.Message
	user, notfound, err := dao.GetUserByTelegramID(nil, m.Sender.ID)

	if err != nil {
		log.Printf("QueryUser tgid = %d error, %s\n", m.Sender.ID, err)
		return err
	}

	if notfound {
		return common.ErrNotBindUser
	}

	p, notfound, err := dao.GetPlanByID(nil, int(user.PlanId))
	if err != nil {
		log.Printf("QueryPlan id = %d error, %s\n", user.PlanId, err)
		return err
	}

	if notfound {
		return common.ErrInvalidPlan
	}

	Email := user.Email
	if m.Chat.ID <= 0 {
		Email = "不可以偷窥哦~~"
	}

	CreatedAt := utils.UnixToStr(user.CreatedAt)
	Balance := user.Balance / 100
	CommissionBalance := user.CommissionBalance / 100
	PlanName := p.Name
	ExpiredAt := utils.UnixToStr(user.ExpiredAt)
	TransferEnable := utils.ByteSize(user.TransferEnable)
	U := utils.ByteSize(user.U)
	D := utils.ByteSize(user.D)
	S := utils.ByteSize(user.TransferEnable - (user.U + user.D))

	if user.PlanId <= 0 {
		msg := fmt.Sprintf("账户信息概况:\n\n当前绑定账户: %s\n注册时间: %s\n账户余额: %d元\n佣金余额: %d元\n\n当前订阅: 当前暂无订阅计划", Email, CreatedAt, Balance, CommissionBalance)
		return ctx.Reply(msg)
	}

	msg := fmt.Sprintf("账户信息概况:\n\n当前绑定账户: %s\n注册时间: %s\n账户余额: %d元\n佣金余额: %d元\n\n当前订阅: %s\n到期时间: %s\n订阅流量: %s\n已用上行: %s\n已用下行: %s\n剩余可用: %s", Email, CreatedAt, Balance, CommissionBalance, PlanName, ExpiredAt, TransferEnable, U, D, S)
	return ctx.Reply(msg)

}
