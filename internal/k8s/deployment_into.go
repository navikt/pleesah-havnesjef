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

func (c Client) DeploymentInfo(ctx context.Context, team string) ([]DeploymentInfo, error) {
	deployments, err := c.client.AppsV1().Deployments(team).List(ctx, metav1.ListOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	var list []DeploymentInfo
	for _, deployment := range deployments.Items {
		list = append(list, DeploymentInfo{
			Name:      deployment.Name,
			Desired:   deployment.Status.Replicas,
			Ready:     deployment.Status.ReadyReplicas,
			Available: deployment.Status.AvailableReplicas,
			Updated:   deployment.Status.UpdatedReplicas,
		})
	}

	return list, nil
}
