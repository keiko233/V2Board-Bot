package service

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
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
}
type DatabaseConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (c *Conf) GetConfig() *Conf {
	path, err := os.Executable()
	path = filepath.Dir(path) + "/uuBot.yaml"
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("read config with path: %s", path)
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("打开配置文件错误...\n错误信息:%s", err)
	}
	if err = yaml.Unmarshal(yamlFile, c); err != nil {
		log.Fatalf("配置文件解析错误... \n错误信息:%s", err)
	}

	// Safe Random
	if c.Bot.MaxByte <= 0 {
		log.Fatalln("config.bot.max_byte must > 0, get ", c.Bot.MaxByte)
	}
	return c
}
