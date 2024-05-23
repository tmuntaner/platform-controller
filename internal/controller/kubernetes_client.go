package controller

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubernetesClient struct {
	*kubernetes.Clientset
}

func NewKubernetesClient(path string) (*KubernetesClient, error) {
	config, err := clientcmd.BuildConfigFromFlags("", path)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &KubernetesClient{
		Clientset: clientset,
	}, nil
}
