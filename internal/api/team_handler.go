package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
)

func (a *api) TeamHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /{team}/create", a.teamCreate)
	mux.HandleFunc("POST /{team}/next-task", a.teamNextTask)
	mux.HandleFunc("GET /{team}/status", a.teamStatus)

	return mux
}

// Example: POST /team/{team}/create
func (a *api) teamCreate(w http.ResponseWriter, r *http.Request) {
	type errorResponse struct {
		Error string `json:"error"`
		Team  string `json:"team"`
	}

	team := r.PathValue("team")
	log := a.log.With("team", team)

	if !validateTeam(team) {
		a.log.Error("team is not valid")
		writeJsonMessage(w, errorResponse{
			Error: "team is not valid",
		}, http.StatusBadRequest)
		return
	}

	k8sconfig, err := a.k8s.SetupTeam(r.Context(), team)
	if err != nil {
		log.Error("failed creating team", "error", err)
		writeJsonMessage(w, errorResponse{
			Error: "failed creating team",
			Team:  team,
		}, http.StatusInternalServerError)
		return
	}

	log.Info("Created new team")

	buffer := new(bytes.Buffer)
	if err = json.Compact(buffer, []byte(k8sconfig)); err != nil {
		a.log.Error("failed minifying kubeconfig", "error", err)
		writeJsonMessage(w, errorResponse{
			Error: "failed creating payload for team",
			Team:  team,
		}, http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = w.Write(buffer.Bytes())
}

func validateTeam(team string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9-]{2,63}$`).Match([]byte(team))
}

// Example: POST /team/{team}/next-task?task=int
func (a *api) teamNextTask(w http.ResponseWriter, r *http.Request) {
	type errorResponse struct {
		Error string `json:"error"`
		Team  string `json:"team"`
		Task  string `json:"task"`
	}

	team := r.PathValue("team")
	log := a.log.With("team", team)

	taskString := r.URL.Query().Get("task")
	if taskString == "" {
		a.log.Error("missing task query parameter")
		writeJsonMessage(w, errorResponse{
			Error: "missing task query parameter",
			Team:  team,
		}, http.StatusBadRequest)

		return
	}

	taskInt, err := strconv.Atoi(taskString)
	if err != nil {
		a.log.Error("task is not int", "error", err, "task", taskString)

		writeJsonMessage(w, errorResponse{
			Error: "can not parse task as int",
			Team:  team,
			Task:  taskString,
		}, http.StatusBadRequest)

		return
	}

	if err := a.k8s.TeamNextTask(r.Context(), team, taskInt); err != "" {
		log.Error("failed storing next task", "error", err)
		writeJsonMessage(w, errorResponse{
			Error: "failed storing task",
			Team:  team,
			Task:  taskString,
		}, http.StatusInternalServerError)

		return
	}

	writeJsonMessage(w, map[string]any{
		"message": "Task was updated",
		"team":    team,
	}, http.StatusOK)
}
