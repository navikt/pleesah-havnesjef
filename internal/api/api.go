package api

import (
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
	mux.HandleFunc("GET /api/v1/status/", a.StatusHandler)
	mux.HandleFunc("POST /api/v1/team/", a.TeamHandler)

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
