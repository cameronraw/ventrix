package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-queue/queue"
	"github.com/golang-queue/queue/core"
	"github.com/golang-queue/redisdb"
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

func CreateWorkerForService(service Service) {
	w := redisdb.NewWorker(
		redisdb.WithAddr("redis:6379"),
		redisdb.WithChannel(service.Name),
		redisdb.WithRunFunc(func(ctx context.Context, m core.QueuedMessage) error {
			log.Printf("Firing worker for service: %s", service.Name)
			v, ok := m.(*Event)
			if !ok {
				if err := json.Unmarshal(m.Bytes(), &v); err != nil {
					return err
				}
			}

			payload, err := json.Marshal(v.Payload)
			if err != nil {
				return err
			}

			postBody, _ := json.Marshal(map[string]string{
				"event_type": v.Type,
				"payload":    string(payload),
			})

			requestBody := bytes.NewBuffer(postBody)

			log.Print("Request body: ", requestBody)

			r, err := http.NewRequest("POST", service.Endpoint, requestBody)
			if err != nil {
				panic(err)
			}

			r.Header.Add("Content-Type", "application/json")

			client := &http.Client{}
			res, err := client.Do(r)
			if err != nil {
				panic(err)
			}

			defer res.Body.Close()

			log.Print("Response from service: ", res.Body)

			return nil
		}),
	)

	q, err := queue.NewQueue(
		queue.WithWorkerCount(10),
		queue.WithWorker(w),
	)
	if err != nil {
		log.Fatal(err)
	}

	q.Start()

	registeredServices[service.Name] = q
}
