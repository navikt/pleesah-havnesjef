package k8s

import (
	"context"
	"fmt"
	"strconv"

	authenticationv1 "k8s.io/api/authentication/v1"
	apiv1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	PLEESAH_TASK = "pleesah.io/task"
)

type Team struct {
	Name        string        `json:"navn"`
	Hexcode     string        `json:"hexKode"`
	Progression []Progression `json:"progresjon"`
}

type Progression struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

func (c Client) TeamNextTask(ctx context.Context, team string, task int) string {
	namespace, err := c.getTeam(ctx, team)
	if err != nil {
		c.log.Error("failed fetching team", "error", err, "team", team)
		return "team was not found"
	}

	oldTaskString := namespace.Annotations[PLEESAH_TASK]
	oldTaskInt, err := strconv.Atoi(oldTaskString)
	if err != nil {
		c.log.Error("task is not int", "error", err, "team", team, "task", task)
		return fmt.Sprintf("failed parsing old task as int: %s", oldTaskString)
	}

	if task <= oldTaskInt {
		return "task was lower than, or equal to previous task"
	}

	namespace.Annotations[PLEESAH_TASK] = fmt.Sprint(task)
	if err := c.UpdateTeam(ctx, namespace); err != nil {
		c.log.Error("failed updating with new task", "error", err, "team", team, "task", task)
		return "failed updating with new task"
	}

	return ""
}

func (c Client) getTeam(ctx context.Context, teamName string) (*apiv1.Namespace, error) {
	return c.client.CoreV1().Namespaces().Get(ctx, teamName, metav1.GetOptions{})
}

func (c Client) SetupTeam(ctx context.Context, team string) (string, error) {
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: team,
			Annotations: map[string]string{
				PLEESAH_TASK: "0",
			},
			Labels: map[string]string{
				"player": "true",
			},
		},
	}

	_, err := c.client.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return "", err
	}

	serviceAccount := &apiv1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: team,
		},
	}

	_, err = c.client.CoreV1().ServiceAccounts(namespace.Name).Create(ctx, serviceAccount, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return "", err
	}

	oneDay := int64(86400)
	tokenRequest := &authenticationv1.TokenRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name: team,
		},
		Spec: authenticationv1.TokenRequestSpec{
			ExpirationSeconds: &oneDay,
		},
	}

	token, err := c.client.CoreV1().ServiceAccounts(namespace.Name).CreateToken(ctx, serviceAccount.Name, tokenRequest, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	roleBinding := rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: team,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:     "Group",
				APIGroup: "rbac.authorization.k8s.io",
				Name:     fmt.Sprintf("system:serviceaccounts:%s", team),
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

	denyAllNetworkPolicy := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: "deny-all",
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeIngress,
				networkingv1.PolicyTypeEgress,
			},
		},
	}
	_, err = c.client.NetworkingV1().NetworkPolicies(namespace.Name).Create(ctx, denyAllNetworkPolicy, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return "", err
	}

	AllowDNSNetworkPolicy := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: "allow-kube-dns",
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			Egress: []networkingv1.NetworkPolicyEgressRule{
				{
					Ports: []networkingv1.NetworkPolicyPort{
						{
							Protocol: new(apiv1.ProtocolUDP),
							Port: &intstr.IntOrString{
								IntVal: 53,
							},
						},
					},
				},
			},
			PolicyTypes: []networkingv1.PolicyType{
				networkingv1.PolicyTypeEgress,
			},
		},
	}
	_, err = c.client.NetworkingV1().NetworkPolicies(namespace.Name).Create(ctx, AllowDNSNetworkPolicy, metav1.CreateOptions{})
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return "", err
	}

	return createKubeconfig(team, token.Status.Token, c.Endpoint, c.CA), nil
}

func (c Client) ListTeams(ctx context.Context) ([]Team, error) {
	namespaces, err := c.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{
		LabelSelector: "player=true",
	})
	if err != nil {
		return nil, err
	}

	teams := make([]Team, len(namespaces.Items))
	for i, item := range namespaces.Items {
		teams[i] = Team{
			Name: item.Name,
		}
	}

	return teams, nil
}

func (c Client) UpdateTeam(ctx context.Context, namespace *apiv1.Namespace) error {
	_, err := c.client.CoreV1().Namespaces().Update(ctx, namespace, metav1.UpdateOptions{})
	return err
}
