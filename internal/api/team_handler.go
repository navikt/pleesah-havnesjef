package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (a *api) TeamHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /{team}/create", a.teamCreate)
	mux.HandleFunc("POST /{team}/next-task", a.teamNextTask)
	mux.HandleFunc("PUT /{team}/coordinates", a.teamAddCoordinates)

	return mux
}

// Example: POST /api/v1/team/{team}/create?hex={code}
func (a *api) teamCreate(w http.ResponseWriter, r *http.Request) {
	team := r.PathValue("team")
	k8sconfig, err := a.k8s.SetupTeam(r.Context(), team)
	if err != nil {
		a.log.Error("failed creating team", "error", err, "team", team)
		writeJsonMessage(w, map[string]any{
			"error": "failed creating team",
		}, http.StatusInternalServerError)

		return
	}

	a.log.Info("Created new team", "team", team)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	buffer := new(bytes.Buffer)
	err = json.Compact(buffer, []byte(k8sconfig))
	w.Write(buffer.Bytes())
}

// Example: PUT /api/v1/team/{team}/coordinates
// Payload: {x: 0, y: 1}
func (a *api) teamAddCoordinates(w http.ResponseWriter, r *http.Request) {
	team := r.PathValue("team")
	type Coordinates struct {
		X int
		Y int
	}

	var coordinates Coordinates
	if err := json.NewDecoder(r.Body).Decode(&coordinates); err != nil {
		a.log.Error("failed parsing body", "team", team, "error", err)
		writeJsonMessage(w, map[string]any{
			"error": "failed parsing body",
		}, http.StatusBadRequest)

		return
	}
	defer r.Body.Close()
	minifiedCoordinates := fmt.Sprintf("%d,%d", coordinates.X, coordinates.Y)

	if err := a.k8s.TeamAddCoordinates(r.Context(), team, minifiedCoordinates); err != "" {
		writeJsonMessage(w, map[string]any{
			"error":       err,
			"team":        team,
			"coordinates": minifiedCoordinates,
		}, http.StatusInternalServerError)

		return
	}

	writeJsonMessage(w, map[string]any{
		"message":     "Coordinates was added",
		"team":        team,
		"coordinates": minifiedCoordinates,
	}, http.StatusOK)
}

// Example: POST /api/v1/team/{team}/next-task?task=int
func (a *api) teamNextTask(w http.ResponseWriter, r *http.Request) {
	team := r.PathValue("team")
	taskString := r.URL.Query().Get("task")
	if taskString == "" {
		a.log.Error("missing task query parameter", "team", team)
		writeJsonMessage(w, map[string]any{
			"error": "missing task query parameter",
		}, http.StatusBadRequest)

		return
	}

	taskInt, err := strconv.Atoi(taskString)
	if err != nil {
		a.log.Error("task is not int", "error", err, "team", team, "task", taskString)

		writeJsonMessage(w, map[string]any{
			"error": "can not parse task as int",
			"team":  team,
			"task":  taskString,
		}, http.StatusBadRequest)

		return
	}

	if err := a.k8s.TeamNextTask(r.Context(), team, taskInt); err != "" {
		writeJsonMessage(w, map[string]any{
			"error": err,
			"team":  team,
		}, http.StatusInternalServerError)

		return
	}

	writeJsonMessage(w, map[string]any{
		"message": "Task was updated",
		"team":    team,
	}, http.StatusOK)
}
