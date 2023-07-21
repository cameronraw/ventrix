package main

import (
	"encoding/json"

	"gorm.io/gorm"
)

type RegisteredService struct {
	gorm.Model
	Name     string
	Endpoint string
	Events   []Event `gorm:"foreignKey:ServiceRefer"`
}

type Event struct {
	gorm.Model
	Type         string
	Payload      map[string]string
	Timeout      int
}

func (event Event) Bytes() []byte {
  bytes, err := json.Marshal(event)
  if err != nil {
    return []byte{}
  }
  return bytes
}

