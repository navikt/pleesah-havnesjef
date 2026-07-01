package k8s

import (
	"context"
	"fmt"

	authenticationv1 "k8s.io/api/authentication/v1"
	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c Client) SetupTeam(ctx context.Context, teamName string) (string, error) {
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: teamName,
		},
	}

	_, err := c.client.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return "", err
	}

	serviceAccount := &apiv1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: teamName,
		},
	}

	_, err = c.client.CoreV1().ServiceAccounts(namespace.Name).Create(ctx, serviceAccount, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return "", err
	}

	oneDay := int64(86400)
	tokenRequest := &authenticationv1.TokenRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: teamName,
		},
		Spec: authenticationv1.TokenRequestSpec{
			ExpirationSeconds: &oneDay,
		},
	}

	token, err := c.client.CoreV1().ServiceAccounts(namespace.Name).CreateToken(ctx, serviceAccount.Name, tokenRequest, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	secret := apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "koordinatene-mine",
		},
		Data: map[string][]byte{
			"KOORDINATER": []byte("59.9124° N, 10.7962° E"),
		},
	}

	_, err = c.client.CoreV1().Secrets(namespace.Name).Create(ctx, &secret, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return "", err
	}

	roleBinding := rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: teamName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:     "Group",
				APIGroup: "rbac.authorization.k8s.io",
				Name:     fmt.Sprintf("system:serviceaccounts:%s", teamName),
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			APIGroup: "rbac.authorization.k8s.io",
			Name:     "pleesah-player",
		},
	}
	_, err = c.client.RbacV1().RoleBindings(namespace.Name).Create(ctx, &roleBinding, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return "", err
	}

	return createKubeconfig(teamName, token.Status.Token, c.Endpoint, c.CA), nil
}
