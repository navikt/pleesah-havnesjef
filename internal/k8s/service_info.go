package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceInfo struct {
	Name      string           `json:"name"`
	Type      v1.ServiceType   `json:"type"`
	ClusterIP string           `json:"clusterIP"`
	Ports     []v1.ServicePort `json:"ports"`
}

func (c Client) ServiceInfo(ctx context.Context, team string) (*ServiceInfo, error) {
	service, err := c.client.CoreV1().Services(team).Get(ctx, "service", metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &ServiceInfo{
		Name:      service.Name,
		Type:      service.Spec.Type,
		ClusterIP: service.Spec.ClusterIP,
		Ports:     service.Spec.Ports,
	}, nil
}
