package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/navikt/pleesah-havnesjef/internal/k8s"
)

func (a *api) TeamHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /{team}/create", a.teamCreate)
	mux.HandleFunc("POST /{team}/next-task", a.teamNextTask)
	mux.HandleFunc("POST /{team}/progression", a.teamAddProgression)
	mux.HandleFunc("GET /{team}/status/{resource}", a.teamResourceStatus)
	mux.HandleFunc("GET /{team}/status", a.teamStatus)

	return mux
}

// Example: POST /team/{team}/create?hex={code}
func (a *api) teamCreate(w http.ResponseWriter, r *http.Request) {
	type errorResponse struct {
		Error   string `json:"error"`
		Team    string `json:"team"`
		Hexcode string `json:"hexcode"`
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

	hexcode := r.URL.Query().Get("hex")
	if !validateHexcode(hexcode) {
		a.log.Error("hex is not valid", "hex", hexcode)
		writeJsonMessage(w, errorResponse{
			Error: "hex is not valid",
			Team:  team,
		}, http.StatusBadRequest)

		return
	}

	k8sconfig, err := a.k8s.SetupTeam(r.Context(), team, hexcode)
	if err != nil {
		log.Error("failed creating team", "error", err)
		writeJsonMessage(w, errorResponse{
			Error:   "failed creating team",
			Team:    team,
			Hexcode: hexcode,
		}, http.StatusInternalServerError)
		return
	}

	log.Info("Created new team")

	go a.workers.StartSecretSignal(team)

	buffer := new(bytes.Buffer)
	if err = json.Compact(buffer, []byte(k8sconfig)); err != nil {
		a.log.Error("failed minifying kubeconfig", "error", err)
		writeJsonMessage(w, errorResponse{
			Error:   "failed creating payload for team",
			Team:    team,
			Hexcode: hexcode,
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

// Example: POST /team/{team}/progression
// Payload: {x: 0, y: 1}
func (a *api) teamAddProgression(w http.ResponseWriter, r *http.Request) {
	type errorResponse struct {
		Error       string          `json:"error"`
		Team        string          `json:"team"`
		Progression k8s.Progression `json:"progression"`
	}

	team := r.PathValue("team")
	log := a.log.With("team", team)

	var progression k8s.Progression
	if err := json.NewDecoder(r.Body).Decode(&progression); err != nil {
		log.Error("failed parsing body", "error", err)
		writeJsonMessage(w, errorResponse{
			Error: "failed parsing progression",
			Team:  team,
		}, http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	if err := a.k8s.TeamAddProgression(r.Context(), team, progression); err != "" {
		writeJsonMessage(w, errorResponse{
			Error:       err,
			Team:        team,
			Progression: progression,
		}, http.StatusInternalServerError)

		return
	}

	writeJsonMessage(w, map[string]any{
		"message":     "Progression was added",
		"team":        team,
		"progression": progression,
	}, http.StatusOK)
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
