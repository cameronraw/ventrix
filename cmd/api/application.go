package api

import (
	"fmt"
	"log"
	"montecristo/cmd/config"
	"net/http"
)

type Application struct {
	config     config.Config
	logger     *log.Logger
	Middleware []Middleware
}

func CreateApplication(config config.Config, logger *log.Logger) *Application {
	app := &Application{
		config:     config,
		logger:     logger,
		Middleware: []Middleware{},
	}

	ConfigureMiddleware(app)

	return app
}

func (app *Application) CreateAddr() string {
	return fmt.Sprintf(":%d", app.config.Port)
}

func (app *Application) CreatePreconfiguredHandler(config *config.Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", app.applyMiddleware(app.healthCheck))
	return mux
}

