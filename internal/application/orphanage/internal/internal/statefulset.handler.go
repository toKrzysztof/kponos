package resourceHandler

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StatefulSetHandler handles finding references to Secrets and ConfigMaps in StatefulSet resources
type StatefulSetHandler struct {
	client.Client
}

// NewStatefulSetHandler creates a new StatefulSetHandler
func NewStatefulSetHandler(client client.Client) *StatefulSetHandler {
	return &StatefulSetHandler{
	Client: client,
	}
}

// FindReferences finds all StatefulSets that reference the given Secret or ConfigMap
func (h *StatefulSetHandler) FindReferences(ctx context.Context, c client.Client, secretName, configMapName string, namespace string) ([]client.Object, error) {
	var statefulSets []client.Object
	statefulSetList := &appsv1.StatefulSetList{}
	if err := c.List(ctx, statefulSetList, client.InNamespace(namespace)); err != nil {
		return statefulSets, err
	}
	
	// TODO: Filter statefulsets that reference the secret/configmap
	
	return statefulSets, nil
}

// GetResourceType returns the resource type this handler processes
func (h *StatefulSetHandler) GetResourceType() string {
	return "StatefulSet"
}

