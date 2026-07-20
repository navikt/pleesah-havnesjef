package k8s

import (
	"context"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentInfo struct {
	Name      string `json:"name"`
	Desired   int32  `json:"desired"`
	Ready     int32  `json:"ready"`
	Available int32  `json:"available"`
	Updated   int32  `json:"updated"`
}

func (c Client) DeploymentInfo(ctx context.Context, team string) (*DeploymentInfo, error) {
	deployment, err := c.client.AppsV1().Deployments(team).Get(ctx, "deployment", metav1.GetOptions{})

	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &DeploymentInfo{
		Name:      deployment.Name,
		Desired:   deployment.Status.Replicas,
		Ready:     deployment.Status.ReadyReplicas,
		Available: deployment.Status.AvailableReplicas,
		Updated:   deployment.Status.UpdatedReplicas,
	}, nil
}
