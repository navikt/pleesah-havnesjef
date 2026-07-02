package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

func (a *api) TeamHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /{team}/create", a.teamCreate)

	return mux
}

// Example: POST /api/v1/team/{team}/create?hex={code}
func (a *api) teamCreate(w http.ResponseWriter, r *http.Request) {
	team := r.PathValue("team")
	k8sconfig, err := a.k8s.SetupTeam(r.Context(), team)
	if err != nil {
		a.log.Error("failed creating team", "error", err, "team", team)
		http.Error(w, "failed creating team", http.StatusInternalServerError)
		return
	}

	a.log.Info("Created new team", "team", team)

	w.Header().Set("Content-Type", "application/json")
	buffer := new(bytes.Buffer)
	err = json.Compact(buffer, []byte(k8sconfig))
	w.Write(buffer.Bytes())
}
