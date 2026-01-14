package main

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	authenticationv1 "k8s.io/api/authentication/v1"
	apiv1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const kubeconfigTemplate = `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVMRENDQXBTZ0F3SUJBZ0lRZkkzYzRQQ0tPc3ZhSkNhazRYUXRuekFOQmdrcWhraUc5dzBCQVFzRkFEQXYKTVMwd0t3WURWUVFERXlSa05UY3lZVGhpWXkxa1pUVXdMVFEyWmprdE9EWTROaTAxTnpka1ltTTFPR1JsWXpndwpJQmNOTWpZd01URXpNVEV4TkRNNFdoZ1BNakExTmpBeE1EWXhNakUwTXpoYU1DOHhMVEFyQmdOVkJBTVRKR1ExCk56SmhPR0pqTFdSbE5UQXRORFptT1MwNE5qZzJMVFUzTjJSaVl6VTRaR1ZqT0RDQ0FhSXdEUVlKS29aSWh2Y04KQVFFQkJRQURnZ0dQQURDQ0FZb0NnZ0dCQUp2TEJqbkxVdEttcEtWOTR3cGxlYXhUbkpYZmdxVHY4MmoyK0VpbwpuUEpibFpKdGdxbmJPSTlNaTFVRzQ1YmNCaFNudzFSeFdKSnVyNUUvbzUyNHRVamlWTlBXb1dDYlFpVE9mblNlCnhsdzRMbjFhd3dGVTYzRlNIajNMMGx1M0xhbnBiWWt3NU1NdlE0a2l4NkQvaWtJQmNUS1kzOXJ5TDdrMmVXbUIKK1pNNHFLWElzUDFXM0d4cGpndkgybGRtVE1DMWwwMVhERC9YNmdWK1hBVGRid2NHTFJrb1h4c2VOaDRWam1xSgpqUGF1T0VQeTRvTUtzSjVWbTZZQWxtcGlOLzBTUFdMUDFPZFpub1k0MWlwQlNCQllPRHAwRUNJSDVidWtaNi9hClhZaEpEU0ltemlCUENNNGNObVNaRFlvTHlRUlBuM1cwWFo5UHJId1BUTGtNaE1YdXRNeWJJVTlvcWdvNXpxTTUKWFJzVnJPQ1BLZDdhV0d3UXNoemN4MFc0Z3FpeHcrc1JHYVIyTm94YVcycFhEOWRtQXFVQ0N0YlFOMk01eXppbworRUlGY0VYZmNGRlJJSndFSzRmbXA3bzNIUnhUL3hpNXlPSWgyTDFpVGc5RW9vTnE2OHp6M1pUM2JkZUpuZzRGCnhGUzd5d0tGQkNXa3grZG1lc2s3WWVHOWNRSURBUUFCbzBJd1FEQU9CZ05WSFE4QkFmOEVCQU1DQWdRd0R3WUQKVlIwVEFRSC9CQVV3QXdFQi96QWRCZ05WSFE0RUZnUVU0WW1HZUpVYWhBcFIyV0t3b0dBbFNBcFNnMm93RFFZSgpLb1pJaHZjTkFRRUxCUUFEZ2dHQkFKRCtpZnFLL3dHRXRNMzM4dmJxSUx3WFBwcVRuNm0yTHhDU2owbVhDdXNHCmh4RjJnNnlsMW5EaU5DREVTcmY5a3NVSFFmczNBWng4cE95ak0vMjBPRzZEcllScmt5WVErTEVHem95bUtnd24KSkl4eUhIcGNmZHpHYzI3dXFnSEJ2VzdzQ04vWnFBcjZYUXMwdjhsdXdxd2pibG9TL1VKS3pCN1JOeHArbGVhYQpXSGoxVVFJYnNZZGREUWJFRlBEbk43djBVbVZzT0c2Ukhvd1JyQTRMSldsQmI5OTdweTRzQ0syOFBjR1BlYUEwCmd0UmpDeWN6RmtJR3ppcEE4Mjhab2p1R0VVck9zMnlxK3RYOWFQVGl3Q1E2NTBuS0o5eTVuc05IV09KYksyenMKeGxvbHJzY1ZIQ2ZOZVltZjFqVjR5aWVHK1I5TlYrNXVjWUxZdzdVOW9TZjhPWUFRdEFYNGYwdEJPQzYyZ1lFRgoxRmVsQmxobXlGOWNXcWtTYVFhK1k3RG52RXJFVmdTd01mSmd0WDkvRHpvcWh2VjdtME16R0VnZ1ppam1KV2xJCmZvN01aaVpwSUNOZHRmVmx1WW54N2VJbFFSaDAycVl5MWl1SUV3MnhabFZDTllZdWVodnUwaEs0b2MrdHZib1UKNWRBU2NqLzJkM3lMT0s5WVEzV21TZz09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K
    server: https://34.51.167.42
  name: gke_leesah-quiz-dev-5cf6_europe-north2_pleesah
contexts:
- context:
    cluster: gke_leesah-quiz-dev-5cf6_europe-north2_pleesah
    namespace: {{ .Name }}
    user: {{ .Name }}
  name: pleesah
current-context: pleesah
kind: Config
preferences: {}
users:
- name: {{ .Name }}
  user:
    token: {{ .Token }}
`

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	api := API{
		log:    log,
		client: clientset,
	}

	http.HandleFunc("GET /", api.IndexHandler)
	http.HandleFunc("POST /", api.TeamHandler)

	log.Info("Running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err.Error())
	}
}

type API struct {
	log    *slog.Logger
	client *kubernetes.Clientset
}

const indexTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Ut på bølgene blå</title>
</head>
<body>
    <h1>Lagnavn</h1>
    <form method="POST" action="/">
        <input type="text" name="team" placeholder="Lagnavn" required>
        <button type="submit">Submit</button>
    </form>
</body>
</html>
`

const postTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Kubeconfig for {{ .Team }}</title>
</head>
<body>
    <h1>Kubeconfig for {{ .Team }}</h1>
    <button onclick="copyToClipboard()">Copy Kubeconfig</button>
    <pre id="kubeconfig" style="border-radius: .5rem;overflow-x: auto;padding: 1rem;background-color: #24292e;color: #e1e4e8;">{{ .Kubeconfig }}</pre>
    <a href="/">Back</a>
    <script>
      function copyToClipboard() {
        const pre = document.getElementById('kubeconfig');
        const text = pre.innerText;
        navigator.clipboard.writeText(text);
      }
    </script>
</body>
</html>
`

func (a API) IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("form").Parse(indexTemplate))
	tmpl.Execute(w, nil)
}

func (a API) TeamHandler(w http.ResponseWriter, r *http.Request) {
	teamName := r.FormValue("team")
	k8sconfig, err := setupTeam(r.Context(), a.client, teamName)
	if err != nil {
		a.log.Error("failed creating team", "error", err)
	}

	tmpl := template.Must(template.New("kube").Parse(postTemplate))
	tmpl.Execute(w, map[string]string{
		"Kubeconfig": k8sconfig,
		"Team":       teamName,
	})
}

func setupTeam(ctx context.Context, client *kubernetes.Clientset, teamName string) (string, error) {
	namespace := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: teamName,
		},
	}

	_, err := client.CoreV1().Namespaces().Create(ctx, namespace, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	serviceAccount := &apiv1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: teamName,
		},
	}

	_, err = client.CoreV1().ServiceAccounts(namespace.Name).Create(ctx, serviceAccount, metav1.CreateOptions{})
	if err != nil {
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
	token, err := client.CoreV1().ServiceAccounts(namespace.Name).CreateToken(ctx, serviceAccount.Name, tokenRequest, metav1.CreateOptions{})
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

	_, err = client.CoreV1().Secrets(namespace.Name).Create(ctx, &secret, metav1.CreateOptions{})
	if err != nil {
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
	_, err = client.RbacV1().RoleBindings(namespace.Name).Create(ctx, &roleBinding, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	tmpl := template.Must(template.New("kubeconfig").Parse(kubeconfigTemplate))
	tmpl.Execute(&sb, map[string]string{
		"Name":  teamName,
		"Token": token.Status.Token,
	})

	return sb.String(), nil
}
