package controller

import (
	"fmt"
	"github.com/keiko233/V2Board-Bot/model"
	tb "gopkg.in/tucnak/telebot.v2"
)

func StartCmdCtr(m *tb.Message) {
	menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	CheckinBtn := menu.Text("👀 每日签到")
	AccountBtn := menu.Text("🚥‍ 账户信息")
	BindBtn := menu.Text("😋 绑定账户")
	UnbindBtn := menu.Text("🤔 解绑账户")
	historyBtn := menu.Text("📅 签到历史")

	menuList := make([]tb.Row, 0)
	// 群聊发起, 不展示解绑和绑定
	if m.Chat.ID < 0 {
		menuList = append(menuList, menu.Row(CheckinBtn, AccountBtn), menu.Row(historyBtn))
	} else {
		menuList = append(menuList, menu.Row(CheckinBtn, AccountBtn), menu.Row(BindBtn, UnbindBtn), menu.Row(historyBtn))
	}

	menu.Reply(menuList...)

	model.Bot.Handle(&CheckinBtn, CheckinCmdCtr)
	model.Bot.Handle(&AccountBtn, AccountCmdCtr)
	model.Bot.Handle(&BindBtn, BindCmdCtr)
	model.Bot.Handle(&UnbindBtn, UnbindCmdCtr)
	model.Bot.Handle(&historyBtn, GetCheckinHistory)

	msg := fmt.Sprintf("%s\n为你提供以下服务:\n\n每日签到 /checkin\n账户信息 /account\n绑定账户 /bind\n解绑账户 /unbind\n签到历史 /history\n请注意, 绑定账号和解绑账号需要私聊我哦~", model.Config.Bot.Name)
	_, _ = model.Bot.Reply(m, msg, menu)
}
