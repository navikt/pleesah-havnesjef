package api

import (
	"net/http"
	"slices"
)

// Example: GET /team/{team}/status/{deployment|pod|service}/?name={string}
func (a *api) teamResourceStatus(w http.ResponseWriter, r *http.Request) {
	type errorResponse struct {
		Error    string `json:"error"`
		Team     string `json:"team"`
		Resource string `json:"resource"`
		Name     string `json:"name"`
	}

	team := r.PathValue("team")
	resource := r.PathValue("resource")
	log := a.log.With("team", team, "resource", resource)

	if !slices.Contains([]string{"deployment", "pod", "service"}, resource) {
		log.Error("resource is not valid")
		writeJsonMessage(w, errorResponse{
			Error:    "resource is not valid",
			Team:     team,
			Resource: resource,
		}, http.StatusBadRequest)

		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		log.Error("missing name query parameter", "name", name)
		writeJsonMessage(w, errorResponse{
			Error:    "missing name query parameter",
			Team:     team,
			Resource: resource,
		}, http.StatusBadRequest)

		return
	}

	var err error
	var running bool
	switch resource {
	case "deployment":
		running, err = a.k8s.IsDeploymentRunning(r.Context(), team, name)
	case "pod":
		running, err = a.k8s.IsPodRunning(r.Context(), team, name)
	case "service":
		running, err = a.k8s.IsServiceRunning(r.Context(), team, name)
	}

	if err != nil {
		a.log.Error("failed checking status", "error", err, "team", team, "name", name, "resources", resource)
		writeJsonMessage(w, errorResponse{
			Error:    "resource is not valid",
			Team:     team,
			Resource: resource,
			Name:     name,
		}, http.StatusInternalServerError)

		return
	}

	writeJsonMessage(w, map[string]any{
		"isRunning": running,
		"team":      team,
		"resource":  resource,
		"name":      name,
	}, http.StatusOK)
}
