package model

import (
	"github.com/keiko233/V2Board-Bot/lib/cache"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Config *Conf
var Cache cache.Cache

var (
	MenuCheckinBtn = "👀 每日签到"
	MenuAccountBtn = "🚥‍ 账户信息"
	MenuBindBtn    = "😋 绑定账户"
	MenuUnbindBtn  = "🤔 解绑账户"
	MenuhistoryBtn = "📅 签到历史"
	MenuReportBtn  = "📊 数据统计"
)

type Conf struct {
	Bot      BotConf      `yaml:"bot"`
	Database DatabaseConf `yaml:"database"`
}

type BotConf struct {
	Token   string `yaml:"token"`
	Name    string `yaml:"name"`
	MinByte int64  `yaml:"min_byte"` // 签到流量的最小值，不配置时为0
	MaxByte int64  `yaml:"max_byte"` // 签到流量的最大值，为负数时为0
	Font    string `yaml:"font"`     // 字体路径，只支持ttf
}
type DatabaseConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
