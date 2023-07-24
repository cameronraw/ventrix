package main

import (
	"github.com/cameronraw/ventrix/cmd/queue"
	"gorm.io/gorm"
)

func AddServiceToDb(service queue.Service) (*gorm.DB, error) {

	result := db.Create(&service)

  return result, result.Error
}

func GetService(name string) (*queue.RegisteredService, error) {
	var service queue.RegisteredService
	result := db.Where("name = ?", name).First(&service)

	return &service, result.Error
}

