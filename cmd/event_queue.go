package main

import (
	"log"
	"time"
)

const (
  TaskStatusPending = "pending"
  TaskStatusComplete = "complete"
  TaskStatusFailed = "failed"
)

type Task struct {
  Service uint
  Event   uint
  Status  string
  SendAt  time.Time
}


func QueueEvent(event Event) error {

  servicesToNotify := registeredEvents[event.Type]

  log.Printf("Found %d services to notify for event %s", len(servicesToNotify), event.Type)

  for _, serviceId := range servicesToNotify {
    log.Printf("Queueing event %s for service %s", event.Type, serviceId)
    q := registeredServices[serviceId]
    err := q.Queue(event)
    if err != nil {
      log.Printf("Error queueing event %s for service %s: %s", event.Type, serviceId, err)
    }
  }
  
  // Fetch services interested in this event type
  // services, err := GetServicesByEventType(event.Type)
  // if err != nil {
  //   return err
  // }

  // Loop over services and queue event for each
  // for _, service := range services {
  //   // Create a task
  //   task := Task{
  //     Service: service.ID,
  //     Event:   event.ID,
  //     Status:  TaskStatusPending,
  //     SendAt:  time.Now().Add(time.Duration(event.Timeout) * time.Second),
  //   }

  //   // Save task to database
  //   result := db.Create(&task)
  //   if result.Error != nil {
  //     return result.Error
  //   }
  // }

  return nil
}

type Service struct {
  ID       uint
  Name     string
  Endpoint string
}

func GetServicesByEventType(eventType string) ([]Service, error) {
  // Fetch services interested in this event type from the database
  // ... implementation depends on your database schema
  return []Service{}, nil
}

