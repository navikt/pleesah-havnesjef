package k8s

import (
	"context"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c Client) IsPodRunning(ctx context.Context, team, service string) (bool, error) {
	_, err := c.client.CoreV1().Pods(team).Get(ctx, service, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (c Client) IsServiceRunning(ctx context.Context, team, service string) (bool, error) {
	_, err := c.client.CoreV1().Services(team).Get(ctx, service, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
