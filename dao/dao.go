package dao

import (
	"errors"

	"github.com/keiko233/V2Board-Bot/model"
	"gorm.io/gorm"
)

func IsNotFound(err error) (bool, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil
	}
	return false, err
}

func GetPlanByID(db *gorm.DB, planId int) (plan *model.Plan, notfound bool, err error) {
	notfound, err = GetOutByQuery(db, &plan, "id = ?", planId)
	return
}

func GetUserByTelegramID(db *gorm.DB, tgId int64) (user *model.User, notfound bool, err error) {
	notfound, err = GetOutByQuery(db, &user, "telegram_id = ?", tgId)
	return
}

func GetUserByToken(db *gorm.DB, token string) (user *model.User, notfound bool, err error) {
	notfound, err = GetOutByQuery(db, &user, "token = ?", token)
	return
}

func GetOutByQuery(db *gorm.DB, out, query interface{}, args ...interface{}) (notfound bool, err error) {
	err = MustGetOutByQuery(db, out, query, args...)
	notfound, err = IsNotFound(err)
	return
}

func MustGetOutByQuery(db *gorm.DB, out, query interface{}, args ...interface{}) (err error) {
	conn := model.DB
	if db != nil {
		conn = db
	}
	err = conn.Where(query, args...).First(out).Error
	return
}

func GetOutsByQuery(db *gorm.DB, out, query interface{}, args ...interface{}) (notfound bool, err error) {
	err = MustGetOutsByQuery(db, query, out, args...)
	notfound, err = IsNotFound(err)
	return
}

func MustGetOutsByQuery(db *gorm.DB, out, query interface{}, args ...interface{}) (err error) {
	conn := model.DB
	if db != nil {
		conn = db
	}
	return conn.Where(query, args...).Find(out).Error
}

func MustGetUserByToken(db *gorm.DB, token string) (user *model.User, err error) {
	err = MustGetOutByQuery(db, &user, "token = ?", token)
	return
}

func MustGetUserByTelegramID(db *gorm.DB, tgId int64) (user *model.User, err error) {
	err = MustGetOutByQuery(db, &user, "telegram_id = ?", tgId)
	return
}

func UpdateModel(db *gorm.DB, m interface{}, query string, value interface{}) error {
	conn := model.DB
	if db != nil {
		conn = db
	}
	return conn.Model(m).Update(query, value).Error
}

func UpdateUser(db *gorm.DB, user *model.User, query string, value interface{}) error {
	return UpdateModel(db, user, query, value)
}

func Update(db *gorm.DB, m interface{}) error {
	conn := model.DB
	if db != nil {
		conn = db
	}
	return conn.Model(&m).Updates(m).Error
}

func GetLatestCheckLogByTelegramID(id int64) (log *model.CheckinLog, notfound bool, err error) {
	err = model.DB.Model(&model.CheckinLog{}).Where("telegram_id = ?", id).Order("created_at DESC").First(&log).Error
	notfound, err = IsNotFound(err)
	return
}

func GetCheckLogsByTelegramID(id int64, pageIndex, pageSize int) (count int64, logs []*model.CheckinLog, err error) {
	builder := model.DB.Model(&model.CheckinLog{}).Where("telegram_id = ?", id).Order("created_at DESC")
	count, err = Page(builder, pageIndex, pageSize, &logs)
	return
}

func Page(db *gorm.DB, pageIndex, pageSize int, out interface{}) (int64, error) {
	var count int64
	err := db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(out).Offset(-1).Limit(-1).Count(&count).Error
	return count, err
}