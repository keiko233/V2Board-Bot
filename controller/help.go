package controller

import (
	"fmt"

	"github.com/keiko233/V2Board-Bot/lib/tgbot"
	"github.com/keiko233/V2Board-Bot/model"
	tb "gopkg.in/tucnak/telebot.v2"
)

func Help(ctx *tgbot.Context) error {
	menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	CheckinBtn := menu.Text(model.MenuCheckinBtn)
	AccountBtn := menu.Text(model.MenuAccountBtn)
	BindBtn := menu.Text(model.MenuBindBtn)
	UnbindBtn := menu.Text(model.MenuUnbindBtn)
	historyBtn := menu.Text(model.MenuhistoryBtn)
	reportBtn := menu.Text(model.MenuReportBtn)

	menuList := make([]tb.Row, 0)
	// 群聊发起, 不展示解绑和绑定
	if ctx.Message.Chat.ID < 0 {
		menuList = append(menuList, menu.Row(CheckinBtn, AccountBtn), menu.Row(historyBtn, reportBtn))
	} else {
		menuList = append(menuList, menu.Row(CheckinBtn, AccountBtn), menu.Row(BindBtn, UnbindBtn), menu.Row(historyBtn, reportBtn))
	}

	menu.Reply(menuList...)

	msg := fmt.Sprintf("%s\n为你提供以下服务:\n\n每日签到 /checkin\n账户信息 /account\n绑定账户 /bind\n解绑账户 /unbind\n签到历史 /history\n签到统计 /report\n请注意, 绑定账号和解绑账号需要私聊我哦~", model.Config.Bot.Name)
	return ctx.Reply(msg, menu)
}
