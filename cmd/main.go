package main

import (
	"log"

	"github.com/cameronraw/ventrix/cmd/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

var cfg config.Config

func main() {
	var err error
	cfg, err = config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	if !cfg.UseInMemory {
		db, err = gorm.Open(mysql.Open(cfg.SqlDsn), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
	}

  InitApiRouter(cfg.Port)
}


