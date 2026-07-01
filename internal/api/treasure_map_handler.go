package api

import (
	"encoding/json"
	"net/http"
)

// Example: GET /api/v1/teams
func (a *api) TreasureMapHandler(w http.ResponseWriter, r *http.Request) {
	teams, err := a.k8s.ListTeams(r.Context())
	if err != nil {
		a.log.Error("failed listing teams", "error", err)
		http.Error(w, "failed listing teams", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(teams); err != nil {
		a.log.Error("failed encoding teams", "error", err)
		http.Error(w, "failed encoding teams", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
}
