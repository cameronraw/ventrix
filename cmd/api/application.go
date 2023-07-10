package api

import (
	"fmt"
	"log"
	"montecristo/cmd/internal/json"
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

func (app *Application) applyMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	if len(app.middleware) == 0 {
		return handler
	}

	whatToWrap := app.middleware[len(app.middleware)-1].Wrap(handler)

	if len(app.middleware) == 1 {
		return whatToWrap
	}

	for i := len(app.middleware) - 2; i >= 0; i-- {
		whatToWrap = app.middleware[i].Wrap(whatToWrap)
	}

	return whatToWrap
}

func (app *Application) healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	data := map[string]string{
		"status":      "available",
		"environment": app.config.Env,
		"version":     "1.0.0",
	}

	err := json.WriteJSON(w, r, http.StatusOK, json.Envelope{"response": data})
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
	}
}
