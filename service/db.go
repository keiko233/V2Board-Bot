package service

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type User struct {
	Id                uint
	TelegramId        uint
	Email             string
	Token             string
	U                 int64
	D                 int64
	PlanId            int64
	Balance           int64
	TransferEnable    int64
	CommissionBalance int64
	ExpiredAt         int64
	CreatedAt         int64
}

type Plan struct {
	Id   uint
	Name string
}

type CheckinLog struct {
	gorm.Model

	UserID         uint
	TelegramID     uint
	CheckinTraffic int64

	OldTraffic int64
	NewTraffic int64
}

var DB *gorm.DB
var c Conf

func InitDB() (*gorm.DB, error) {
	c.GetConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Database.Username, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "v2_",
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("InitDB Open %w", err)
	}
	if err = db.AutoMigrate(&CheckinLog{}); err != nil {
		return nil, fmt.Errorf("InitDB AutoMigrate %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("InitDB Get DB %w", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Second)

	DB = db
	return db, nil
}

func GetLatestCheckLogByTelegramID(id int64) (log *CheckinLog, notfound bool, err error) {
	err = DB.Model(&CheckinLog{}).Where("telegram_id = ?", id).Order("created_at DESC").First(&log).Error
	notfound, err = IsNotFound(err)
	return
}

func GetCheckLogsByTelegramID(id int64, pageIndex, pageSize int) (count int64, logs []*CheckinLog, err error) {
	builder := DB.Model(&CheckinLog{}).Where("telegram_id = ?", id).Order("created_at DESC")
	count, err = Page(builder, pageIndex, pageSize, &logs)
	return
}

func Page(db *gorm.DB, pageIndex, pageSize int, out interface{}) (int64, error) {
	var count int64
	err := db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(out).Offset(-1).Limit(-1).Count(&count).Error
	return count, err
}
