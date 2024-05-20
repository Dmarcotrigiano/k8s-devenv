package kubernetes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubeConfig() (*rest.Config, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n", err)
		return nil, err
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("error getting Kubernetes config: %v\n", err)
		return nil, err
	}

	return kubeConfig, nil
}

func GetScyllaLoadBalancerAddrs(clientset k8s.Interface, namespace string) ([]string, error) {
	services, err := clientset.CoreV1().Services(namespace).List(context.Background(), v1meta.ListOptions{})
	if err != nil {
		return nil, err
	}
	var servicesExternalIps []string
	for _, service := range services.Items {
		if service.Spec.Type == v1.ServiceTypeLoadBalancer {
			servicesExternalIps = append(servicesExternalIps, service.Status.LoadBalancer.Ingress[0].IP)
		}
	}
	return servicesExternalIps, nil
}
