package api

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/navikt/pleesah-havnesjef/internal/k8s"
)

func New(client k8s.Client, log *slog.Logger) api {
	a := api{
		k8s: client,
		log: log,
	}

	http.HandleFunc("GET /", a.IndexHandler)
	http.HandleFunc("POST /", a.TeamHandler)

	return a
}

func (a api) Run() {
	a.log.Info("Running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err.Error())
	}
}

type api struct {
	k8s k8s.Client
	log *slog.Logger
}

func (a api) IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("form").Parse(indexTemplate))
	tmpl.Execute(w, nil)
}

func (a api) TeamHandler(w http.ResponseWriter, r *http.Request) {
	teamName := r.FormValue("team")
	k8sconfig, err := a.k8s.SetupTeam(r.Context(), teamName)
	if err != nil {
		a.log.Error("failed creating team", "error", err)
		tmpl := template.Must(template.New("kube").Parse(errorTemplate))
		tmpl.Execute(w, map[string]string{
			"Error": err.Error(),
		})

		return
	}

	tmpl := template.Must(template.New("kube").Parse(postTemplate))
	tmpl.Execute(w, map[string]string{
		"Kubeconfig": k8sconfig,
		"Team":       teamName,
	})

	a.log.Info("Created new team", "team", teamName)
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
    <style>
        code {
            border-radius: .5rem;
            overflow-x: auto;
            padding: .25rem;
            background-color: #24292e;
            color: #e1e4e8;
        }

        pre {
            border-radius: .5rem;
            overflow-x: auto;
            padding: 1rem;
            background-color: #24292e;
            color: #e1e4e8;"
        }
    </style>
</head>
<body>
    <h1>Kubeconfig for ✨{{ .Team }}✨</h1>
    <p>
        <ol>
            <li>Opprett en fil som heter <code>config</code></li>
            <li>Lim innholdet nedenfor inn i filen</li>
            <li>Kjør <code>export KUBECONFIG=./config</code> i din terminal</li>
        </ol>

        PS: Hvis du lukker terminalen din må du kjøre <code>export KUBECONFIG=./config</code> på nytt.
    </p>
    <button onclick="copyToClipboard()">Copy Kubeconfig</button>
    <pre id="kubeconfig">{{ .Kubeconfig }}</pre>
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

const errorTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Error!</title>
</head>
<body>
    <h1>Skipet ditt sank</h1>
    <p>{{ .Error }}</p>
</body>
</html>
`
