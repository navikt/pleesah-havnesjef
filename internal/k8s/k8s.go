package k8s

import (
	"log/slog"

	"k8s.io/client-go/kubernetes"
)

type Client struct {
	client *kubernetes.Clientset
	log    *slog.Logger
}

func New(client *kubernetes.Clientset, log *slog.Logger) Client {
	return Client{
		client: client,
		log:    log,
	}
}
