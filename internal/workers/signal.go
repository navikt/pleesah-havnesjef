package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (w Worker) StartSecretSignal(team string) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	ctx := context.Background()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.log.Info("tick")
			if err := w.sendSecretSignal(ctx, team); err != nil {
				w.log.Error("failed sending secret signal", "error", err)
			}
		}
	}
}

func (w Worker) sendSecretSignal(ctx context.Context, team string) error {
	service, err := w.k8s.FetchService(ctx, team, "tobias")
	if err != nil {
		return err
	}

	if service == nil {
		w.log.Info("no service")
		return nil
	}

	if len(service.Status.LoadBalancer.Ingress) != 1 {
		w.log.Info("not one ip", "ips", service.Spec.ExternalIPs)
		return nil
	}

	client := http.Client{
		Timeout: time.Second * 10,
	}

	externalIP := service.Status.LoadBalancer.Ingress[0].IP
	url := fmt.Sprintf("http://%s/notify", externalIP)
	payload, err := json.Marshal(map[string]string{
		"message": "59.9124° N, 10.7962° E",
	})
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	w.log.Info("message sent", "code", resp.StatusCode, "team", team)
	return nil
}
