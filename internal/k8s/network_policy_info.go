package k8s

import (
	"context"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NetworkPolicyInfo struct {
	Name string `json:"name"`
}

func (c Client) NetworkPolicyInfo(ctx context.Context, team string) ([]NetworkPolicyInfo, error) {
	networkPolicies, err := c.client.NetworkingV1().NetworkPolicies(team).List(ctx, metav1.ListOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	var list []NetworkPolicyInfo
	for _, policy := range networkPolicies.Items {
		list = append(list, NetworkPolicyInfo{
			Name: policy.Name,
		})
	}

	return list, nil
}
