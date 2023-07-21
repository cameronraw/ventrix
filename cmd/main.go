package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"montecristo/cmd/config"
	"net/http"

	"github.com/golang-queue/queue"
	"github.com/golang-queue/queue/core"
	"github.com/golang-queue/redisdb"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

// Map of service name to service worker
var registeredServices = make(map[string]*queue.Queue)

// Map of event type to list of service names
var registeredEvents = make(map[string][]string)

var client = &http.Client{}

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

	router := mux.NewRouter()

	router.HandleFunc("/api/service", RegisterServiceHandler).Methods("POST")
	router.HandleFunc("/api/event", QueueEventHandler).Methods("POST")
	router.HandleFunc("/api/registerevent", RegisterEventHandler).Methods("POST")
	router.HandleFunc("/api/listen", ListenEventHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}

type RegisterServiceRequest struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
}

func createWorkerForService(service Service) {
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

func RegisterServiceHandler(w http.ResponseWriter, r *http.Request) {
	var request RegisterServiceRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newService := Service{
		Name:     request.Name,
		Endpoint: request.Endpoint,
	}

	log.Print("Registering service: ", newService)

	createWorkerForService(newService)

	log.Print("After creating worker: ", newService)

	if !cfg.UseInMemory {
		_, err = AddServiceToDb(newService)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(newService)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type QueueEventRequest struct {
	Type    string            `json:"type"`
	Payload map[string]string `json:"payload"`
	Timeout int               `json:"timeout"`
}

func (request *QueueEventRequest) ToEvent() Event {
	return Event{
		Type:    request.Type,
		Payload: request.Payload,
		Timeout: request.Timeout,
	}
}

func QueueEventHandler(w http.ResponseWriter, r *http.Request) {
	var request QueueEventRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event := request.ToEvent()

	err = QueueEvent(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type ListenEventRequest struct {
	ServiceName string `json:"service_name"`
	Type        string `json:"type"`
}

type RegisterEventRequest struct {
	EventType string `json:"event_type"`
}

func RegisterEventHandler(w http.ResponseWriter, r *http.Request) {
	var request RegisterEventRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if value, exists := registeredEvents[request.EventType]; !exists {
		log.Print("Registering new event type: ", request.EventType)
		registeredEvents[request.EventType] = []string{}
	} else {
		log.Print("Event type already registered: ", request.EventType)
		log.Print("Registered events: ", value)
	}

	log.Print("Registered events: ", registeredEvents)
}

func ListenEventHandler(w http.ResponseWriter, r *http.Request) {
	var request ListenEventRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if value, exists := registeredEvents[request.Type]; exists {
		registeredEvents[request.Type] = append(value, request.ServiceName)
		log.Print("Registered events: ", registeredEvents)
	} else {
		log.Print("Couldn't find event type: ", request.Type)
		log.Print("Registered events: ", registeredEvents)
	}
}
