package main

import (
	"context"
	"fmt"
	k8sres "github.com/kube-operators/operator-helpers/pkg/k8s-resources"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
)

const (
	confDir  string = ".kube"
	confName string = "config"

	testSAName        = "test"
	testClusterRole   = "admin"
	testClusterRBName = "default-admin"
	testSaName        = "test"
	testNamespace     = "test"
)

func localTestClient() (*kubernetes.Clientset, error) {
	homeDir := os.Getenv("HOME")
	kubeConfPath := fmt.Sprintf("%s/%s/%s", homeDir, confDir, confName)
	fmt.Printf("Kubernetes config file path: %s\n", kubeConfPath)
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

func TestNS(client *kubernetes.Clientset, context context.Context) {
	ns, err := k8sres.CommonNamespace(testNamespace)
	if err != nil {
		log.Fatal(err)
	}
	existNS, err := client.CoreV1().Namespaces().List(context, metav1.ListOptions{FieldSelector: fmt.Sprintf("metadata.name=%s", ns.Name)})
	if err != nil {
		log.Fatal(err)
	}
	if len(existNS.Items) != 1 {
		_, err = client.CoreV1().Namespaces().Create(context, &ns, metav1.CreateOptions{})
		if err != nil {
			log.Fatal("%s\n", err)
		}
	}
	_, err = client.CoreV1().Namespaces().Update(context, &ns, metav1.UpdateOptions{})
	if err != nil {
		log.Fatal(err)
	}
}

func TestSA(client *kubernetes.Clientset, context context.Context) {
	sa, err := k8sres.CommonServiceAccount(testSAName)
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.CoreV1().ServiceAccounts(testNamespace).Create(context, &sa, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
}

func TestSAClusterRoleBinding(client *kubernetes.Clientset, context context.Context) {
	rb := k8sres.SaClusterRoleBindingToNamespace(testNamespace, testClusterRBName, testSAName, testClusterRole)
	_, err := client.RbacV1().ClusterRoleBindings().Create(context, &rb, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	client, err := localTestClient()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.TODO()
	TestNS(client, ctx)
}
