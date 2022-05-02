package route

import (
	"github.com/keiko233/V2Board-Bot/common"
	"github.com/keiko233/V2Board-Bot/controller"
	"github.com/keiko233/V2Board-Bot/lib/tgbot"
	"github.com/keiko233/V2Board-Bot/model"
)

func Init() error {
	engine, err := tgbot.NewDefaultEngine(model.Config.Bot.Token)
	if err != nil {
		return err
	}

	route := engine.Use(handleErr)
	{
		route.Handle(controller.Help, "/help")
		route.Handle(controller.Checkin, "/checkin", model.MenuCheckinBtn)
		route.Handle(controller.Bind, "/bind", model.MenuBindBtn)
		route.Handle(controller.Unbind, "/unbind", model.MenuUnbindBtn)
		route.Handle(controller.Account, "/account", model.MenuAccountBtn)

		route.Handle(controller.Report, "/report", model.MenuReportBtn)
		route.HandleCallback(controller.ReportCallback, "report_btn")

		route.Handle(controller.CheckinHistory, "/history", model.MenuhistoryBtn)
		route.HandleCallback(controller.CheckinHistoryCallback, "history_page")
	}

	route.Run()
	return nil
}

func handleErr(h tgbot.HandleFunc) tgbot.HandleFunc {
	return func(c *tgbot.Context) error {
		if err := h(c); err != nil {
			common.ErrorResult(c, err)
		}
		return nil
	}
}
