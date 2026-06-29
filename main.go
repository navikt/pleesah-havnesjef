package main

import (
	"fmt"
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

	endpoint := os.Getenv("ENDPOINT")
	if endpoint == "" {
		panic(fmt.Errorf("endpoint is not set\n"))
	}

	ca := os.Getenv("CA")
	if ca == "" {
		panic(fmt.Errorf("ca is not set\n"))
	}

	kubeconfig := findKubeconfig(log, endpoint, ca)

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	api := api.New(k8s.New(clientset, log.WithGroup("k8s"), endpoint, ca), log.WithGroup("api"))

	api.Run()
}

func findKubeconfig(log *slog.Logger, endpoint, ca string) string {
	var kubeconfig string
	kubeconfigToken := os.Getenv("KUBECONFIG_TOKEN")
	if kubeconfigToken != "" {
		log.Info("KUBECONFIG_TOKEN is set, creating kubeconfig")
		var err error
		kubeconfig, err = k8s.CreateHavnesjefConfig(kubeconfigToken, endpoint, ca)
		if err != nil {
			panic(fmt.Errorf("failed creating kubeconfig: %s", err))
		}
	} else {
		kubeconfigEnv := os.Getenv("KUBECONFIG")
		if kubeconfigEnv != "" {
			log.Info("Using config from env")
			kubeconfig = kubeconfigEnv
		} else if home := homedir.HomeDir(); home != "" {
			log.Info("Using config from .kube")
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
	}

	return kubeconfig
}
