package service

import (
	"time"

	"github.com/keiko233/V2Board-Bot/dao"
	"github.com/keiko233/V2Board-Bot/model"
	"github.com/keiko233/V2Board-Bot/utils"
)

func Report(reporyType model.ReportType) (*model.Report, bool, time.Time, time.Time, error) {
	var start, end time.Time
	switch reporyType {
	case model.DailyReport:
		start, end = utils.Today()
	case model.WeeklyReport:
		start, end = utils.ThisWeek()
	case model.MonthlyReport:
		start, end = utils.ThisMonth()
	}

	r, b, err := dao.GetReportByTime(start, end)
	return r, b, start, end, err
}
