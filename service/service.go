package service

import (
	"time"

	"github.com/keiko233/V2Board-Bot/dao"
	"github.com/keiko233/V2Board-Bot/lib/rand"
	"github.com/keiko233/V2Board-Bot/model"
	"gorm.io/gorm"
)

func BindUser(token string, tgId int64) (*model.User, error) {
	user, err := dao.MustGetUserByToken(nil, token)
	if err != nil {
		return nil, err
	}

	if user.TelegramId != 0 {
		return user, nil
	}

	err = dao.UpdateUser(nil, user, "telegram_id", tgId)
	return user, err
}

func UnbindUser(tgId int64) (notfound bool, err error) {
	user, notfound, err := dao.GetUserByTelegramID(nil, tgId)
	if err != nil {
		return false, err
	}

	if notfound {
		return true, err
	}

	err = dao.UpdateUser(nil, user, "telegram_id", nil)
	return false, err
}

func CheckinTime(tgId int64) (todayNotCheckin bool, err error) {

	l, notfound, err := dao.GetLatestCheckLogByTelegramID(tgId)
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

func CheckinUser(tgId int64) (log *model.CheckinLog, err error) {

	user, err := dao.MustGetUserByTelegramID(nil, tgId)
	if err != nil {
		return nil, err
	}

	b := rand.RandInt(model.Config.Bot.MaxByte, model.Config.Bot.MinByte)
	checkIns := b * 1024 * 1024
	oldTraffic := user.TransferEnable
	newTraffic := user.TransferEnable + checkIns
	err = model.DB.Transaction(func(tx *gorm.DB) error {
		if err := dao.UpdateUser(tx, user, "transfer_enable", newTraffic); err != nil {
			return err
		}
		log = &model.CheckinLog{
			UserID:         user.Id,
			TelegramID:     uint(tgId),
			CheckinTraffic: checkIns,
			OldTraffic:     oldTraffic,
			NewTraffic:     newTraffic,
		}
		return tx.Model(&model.CheckinLog{}).Create(log).Error
	})

	return
}
