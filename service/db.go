package service

import (
	"fmt"
	"log"
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

type UUBot struct {
	Id             uint `gorm:"primaryKey"`
	UserId         uint `gorm:"unique"`
	TelegramId     uint `gorm:"unique" `
	CheckinTraffic int64
	CheckinAt      int64
	NextAt         int64
}

var DB *gorm.DB
var c Conf

func init() {
	c.GetConfig()
}

func InitDB() (*gorm.DB, error) {
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
	if err = db.AutoMigrate(&UUBot{}); err != nil {
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

func QueryPlan(planId int) Plan {
	var plan Plan
	if err := DB.Where("id = ?", planId).First(&plan).Error; err != nil {
		log.Printf("QueryPlan id = %d error, %s\n", planId, err)
	}
	return plan
}

func QueryUser(tgId int64) User {
	var user User
	if err := DB.Where("telegram_id = ?", tgId).First(&user).Error; err != nil {
		log.Printf("QueryUser tgid = %d error, %s\n", tgId, err)
	}
	return user
}

func BindUser(token string, tgId int64) User {
	var user User
	if err := DB.Where("token = ?", token[6:38]).First(&user).Error; err != nil {
		log.Printf("BindUser Select User By tgid = %d error, %s\n", tgId, err)
	}
	if user.Id <= 0 {
		return user
	}
	if user.TelegramId <= 0 {
		if err := DB.Model(&user).Update("telegram_id", tgId).Error; err != nil {
			log.Printf("BindUser Update User By tgid = %d error, %s\n", tgId, err)
		}
	}
	return user
}

func unbindUser(tgId int64) User {
	var user User
	if err := DB.Where("telegram_id = ?", tgId).First(&user).Error; err != nil {
		log.Printf("unbindUser Select User By tgid = %d error, %s\n", tgId, err)
	}
	if user.Id > 0 {
		if err := DB.Model(&user).Update("telegram_id", nil).Error; err != nil {
			log.Printf("unbindUser Update User By tgid = %d error, %s\n", tgId, err)
		}
		return user
	}
	return user
}

func CheckinTime(tgId int64) bool {
	var uu UUBot
	if err := DB.Where("telegram_id = ?", tgId).First(&uu).Error; err != nil {
		log.Printf("CheckinTime Select User By tgid = %d error, %s\n", tgId, err)
		return true
	}
	checkin := time.Unix(uu.CheckinAt, 0)
	tomorrow, _ := time.ParseInLocation("2006-01-02", checkin.Format("2006-01-02"), time.Local)
	if tomorrow.AddDate(0, 0, 1).After(time.Now()) {
		return false
	}
	return true
}

func checkinUser(tgId int64) (UUBot, error) {
	var user User
	var uu UUBot
	if err := DB.Where("telegram_id = ?", tgId).First(&user).Error; err != nil {
		log.Printf("checkinUser Select User By tgid = %d error, %s\n", tgId, err)
		return uu, err
	}
	if err := DB.Where("telegram_id = ?", tgId).First(&uu).Error; err != nil {
		log.Printf("checkinUser Select UUBot By tgid = %d error, %s\n", tgId, err)
	}

	b := RandInt(c.Bot.MaxByte, c.Bot.MinByte)
	CheckIns := b * 1024 * 1024
	T := user.TransferEnable + CheckIns
	checkInAt := time.Now()
	nextAt, _ := time.ParseInLocation("2006-01-02", checkInAt.Format("2006-01-02"), time.Local)
	if uu.Id <= 0 {
		newUU := UUBot{
			UserId:         user.Id,
			TelegramId:     user.TelegramId,
			CheckinAt:      checkInAt.Unix(),
			NextAt:         nextAt.AddDate(0, 0, 1).Unix(),
			CheckinTraffic: 0,
		}
		if err := DB.Create(&newUU).Error; err != nil {
			log.Printf("checkinUser Create UUBot: %+v error, %s\n", newUU, err)
			return uu, err
		}
	} else {
		if err := DB.Model(&uu).Updates(UUBot{
			CheckinAt:      checkInAt.Unix(),
			NextAt:         nextAt.AddDate(0, 0, 1).Unix(),
			CheckinTraffic: CheckIns,
		}).Error; err != nil {
			log.Printf("checkinUser Update UUBot By %+v error, %s\n", uu, err)
			return uu, err
		}
	}
	if err := DB.Model(&user).Update("transfer_enable", T).Error; err != nil {
		log.Printf("checkinUser Update User By %+v error, %s\n", user, err)
		return uu, err
	}

	return uu, nil
}
