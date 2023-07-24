package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/cameronraw/ventrix/cmd/queue"
)

func RegisterServiceHandler(w http.ResponseWriter, r *http.Request) {
	var request RegisterServiceRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newService := queue.Service{
		Name:     request.Name,
		Endpoint: request.Endpoint,
	}

	log.Print("Registering service: ", newService)

	queue.CreateWorkerForService(newService)

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

	queue.RegisterEvent(request.EventType)
}

func ListenEventHandler(w http.ResponseWriter, r *http.Request) {
	var request ListenEventRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queue.ListenToEvent(request.ServiceName, request.Type)
}

func QueueEventHandler(w http.ResponseWriter, r *http.Request) {
	var request QueueEventRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event := request.ToEvent()

	err = queue.QueueEvent(event)
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
