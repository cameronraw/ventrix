package main

import (
	"log"
	"montecristo/cmd/config"

	"github.com/golang-queue/queue"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// Map of service name to service worker
var registeredServices = make(map[string]*queue.Queue)

// Map of event type to list of service names
var registeredEvents = make(map[string][]string)

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


