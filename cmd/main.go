package main

import (
	"log"
	"montecristo/cmd/api"
	"net/http"
	"os"
	"time"
)

func main() {
	config := api.NewConfig()
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app := api.CreateApplication(config, logger)

	srv := &http.Server{
		Addr:         app.CreateAddr(),
		Handler:      app.CreatePreconfiguredHandler(&config),
		IdleTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		logger.Fatal(err)
	}
}
