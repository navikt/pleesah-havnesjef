package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/navikt/pleesah-havnesjef/internal/k8s"
)

type api struct {
	k8s    k8s.Client
	log    *slog.Logger
	server *http.Server
}

func New(client k8s.Client, log *slog.Logger) api {
	a := api{
		k8s: client,
		log: log,
	}

	mux := http.NewServeMux()
	mux.Handle("/api/v1/team/", http.StripPrefix("/api/v1/team", a.TeamHandler()))
	mux.HandleFunc("GET /api/v1/teams", a.TreasureMapHandler)

	server := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	a.server = server
	return a
}

func (a api) Run() {
	a.log.Info("Running on :8080")
	if err := a.server.ListenAndServe(); err != nil {
		panic(err.Error())
	}
}

func writeJsonMessage(w http.ResponseWriter, blob map[string]any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(blob)
}
