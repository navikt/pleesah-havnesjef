package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceInfo struct {
	Name      string
	Type      v1.ServiceType
	ClusterIP string
	Ports     []v1.ServicePort
}

func (c Client) ServiceInfo(ctx context.Context, team, service string) (*ServiceInfo, error) {
	svc, err := c.client.CoreV1().Services(team).Get(ctx, service, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &ServiceInfo{
		Name:      svc.Name,
		Type:      svc.Spec.Type,
		ClusterIP: svc.Spec.ClusterIP,
		Ports:     svc.Spec.Ports,
	}, nil
}
