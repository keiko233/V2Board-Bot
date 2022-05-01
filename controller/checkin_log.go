package controller

import (
	"bytes"
	"fmt"
	"image/png"
	"log"
	"strconv"
	"strings"

	"github.com/keiko233/V2Board-Bot/dao"
	"github.com/keiko233/V2Board-Bot/lib/image"
	"github.com/keiko233/V2Board-Bot/model"
	"github.com/keiko233/V2Board-Bot/service"
	"github.com/keiko233/V2Board-Bot/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

func GetCheckinHistory(m *tb.Message) {
	history(m, int(m.Sender.ID), 1, false)
}

func HistoryCallback(q *tb.Callback) {
	list := strings.Split(q.Data, ":")
	n, _ := strconv.Atoi(list[0])
	id, _ := strconv.Atoi(list[2])
	history(q.Message, id, n, true)
}

func history(m *tb.Message, id, n int, isCallBack bool) {
	count, out, err := dao.GetCheckLogsByTelegramID(int64(id), n, 5)
	if err != nil {
		log.Println("history GetCheckLogsByTelegramID err", err)
		model.Bot.Reply(m, "获取失败")
		return
	}

	sum, _, err := dao.GetCheckinLogsTrafficSumByTelegramID(int64(id))
	if err != nil {
		log.Println("history GetCheckinLogsTrafficSumByTelegramID err", err)
		model.Bot.Reply(m, "获取失败")
		return
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
		log.Println("test2 err", err)
		model.Bot.Reply(m, "生成图片失败")
		return
	}
	var b []byte
	bf := bytes.NewBuffer(b)
	err = png.Encode(bf, img.GetImage())
	if err != nil {
		log.Println("test3 err", err)
		_, err = model.Bot.Reply(m, "生成图片失败")
		return
	}

	if isCallBack {
		model.Bot.Edit(m, &tb.Photo{
			File:    tb.FromReader(bf),
			Caption: s,
		}, historyPage(n-1, n+1, int(max), id))
	} else {
		model.Bot.Reply(m, &tb.Photo{
			File:    tb.FromReader(bf),
			Caption: s,
		}, historyPage(n-1, n+1, int(max), id))
	}
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

func Report(m *tb.Message) {
	report(m, model.DailyReport, false)
}

func ReportCallback(c *tb.Callback) {
	report(c.Message, model.ReportType(c.Data), true)
}

func report(m *tb.Message, r model.ReportType, isCallBack bool) {
	report, notfound, start, end, err := service.Report(r)
	if err != nil {
		log.Printf("操作失败 %s\n", err)
		msg := "操作失败！请联系管理员！"
		reply(isCallBack, m, msg, reportBtn(r))
		return
	}

	if notfound {
		log.Printf("report: %+v\n", report)
		msg := "今天还没有人签到哦~"
		reply(isCallBack, m, msg, reportBtn(r))
		return
	}
	var max, min *tb.ChatMember
	if m.Chat.ID > 0 {
		max = new(tb.ChatMember)
		max.User = new(tb.User)
		max.User.Username = "null"
		max.User.FirstName = "匿名"

		min = new(tb.ChatMember)
		min.User = new(tb.User)
		min.User.Username = "null"
		min.User.FirstName = "匿名"
	} else {

		max, err = model.Bot.ChatMemberOf(m.Chat, &tb.User{ID: report.MaxUser})
		if err != nil {
			log.Println("bot ChatMemberOf err", err)
			msg := "操作失败！请联系管理员！"
			reply(isCallBack, m, msg, reportBtn(r))
			return
		}

		min, err = model.Bot.ChatMemberOf(m.Chat, &tb.User{ID: report.MinUser})
		if err != nil {
			log.Println("bot ChatMemberOf err", err)
			msg := "操作失败！请联系管理员！"
			reply(isCallBack, m, msg, reportBtn(r))
			return
		}
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

	reply(isCallBack, m, msg, reportBtn(r))
}

func reply(isCallback bool, to *tb.Message, what interface{}, options ...interface{}) (*tb.Message, error) {
	if isCallback {
		return model.Bot.Edit(to, what, options...)
	} else {
		return model.Bot.Reply(to, what, options...)
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
