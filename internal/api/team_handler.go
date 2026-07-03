package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

func (a *api) TeamHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /{team}/create", a.teamCreate)
	mux.HandleFunc("POST /{team}/next-task", a.teamNextTask)
	mux.HandleFunc("POST /{team}/progression", a.teamAddProgression)
	mux.HandleFunc("GET /{team}/status/{resource}", a.teamResourceStatus)

	return mux
}

// Example: POST /api/v1/team/{team}/create?hex={code}
func (a *api) teamCreate(w http.ResponseWriter, r *http.Request) {
	team := r.PathValue("team")
	log := a.log.With("team", team)

	if !validateTeam(team) {
		a.log.Error("team is not valid")
		writeJsonMessage(w, map[string]any{
			"error": "team is not valid",
		}, http.StatusBadRequest)

		return
	}

	hexcode := r.URL.Query().Get("hex")
	if !validateHexcode(hexcode) {
		a.log.Error("hex is not valid", "hex", hexcode)
		writeJsonMessage(w, map[string]any{
			"error": "hex is not valid",
		}, http.StatusBadRequest)

		return
	}

	k8sconfig, err := a.k8s.SetupTeam(r.Context(), team, hexcode)
	if err != nil {
		log.Error("failed creating team", "error", err)
		writeJsonMessage(w, map[string]any{
			"error":   "failed creating team",
			"team":    team,
			"hexcode": hexcode,
		}, http.StatusInternalServerError)

		return
	}

	log.Info("Created new team")

	buffer := new(bytes.Buffer)
	if err = json.Compact(buffer, []byte(k8sconfig)); err != nil {
		a.log.Error("failed minifying kubeconfig", "error", err)
		writeJsonMessage(w, map[string]any{
			"error": "",
			"team":  team,
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

func validateHexcode(hex string) bool {
	return regexp.MustCompile(`^#?(?:[a-fA-F0-9]{6}|[a-fA-F0-9]{3})$`).Match([]byte(hex))
}

// Example: POST /api/v1/team/{team}/progression
// Payload: {x: 0, y: 1}
func (a *api) teamAddProgression(w http.ResponseWriter, r *http.Request) {
	team := r.PathValue("team")
	log := a.log.With("team", team)
	type Progression struct {
		X float32
		Y float32
	}

	var progression Progression
	if err := json.NewDecoder(r.Body).Decode(&progression); err != nil {
		log.Error("failed parsing body", "error", err)
		writeJsonMessage(w, map[string]any{
			"error": "failed parsing progression",
		}, http.StatusBadRequest)

		return
	}
	defer r.Body.Close()
	minifiedProgression := fmt.Sprintf("%f,%f", progression.X, progression.Y)

	if err := a.k8s.TeamAddProgression(r.Context(), team, minifiedProgression); err != "" {
		writeJsonMessage(w, map[string]any{
			"error":       err,
			"team":        team,
			"progression": minifiedProgression,
		}, http.StatusInternalServerError)

		return
	}

	writeJsonMessage(w, map[string]any{
		"message":     "Progression was added",
		"team":        team,
		"progression": minifiedProgression,
	}, http.StatusOK)
}

// Example: POST /api/v1/team/{team}/next-task?task=int
func (a *api) teamNextTask(w http.ResponseWriter, r *http.Request) {
	team := r.PathValue("team")
	log := a.log.With("team", team)

	taskString := r.URL.Query().Get("task")
	if taskString == "" {
		a.log.Error("missing task query parameter")
		writeJsonMessage(w, map[string]any{
			"error": "missing task query parameter",
		}, http.StatusBadRequest)

		return
	}

	taskInt, err := strconv.Atoi(taskString)
	if err != nil {
		a.log.Error("task is not int", "error", err, "task", taskString)

		writeJsonMessage(w, map[string]any{
			"error": "can not parse task as int",
			"team":  team,
			"task":  taskString,
		}, http.StatusBadRequest)

		return
	}

	if err := a.k8s.TeamNextTask(r.Context(), team, taskInt); err != "" {
		log.Error("failed storing next task", "error", err)
		writeJsonMessage(w, map[string]any{
			"error": "failed storing task",
			"team":  team,
		}, http.StatusInternalServerError)

		return
	}

	writeJsonMessage(w, map[string]any{
		"message": "Task was updated",
		"team":    team,
	}, http.StatusOK)
}
