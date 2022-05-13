package model

import "gorm.io/gorm"

type CheckinLog struct {
	gorm.Model

	UserID         uint
	TelegramID     uint
	CheckinTraffic int64

	OldTraffic int64
	NewTraffic int64

	Fortune FortuneType
}

type ReportType string

const (
	DailyReport   ReportType = "day"
	WeeklyReport  ReportType = "week"
	MonthlyReport ReportType = "month"
)

type Report struct {
	Sum       int64 `json:"sum"`      // 总流量
	UserCount int64 `json:"count"`    // 总人数
	MaxUser   int64 `json:"max_user"` // 欧狗
	Max       int64 `json:"max"`
	MinUser   int64 `json:"min_user"` // 非酋
	Min       int64 `json:"mix"`
}
