package service

import (
	"time"

	"gorm.io/gorm"
)

func BindUser(token string, tgId int64) (*User, error) {
	user, err := MustGetUserByToken(token)
	if err != nil {
		return nil, err
	}

	if user.TelegramId != 0 {
		return user, nil
	}

	err = UpdateUser(user, "telegram_id", tgId)
	return user, err
}

func unbindUser(tgId int64) (notfound bool, err error) {
	user, notfound, err := GetUserByTelegramID(tgId)
	if err != nil {
		return false, err
	}

	if notfound {
		return true, err
	}

	err = UpdateUser(user, "telegram_id", nil)
	return false, err
}

func CheckinTime(tgId int64) (todayNotCheckin bool, err error) {

	l, notfound, err := GetLatestCheckLogByTelegramID(tgId)
	if err != nil {
		return false, err
	}
	if notfound {
		return true, nil
	}

	checkDay, err := time.ParseInLocation("2006-01-02", l.CreatedAt.Format("2006-01-02"), time.Local)
	if err != nil {
		return false, err
	}
	return checkDay.AddDate(0, 0, 1).Before(time.Now()), nil
}

func checkinUser(tgId int64) (log *CheckinLog, err error) {

	user, err := MustGetUserByTelegramID(tgId)
	if err != nil {
		return nil, err
	}

	b := RandInt(c.Bot.MaxByte, c.Bot.MinByte)
	checkIns := b * 1024 * 1024
	oldTraffic := user.TransferEnable
	newTraffic := user.TransferEnable + checkIns
	err = DB.Transaction(func(tx *gorm.DB) error {
		if err := UpdateUser(user, "transfer_enable", newTraffic); err != nil {
			return err
		}
		log = &CheckinLog{
			UserID:         user.Id,
			TelegramID:     uint(tgId),
			CheckinTraffic: checkIns,
			OldTraffic:     oldTraffic,
			NewTraffic:     newTraffic,
		}
		return tx.Model(&CheckinLog{}).Create(log).Error
	})

	return
}
