package service

import (
	"errors"

	"gorm.io/gorm"
)

func IsNotFound(err error) (bool, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true, nil
	}
	return false, err
}

func GetPlanByID(planId int) (plan *Plan, notfound bool, err error) {
	notfound, err = GetOutByQuery(&plan, "id = ?", planId)
	return
}

func GetUserByTelegramID(tgId int64) (user *User, notfound bool, err error) {
	notfound, err = GetOutByQuery(&user, "telegram_id = ?", tgId)
	return
}

func GetUserByToken(token string) (user *User, notfound bool, err error) {
	notfound, err = GetOutByQuery(&user, "token = ?", token)
	return
}

func GetOutByQuery(out, query interface{}, args ...interface{}) (notfound bool, err error) {
	err = MustGetOutByQuery(out, query, args...)
	notfound, err = IsNotFound(err)
	return
}

func MustGetOutByQuery(out, query interface{}, args ...interface{}) (err error) {
	err = DB.Where(query, args...).First(out).Error
	return
}

func GetOutsByQuery(out, query interface{}, args ...interface{}) (notfound bool, err error) {
	err = MustGetOutsByQuery(query, out, args...)
	notfound, err = IsNotFound(err)
	return
}

func MustGetOutsByQuery(out, query interface{}, args ...interface{}) (err error) {
	return DB.Where(query, args...).Find(out).Error
}

func MustGetUserByToken(token string) (user *User, err error) {
	err = MustGetOutByQuery(&user, "token = ?", token)
	return
}

func MustGetUserByTelegramID(tgId int64) (user *User, err error) {
	err = MustGetOutByQuery(&user, "telegram_id = ?", tgId)
	return
}

func UpdateModel(model interface{}, query string, value interface{}) error {
	return DB.Model(model).Update(query, value).Error
}

func UpdateUser(user *User, query string, value interface{}) error {
	return UpdateModel(user, query, value)
}

func Update(model interface{}) error {
	return DB.Model(&model).Updates(model).Error
}
