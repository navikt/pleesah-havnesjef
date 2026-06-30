package api

import (
	"encoding/json"
	"net/http"
)

// Example: GET /api/v1/status/service?team=team-a&service=my-service
func (a *api) StatusHandler(w http.ResponseWriter, r *http.Request) {
	team := r.URL.Query().Get("team")
	name := r.URL.Query().Get("name")
	resource := r.URL.Query().Get("resource")

	if team == "" || name == "" || resource == "" {
		a.log.Error("missing team, name, resource query parameter", "team", team, "name", name, "resource", resource)
		http.Error(w, "missing team, name, resource query parameter", http.StatusBadRequest)
		return
	}

	var err error
	var running bool
	switch resource {
	case "service":
		running, err = a.k8s.IsServiceRunning(r.Context(), team, name)
		if err != nil {
			a.log.Error("failed checking service", "error", err, "team", team, "service", name)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	case "pod":
		running, err = a.k8s.IsPodRunning(r.Context(), team, name)
		if err != nil {
			a.log.Error("failed checking pod", "error", err, "team", team, "name", name)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	emoji := "❌"
	if running {
		emoji = "✅"
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"running":  emoji,
		"resource": resource,
		"name":     name,
	})
}
