package main

import (
	"encoding/json"
	"log"
	"net/http"
)


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

	CreateWorkerForService(newService)

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
