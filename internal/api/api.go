package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/navikt/pleesah-havnesjef/internal/k8s"
	"github.com/navikt/pleesah-havnesjef/internal/workers"
)

type api struct {
	k8s     k8s.Client
	log     *slog.Logger
	server  *http.Server
	workers workers.Worker
}

func New(client k8s.Client, log *slog.Logger) api {
	a := api{
		k8s:     client,
		log:     log,
		workers: workers.New(client, log.WithGroup("workers")),
	}

	mux := http.NewServeMux()
	mux.Handle("/team/", http.StripPrefix("/team", a.TeamHandler()))
	mux.HandleFunc("GET /teams", a.TreasureMapHandler)

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

func writeJsonMessage(w http.ResponseWriter, blob any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(blob)
}
