package k8s

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceInfo struct {
	Name      string               `json:"name"`
	Type      corev1.ServiceType   `json:"type"`
	ClusterIP string               `json:"clusterIP"`
	Ports     []corev1.ServicePort `json:"ports"`
}

func (c Client) FetchService(ctx context.Context, team, serviceName string) (*corev1.Service, error) {
	service, err := c.client.CoreV1().Services(team).Get(ctx, serviceName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return service, nil
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
