package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	authenticationv1 "k8s.io/api/authentication/v1"
	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	PLEESAH_TASK        = "pleesah.io/task"
	PLEESAH_HEXCODE     = "pleesah.io/hexcode"
	PLEESAH_COORDINATES = "pleesah.io/coordinates"
)

type Team struct {
	Name        string   `json:"navn"`
	Hexcode     string   `json:"hexKode"`
	Progression []string `json:"progresjon"`
}

func (c Client) TeamAddCoordinates(ctx context.Context, team, minifiedCoordinates string) string {
	log := c.log.With("team", team)
	namespace, err := c.getTeam(ctx, team)
	if err != nil {
		log.Error("failed fetching team", "error", err)
		return "team was not found"
	}

	coordinatesString := namespace.Annotations[PLEESAH_COORDINATES]
	var coordinates []string
	if err := json.Unmarshal([]byte(coordinatesString), &coordinates); err != nil {
		log.Error("failed unmarshaling coordinates", "error", err, "coordinates", coordinatesString)
		return "failed reading coordinates"
	}

	coordinates = append(coordinates, minifiedCoordinates)
	payload, err := json.Marshal(coordinates)
	if err != nil {
		log.Error("failed marshaling coordinates", "error", err, "coordinates", coordinates[len(coordinates)-1])
		return "failed writing coordinates"
	}

	namespace.Annotations[PLEESAH_COORDINATES] = string(payload)
	if err := c.UpdateTeam(ctx, namespace); err != nil {
		c.log.Error("failed storing team", "error", err)
	}

	return ""
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
		return fmt.Sprintf("failed parsing old task as int: %s", &oldTaskString)
	}

	if task <= oldTaskInt {
		return "task was lower than previous task"
	}

	namespace.Annotations[PLEESAH_TASK] = fmt.Sprint(task)
	if err := c.UpdateTeam(ctx, namespace); err != nil {
		c.log.Error("failed updating with new task", "error", err, "team", team, "task", task)
		//=http.Error(w, , http.StatusInternalServerError)
		return "failed updating with new task"
	}

	return ""
}

func (c Client) getTeam(ctx context.Context, teamName string) (*apiv1.Namespace, error) {
	return c.client.CoreV1().Namespaces().Get(ctx, teamName, metav1.GetOptions{})
}

func (c Client) SetupTeam(ctx context.Context, teamName string) (string, error) {
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: teamName,
			Annotations: map[string]string{
				PLEESAH_TASK:        "0",
				PLEESAH_HEXCODE:     "#123321",
				PLEESAH_COORDINATES: "[]",
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

func (c Client) ListTeams(ctx context.Context) ([]Team, error) {
	namespaces, err := c.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{
		LabelSelector: "player=true",
	})
	if err != nil {
		return nil, err
	}

	teams := make([]Team, len(namespaces.Items))
	for i, item := range namespaces.Items {
		team := namespaceToTeam(item)
		teams[i] = team
	}

	return teams, nil
}

func namespaceToTeam(namespace apiv1.Namespace) Team {
	annotations := namespace.GetAnnotations()
	var progression []string
	json.Unmarshal([]byte(annotations[PLEESAH_COORDINATES]), &progression)
	return Team{
		Name:        namespace.Name,
		Hexcode:     annotations[PLEESAH_HEXCODE],
		Progression: progression,
	}
}

func (c Client) UpdateTeam(ctx context.Context, namespace *apiv1.Namespace) error {
	_, err := c.client.CoreV1().Namespaces().Update(ctx, namespace, metav1.UpdateOptions{})
	return err
}
