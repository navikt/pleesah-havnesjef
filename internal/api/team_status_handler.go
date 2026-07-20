package api

import "net/http"

// Example: GET /api/v1/{team}/status/
func (a *api) teamStatus(w http.ResponseWriter, r *http.Request) {
	team := r.PathValue("team")
	log := a.log.With("team", team)

	pods, podsErr := a.k8s.PodInfo(r.Context(), team)
	services, servicesErr := a.k8s.ServiceInfo(r.Context(), team)
	deployments, deploymentsErr := a.k8s.DeploymentInfo(r.Context(), team)

	if podsErr != nil || servicesErr != nil || deploymentsErr != nil {
		log.Error("failed checking status", "podsErr", podsErr, "servicesErr", servicesErr, "deploymentsErr", deploymentsErr)
		writeJsonMessage(w, map[string]any{
			"error":          "failed checking status",
			"podsErr":        podsErr,
			"servicesErr":    servicesErr,
			"deploymentsErr": deploymentsErr,
		}, http.StatusInternalServerError)

		return
	}

	writeJsonMessage(w, map[string]any{
		"pods":        pods,
		"services":    services,
		"deployments": deployments,
	}, http.StatusOK)
}
