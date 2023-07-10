package api

import (
	"fmt"
	"log"
	"net/http"
)

type Application struct {
	config     Config
	logger     *log.Logger
	middleware []Middleware
}

func CreateApplication(config Config, logger *log.Logger) *Application {
	app := &Application{
		config:     config,
		logger:     logger,
		middleware: []Middleware{},
	}

	app.ConfigureMiddleware()

	return app
}

func (app *Application) CreateAddr() string {
	return fmt.Sprintf(":%d", app.config.Port)
}

func (app *Application) CreatePreconfiguredHandler(config *Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", app.applyMiddleware(app.healthCheck))
	return mux
}

