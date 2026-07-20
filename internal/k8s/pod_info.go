package k8s

import (
	"context"
	"time"

	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodInfo struct {
	Name     string        `json:"name"`
	Phase    v1.PodPhase   `json:"phase"`
	Restarts int32         `json:"restarts"`
	Node     string        `json:"node"`
	Age      time.Duration `json:"age"`
}

func (c Client) PodInfo(ctx context.Context, team string) ([]PodInfo, error) {
	pods, err := c.client.CoreV1().Pods(team).List(ctx, metav1.ListOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	var list []PodInfo
	for _, pod := range pods.Items {
		var restarts int32
		for _, cs := range pod.Status.ContainerStatuses {
			restarts += cs.RestartCount
		}
		list = append(list, PodInfo{
			Name:     pod.Name,
			Phase:    pod.Status.Phase,
			Restarts: restarts,
			Node:     pod.Spec.NodeName,
			Age:      time.Since(pod.CreationTimestamp.Time),
		})
	}

	return list, nil
}
