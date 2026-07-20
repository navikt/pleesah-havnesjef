package k8s

import (
	"context"
	"time"

	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodInfo struct {
	Name     string
	Phase    v1.PodPhase
	Restarts int32
	Node     string
	Age      time.Duration
}

func (c Client) PodInfo(ctx context.Context, team string) (*PodInfo, error) {
	pod, err := c.client.CoreV1().Pods(team).Get(ctx, "pod", metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	var restarts int32
	for _, cs := range pod.Status.ContainerStatuses {
		restarts += cs.RestartCount
	}

	return &PodInfo{
		Name:     pod.Name,
		Phase:    pod.Status.Phase,
		Restarts: restarts,
		Node:     pod.Spec.NodeName,
		Age:      time.Since(pod.CreationTimestamp.Time),
	}, nil
}
