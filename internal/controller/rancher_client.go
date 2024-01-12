package controller

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	platformv1alpha1 "tmuntaner/platform-controller/api/v1alpha1"
)

type RancherClient struct {
	client   client.Client
	ranchers map[string]*KubernetesClient
}

func NewRancherClient(client client.Client) *RancherClient {
	return &RancherClient{
		client:   client,
		ranchers: make(map[string]*KubernetesClient),
	}
}

func (r *RancherClient) GetKubeConfigSecret(ctx context.Context, req ctrl.Request, instance *platformv1alpha1.UpstreamRancher) (*KubernetesClient, error) {
	if r.ranchers[instance.Spec.Name] != nil {
		return r.ranchers[instance.Spec.Name], nil
	}

	secretRef := instance.Spec.KubeConfig.ValueFrom.SecretKeyRef
	secretNamespacedName := types.NamespacedName{
		Name:      secretRef.Name,
		Namespace: req.Namespace,
	}
	secret := &v1.Secret{}
	if err := r.client.Get(ctx, secretNamespacedName, secret); err != nil {
		return nil, client.IgnoreNotFound(err)
	}

	kubeConfig := secret.Data[secretRef.Key]
	f, err := os.CreateTemp("", fmt.Sprintf("%s-kubeconfig.yml", instance.Spec.Name))
	if err != nil {
		return nil, err
	}

	_, err = f.Write(kubeConfig)
	if err != nil {
		return nil, err
	}

	clientset, err := NewKubernetesClient(f.Name())
	if err != nil {
		return nil, err
	}
	r.ranchers[instance.Spec.Name] = clientset

	return r.ranchers[instance.Spec.Name], nil
}
