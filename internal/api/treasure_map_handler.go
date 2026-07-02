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
		writeJsonMessage(w, map[string]any{
			"error": "failed listing teams",
		}, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_ = json.NewEncoder(w).Encode(teams)
}
