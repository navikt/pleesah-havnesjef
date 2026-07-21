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

func (c Client) ServiceInfo(ctx context.Context, team string) ([]ServiceInfo, error) {
	services, err := c.client.CoreV1().Services(team).List(ctx, metav1.ListOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	var list []ServiceInfo
	for _, svc := range services.Items {
		list = append(list, ServiceInfo{
			Name:      svc.Name,
			Type:      svc.Spec.Type,
			ClusterIP: svc.Spec.ClusterIP,
			Ports:     svc.Spec.Ports,
		})
	}

	return list, nil
}
