package k8s

import (
	"log/slog"

	"k8s.io/client-go/kubernetes"
)

type Client struct {
	client   *kubernetes.Clientset
	log      *slog.Logger
	Endpoint string
	CA       string
}

func New(client *kubernetes.Clientset, log *slog.Logger, endpoint, ca string) Client {
	return Client{
		client:   client,
		log:      log,
		Endpoint: endpoint,
		CA:       ca,
	}
}
