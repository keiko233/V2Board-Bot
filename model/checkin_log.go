package model

import "gorm.io/gorm"

type CheckinLog struct {
	gorm.Model

	UserID         uint
	TelegramID     uint
	CheckinTraffic int64

	OldTraffic int64
	NewTraffic int64
}