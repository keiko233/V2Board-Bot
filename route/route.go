package route

import (
	"log"
	"time"

	"github.com/keiko233/V2Board-Bot/controller"
	"github.com/keiko233/V2Board-Bot/model"
	tb "gopkg.in/tucnak/telebot.v2"
)

func Start() {
	bot, err := tb.NewBot(tb.Settings{
		URL:    "https://api.telegram.org",
		Token:  model.Config.Bot.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatalf("Bot 启动失败啦...... \n当前Token [ %s ] \n错误信息:  %s", model.Config.Bot.Token, err)
	}
	model.Bot = bot
	setHandle(bot)
	bot.Start()
}

func setHandle(bot *tb.Bot) {
	bot.Handle("/help", controller.StartCmdCtr)
	bot.Handle("/checkin", controller.CheckinCmdCtr)
	bot.Handle("/account", controller.AccountCmdCtr)
	bot.Handle("/bind", controller.BindCmdCtr)
	bot.Handle("/unbind", controller.UnbindCmdCtr)
	bot.Handle("/history", controller.GetCheckinHistory)
	bot.Handle("\fhistory_page", controller.HistoryCallback)
	bot.Handle("/report", controller.Report)
	bot.Handle("\freport_btn", controller.ReportCallback)

	bot.Handle(model.MenuCheckinBtn, controller.CheckinCmdCtr)
	bot.Handle(model.MenuAccountBtn, controller.AccountCmdCtr)
	bot.Handle(model.MenuBindBtn, controller.BindCmdCtr)
	bot.Handle(model.MenuUnbindBtn, controller.UnbindCmdCtr)
	bot.Handle(model.MenuhistoryBtn, controller.GetCheckinHistory)
	bot.Handle(model.MenureportBtn, controller.Report)

}
