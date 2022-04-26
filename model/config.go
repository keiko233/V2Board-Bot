package model

import (
	"gorm.io/gorm"
	tb "gopkg.in/tucnak/telebot.v2"
)

var DB *gorm.DB
var Config *Conf
var Bot *tb.Bot

type Conf struct {
	Bot      BotConf      `yaml:"bot"`
	Database DatabaseConf `yaml:"database"`
}

type BotConf struct {
	Token   string `yaml:"token"`
	Name    string `yaml:"name"`
	MinByte int64  `yaml:"min_byte"` // 签到流量的最小值，不配置时为0
	MaxByte int64  `yaml:"max_byte"` // 签到流量的最大值，为负数时为0
}
type DatabaseConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
