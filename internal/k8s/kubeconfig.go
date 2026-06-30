package k8s

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func CreateHavnesjefConfig(token, endpoint, ca string) (string, error) {
	kubeconfig := createKubeconfig("havnesjef", token, endpoint, ca)
	path := filepath.Join(os.TempDir(), ".config")
	err := os.WriteFile(path, []byte(kubeconfig), 0o600)

	return path, err
}

func createKubeconfig(team, token, endpoint, ca string) string {
	var sb strings.Builder
	tmpl := template.Must(template.New("kubeconfig").Parse(kubeconfigTemplate))
	_ = tmpl.Execute(&sb, map[string]string{
		"Name":     team,
		"Token":    token,
		"Endpoint": endpoint,
		"CA":       ca,
	})

	return sb.String()
}

const kubeconfigTemplate = `apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: {{ .CA }}
    server: https://{{ .Endpoint }}
  name: pleesah
contexts:
- context:
    cluster: pleesah
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
