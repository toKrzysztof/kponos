package internal

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PodHandler handles finding references to Secrets and ConfigMaps in Pod resources
type PodHandler struct {
	client.Client
}

// NewPodHandler creates a new PodHandler
func NewPodHandler(client client.Client) *PodHandler {
	return &PodHandler{
		Client: client,
	}
}

// FindReferences finds all Pods that reference the given Secret or ConfigMap
func (h *PodHandler) FindReferences(ctx context.Context, c client.Client, secretName, configMapName string, namespace string) ([]client.Object, error) {
	var pods []client.Object
	podList := &corev1.PodList{}
	if err := c.List(ctx, podList, client.InNamespace(namespace)); err != nil {
		return pods, err
	}
	
	// TODO: Filter pods that reference the secret/configmap
	
	return pods, nil
}

// GetResourceType returns the resource type this handler processes
func (h *PodHandler) GetResourceType() string {
	return "Pod"
}

