package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c Client) IsDeploymentRunning(ctx context.Context, team, name string) (bool, error) {
	_, err := c.client.AppsV1().Deployments(team).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (c Client) IsPodRunning(ctx context.Context, team, service string, isReady bool) (bool, error) {
	pod, err := c.client.CoreV1().Pods(team).Get(ctx, service, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	if isReady {
		return isPodReady(pod), nil
	}

	return true, nil
}

func isPodReady(pod *v1.Pod) bool {
	return isPodReadyConditionTrue(pod.Status)
}

func isPodReadyConditionTrue(status v1.PodStatus) bool {
	condition := getPodReadyCondition(status)
	return condition != nil && condition.Status == v1.ConditionTrue
}

func getPodReadyCondition(status v1.PodStatus) *v1.PodCondition {
	for i := range status.Conditions {
		if status.Conditions[i].Type == v1.PodReady {
			return &status.Conditions[i]
		}
	}
	return nil
}

func (c Client) IsServiceRunning(ctx context.Context, team, serviceName string) (bool, error) {
	service, err := c.FetchService(ctx, team, serviceName)
	if err != nil {
		return false, err
	}

	return service != nil, nil
}
