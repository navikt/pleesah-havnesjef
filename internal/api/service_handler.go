package api

import (
	"encoding/json"
	"net/http"
)

// Example: GET /api/v1/service/status?team=team-a&service=my-service
func (a *api) ServiceRunningHandler(w http.ResponseWriter, r *http.Request) {
	team := r.URL.Query().Get("team")
	service := r.URL.Query().Get("service")
	if team == "" || service == "" {
		a.log.Error("missing team or service query parameter", "team", team, "service", service)
		http.Error(w, "missing team or service query parameter", http.StatusBadRequest)
		return
	}

	running, err := a.k8s.IsServiceRunning(r.Context(), team, service)
	if err != nil {
		a.log.Error("failed checking service", "error", err, "team", team, "service", service)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	emoji := "❌"
	if running {
		emoji = "✅"
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"running": emoji,
	})
}
