package api

import (
	"log/slog"
	"net/http"

	"github.com/navikt/pleesah-havnesjef/internal/k8s"
)

func New(client k8s.Client, log *slog.Logger) api {
	a := api{
		k8s: client,
		log: log,
	}

	http.HandleFunc("GET /api/v1/service/status", a.ServiceRunningHandler)
	http.HandleFunc("POST /api/v1/team/", a.TeamHandler)

	return a
}

func (a api) Run() {
	a.log.Info("Running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err.Error())
	}
}

type api struct {
	k8s k8s.Client
	log *slog.Logger
}
