package internal

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DaemonSetHandler handles finding references to Secrets and ConfigMaps in DaemonSet resources
type DaemonSetHandler struct {
	client.Client
}

// NewDaemonSetHandler creates a new DaemonSetHandler
func NewDaemonSetHandler(client client.Client) *DaemonSetHandler {
	return &DaemonSetHandler{
		Client: client,
	}
}

// FindReferences finds all DaemonSets that reference the given Secret or ConfigMap
func (h *DaemonSetHandler) FindReferences(ctx context.Context, c client.Client, secretName, configMapName string, namespace string) ([]client.Object, error) {

	var daemonSets []client.Object
	daemonSetList := &appsv1.DaemonSetList{}
	if err := c.List(ctx, daemonSetList, client.InNamespace(namespace)); err != nil {
		return daemonSets, err
	}
	
	// TODO: Filter daemonsets that reference the secret/configmap
	
	return daemonSets, nil
}

// GetResourceType returns the resource type this handler processes
func (h *DaemonSetHandler) GetResourceType() string {
	return "DaemonSet"
}

