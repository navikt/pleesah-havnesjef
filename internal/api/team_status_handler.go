package api

import "net/http"

// Example: GET /api/v1/{team}/status/
func (a *api) teamStatus(w http.ResponseWriter, r *http.Request) {
	team := r.PathValue("team")
	log := a.log.With("team", team)

}
