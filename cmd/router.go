package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func InitApiRouter(port int) {
	router := mux.NewRouter()

	router.HandleFunc("/api/service", RegisterServiceHandler).Methods("POST")
	router.HandleFunc("/api/event", QueueEventHandler).Methods("POST")
	router.HandleFunc("/api/registerevent", RegisterEventHandler).Methods("POST")
	router.HandleFunc("/api/listen", ListenEventHandler).Methods("POST")

  portAsString := strconv.Itoa(port)
  portAsString = fmt.Sprintf(":%s", portAsString)

	log.Fatal(http.ListenAndServe(portAsString, router))
}
