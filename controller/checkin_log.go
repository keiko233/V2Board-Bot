package controller

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"strconv"
	"strings"

	"github.com/keiko233/V2Board-Bot/common"
	"github.com/keiko233/V2Board-Bot/dao"
	"github.com/keiko233/V2Board-Bot/lib/image"
	"github.com/keiko233/V2Board-Bot/lib/tgbot"
	"github.com/keiko233/V2Board-Bot/model"
	"github.com/keiko233/V2Board-Bot/service"
	"github.com/keiko233/V2Board-Bot/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

func CheckinHistory(ctx *tgbot.Context) error {
	return history(ctx, int(ctx.Message.Sender.ID), 1)
}

func CheckinHistoryCallback(ctx *tgbot.Context) error {
	list := strings.Split(ctx.Callback.Data, ":")
	n, _ := strconv.Atoi(list[0])
	id, _ := strconv.Atoi(list[2])
	return history(ctx, id, n)
}

func history(ctx *tgbot.Context, id, n int) error {
	count, out, err := dao.GetCheckLogsByTelegramID(int64(id), n, 5)
	if err != nil {
		log.Println("history GetCheckLogsByTelegramID err", err)
		return err
	}

	sum, _, err := dao.GetCheckinLogsTrafficSumByTelegramID(int64(id))
	if err != nil {
		log.Println("history GetCheckinLogsTrafficSumByTelegramID err", err)
		return err
	}

	max := count / 5
	if count%5 != 0 {
		max += 1
	}

	s := fmt.Sprintf("当前位于第%d页, 总条数%d, 总页数%d, 总签到流量%s", n, count, max, utils.ByteSize(sum))
	ss := make([][]string, 0)
	s1 := make([]string, 0)
	s2 := make([]string, 0)
	s1 = append(s1, "签到时间")
	s2 = append(s2, "获得流量")
	for _, i := range out {
		s1 = append(s1, i.CreatedAt.Format("2006-01-02 15:04:05"))
		s2 = append(s2, utils.ByteSize(i.CheckinTraffic))
	}
	ss = append(ss, s1, s2)
	img, err := image.NewDefaultTable(ss, model.Config.Bot.Font)
	if err != nil {
		return err
	}
	var b []byte
	bf := bytes.NewBuffer(b)
	err = png.Encode(bf, img.GetImage())
	if err != nil {
		return err
	}

	if ctx.IsCallback() {
		return ctx.Edit(&tb.Photo{
			File:    tb.FromReader(bf),
			Caption: s,
		}, historyPage(n-1, n+1, int(max), id))
	}

	return ctx.Reply(&tb.Photo{
		File:    tb.FromReader(bf),
		Caption: s,
	}, historyPage(n-1, n+1, int(max), id))

}

func historyPage(perv, next, max, id int) *tb.ReplyMarkup {
	r := make([][]tb.InlineButton, 0)
	r1 := make([]tb.InlineButton, 0)
	r2 := tb.InlineButton{
		Unique: "history_page",
		Data:   strconv.Itoa(perv) + ":" + strconv.Itoa(max) + ":" + strconv.Itoa(id),
		Text:   "上一页",
	}
	if perv > 0 {
		r1 = append(r1, r2)
	}
	r2.Data = strconv.Itoa(next) + ":" + strconv.Itoa(max) + ":" + strconv.Itoa(id)
	r2.Text = "下一页"
	if max != 0 && next <= max {
		r1 = append(r1, r2)
	}
	r = append(r, r1)
	return &tb.ReplyMarkup{
		InlineKeyboard: r,
	}
}

func Report(ctx *tgbot.Context) error {
	return report(ctx, model.DailyReport)
}

func ReportCallback(ctx *tgbot.Context) error {
	return report(ctx, model.ReportType(ctx.Callback.Data))
}

func report(ctx *tgbot.Context, r model.ReportType) error {
	report, notfound, start, end, err := service.Report(r)
	if err != nil {
		return err
	}

	if notfound {
		return common.ErrNotFoundCheckinUsers
	}

	max, err := ctx.ChatMemberOf(report.MaxUser)
	if err != nil {
		log.Println("bot ChatMemberOf err", err)
		max = new(tb.ChatMember)
		max.User = new(tb.User)
		max.User.Username = "null"
		max.User.FirstName = "匿名"
	}

	min, err := ctx.ChatMemberOf(report.MinUser)
	if err != nil {
		log.Println("bot ChatMemberOf err", err)
		min = new(tb.ChatMember)
		min.User = new(tb.User)
		min.User.Username = "null"
		min.User.FirstName = "匿名"
	}

	msg := fmt.Sprintf("%s: \n\n统计时间: \n开始: %s\n结束: %s\n\n签到总流量: %s\n签到总次数: %d\n\n欧皇: %s\n@%s\n获得: %s\n\n非酋: %s\n@%s\n获得: %s",
		checkReportType(r),
		start.Format("2006-01-02 15:04:05"),
		end.Format("2006-01-02 15:04:05"),
		utils.ByteSize(report.Sum),
		report.UserCount,
		max.User.FirstName+max.User.LastName,
		max.User.Username,
		utils.ByteSize(report.Max),
		min.User.FirstName+min.User.LastName,
		min.User.Username,
		utils.ByteSize(report.Min),
	)

	return reply(ctx, msg, reportBtn(r))
}

func reply(ctx *tgbot.Context, what interface{}, options ...interface{}) error {
	if ctx.IsCallback() {
		return ctx.Edit(what, options...)
	} else {
		return ctx.Reply(what, options...)
	}
}

func reportBtn(t model.ReportType) *tb.ReplyMarkup {

	ss := make([]model.ReportType, 0, 3)
	ss = append(ss, model.DailyReport)
	ss = append(ss, model.WeeklyReport)
	ss = append(ss, model.MonthlyReport)

	r := make([][]tb.InlineButton, 0)
	r1 := make([]tb.InlineButton, 0)
	r2 := tb.InlineButton{
		Unique: "report_btn",
	}
	for _, i := range ss {
		if i == t {
			continue
		}
		r2.Data = string(i)
		r2.Text = checkReportType(i)
		r1 = append(r1, r2)
	}

	return &tb.ReplyMarkup{
		InlineKeyboard: append(r, r1),
	}
}

func checkReportType(t model.ReportType) string {
	switch t {
	case model.DailyReport:
		return "本日汇报"
	case model.WeeklyReport:
		return "本周汇报"
	case model.MonthlyReport:
		return "本月汇报"
	}
	return ""
}
