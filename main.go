package main

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/navikt/pleesah-havnesjef/internal/api"
	"github.com/navikt/pleesah-havnesjef/internal/k8s"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	var kubeconfig *string
	kubeconfigEnv := os.Getenv("KUBECONFIG")
	if kubeconfigEnv != "" {
		log.Info("Using config from env")
		kubeconfig = &kubeconfigEnv
	} else if home := homedir.HomeDir(); home != "" {
		log.Info("Using config from .kube")
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		log.Info("Using config from flag")
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

	api := api.New(k8s.New(clientset, log.WithGroup("k8s")), log.WithGroup("api"))

	api.Run()
}
