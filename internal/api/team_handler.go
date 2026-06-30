package api

import (
	"encoding/json"
	"net/http"
)

// Example: POST /api/v1/team/?team=team-a
func (a *api) TeamHandler(w http.ResponseWriter, r *http.Request) {
	team := r.URL.Query().Get("team")
	if team == "" {
		a.log.Error("missing team query parameter")
		http.Error(w, "missing team query parameter", http.StatusBadRequest)
		return
	}

	k8sconfig, err := a.k8s.SetupTeam(r.Context(), team)
	if err != nil {
		a.log.Error("failed creating team", "error", err, "team", team)
		http.Error(w, "failed creating team", http.StatusBadRequest)
		return
	}

	a.log.Info("Created new team", "team", team)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"kubeconfig": k8sconfig,
	})
}
