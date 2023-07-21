package main

import "gorm.io/gorm"

func AddServiceToDb(service Service) (*gorm.DB, error) {

	result := db.Create(&service)

  return result, result.Error
}

func GetService(name string) (*RegisteredService, error) {
	var service RegisteredService
	result := db.Where("name = ?", name).First(&service)

	return &service, result.Error
}

