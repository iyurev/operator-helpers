package main

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

const (
	confDir  string = ".kube"
	confName string = "config"
)

func localTestClient() (*kubernetes.Clientset, error) {
	homeDir := os.Getenv("HOME")
	kubeConfPath := fmt.Sprintf("%s/%s/%s", homeDir, confDir, confName)
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfPath)
	if err != nil {
		return &kubernetes.Clientset{}, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return &kubernetes.Clientset{}, err
	}
	return client, err
}

func main() {

}
