package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/navikt/pleesah-havnesjef/internal/k8s"
)

func (a *api) TeamHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /{team}/create", a.teamCreate)
	mux.HandleFunc("POST /{team}/next-task", a.teamNextTask)

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

// Example: POST /api/v1/team/{team}/next-task?task=int
func (a *api) teamNextTask(w http.ResponseWriter, r *http.Request) {
	team := r.PathValue("team")
	namespace, err := a.k8s.GetTeam(r.Context(), team)
	if err != nil {
		a.log.Error("failed fetching team", "error", err, "team", team)
		http.Error(w, "team was not found", http.StatusNotFound)
		return
	}

	taskString := r.URL.Query().Get("task")
	if taskString == "" {
		a.log.Error("missing task query parameter", "team", team)
		http.Error(w, "missing task query parameter", http.StatusBadRequest)
		return
	}

	taskInt, err := strconv.Atoi(taskString)
	if err != nil {
		a.log.Error("task is not int", "error", err, "team", team, "task", taskString)
		http.Error(w, "can not parse task as int", http.StatusBadRequest)
		return
	}

	oldTaskString := namespace.Annotations[k8s.PLEESAH_TASK]
	oldTaskInt, err := strconv.Atoi(oldTaskString)
	if err != nil {
		a.log.Error("task is not int", "error", err, "team", team, "task", taskString)
		http.Error(w, "can not parse task as int", http.StatusBadRequest)
		return
	}

	if taskInt <= oldTaskInt {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message": "Task was lower than previous task",
			"team":    team,
		})

		return
	}

	namespace.Annotations[k8s.PLEESAH_TASK] = taskString
	if err := a.k8s.UpdateTeam(r.Context(), namespace); err != nil {
		a.log.Error("failed updating with new task", "error", err, "team", team, "task", taskString)
		http.Error(w, "failed updating with new task", http.StatusInternalServerError)
		return

	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": "Task was updated",
		"team":    team,
	})
}
