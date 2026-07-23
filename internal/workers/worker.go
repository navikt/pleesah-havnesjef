package workers

import (
	"log/slog"

	"github.com/navikt/pleesah-havnesjef/internal/k8s"
)

type Worker struct {
	log *slog.Logger
	k8s k8s.Client
}

func New(client k8s.Client, log *slog.Logger) Worker {
	return Worker{
		k8s: client,
		log: log,
	}
}
