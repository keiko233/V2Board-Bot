package config

import (
	"io/ioutil"
	"log"

	"github.com/keiko233/V2Board-Bot/model"
	"gopkg.in/yaml.v2"
)

func GetConfig(path string) (c *model.Conf) {
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
	return
}
