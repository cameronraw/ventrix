package api

import (
	"montecristo/cmd/internal/json"
	"net/http"
)

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
