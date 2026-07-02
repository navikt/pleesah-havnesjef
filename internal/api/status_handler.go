package api

import (
	"net/http"
	"slices"
)

// Example: GET /api/v1/{team}/status/{deployment|pod|service}/?name={string}
func (a *api) teamResourceStatus(w http.ResponseWriter, r *http.Request) {
	team := r.PathValue("team")
	resource := r.PathValue("resource")
	log := a.log.With("team", team, "resource", resource)

	if !slices.Contains([]string{"deployment", "pod", "service"}, resource) {
		log.Error("resource is not valid")
		writeJsonMessage(w, map[string]any{
			"err": "resource is not valid",
		}, http.StatusBadRequest)

		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		log.Error("missing name query parameter", "name", name)
		writeJsonMessage(w, map[string]any{
			"err": "missing name query parameter",
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
		writeJsonMessage(w, map[string]any{
			"error":    err,
			"resource": resource,
			"name":     name,
		}, http.StatusInternalServerError)

		return
	}

	emoji := "❌"
	if running {
		emoji = "✅"
	}

	writeJsonMessage(w, map[string]any{
		"running":  emoji,
		"resource": resource,
		"name":     name,
	}, http.StatusOK)
}
