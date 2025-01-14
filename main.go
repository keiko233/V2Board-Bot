package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/keiko233/V2Board-Bot/lib/cache"
	"github.com/keiko233/V2Board-Bot/lib/config"
	"github.com/keiko233/V2Board-Bot/lib/db"
	"github.com/keiko233/V2Board-Bot/model"
	"github.com/keiko233/V2Board-Bot/route"
)

func main() {
	path, err := os.Executable()
	path = filepath.Dir(path) + "/uuBot.yaml"
	if err != nil {
		log.Fatalln(err)
	}
	model.Config = config.GetConfig(path)
	model.Cache = cache.NewMapCache()
	model.DB, err = db.InitDB(model.Config.Database)
	if err != nil {
		log.Fatalln(err)
	}

	if err := route.Init(); err != nil {
		log.Fatalln(err)
	}
}
