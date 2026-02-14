package resourceHandler

import (
	"context"
	"fmt"

	core "github.com/toKrzysztof/kponos/internal/core/reference_analyzer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DaemonSetHandler handles finding references to Secrets and ConfigMaps in DaemonSet resources
type DaemonSetHandler struct {
	client.Client
	referenceAnalyzer *core.ReferenceAnalyzer
	finders           map[string]ResourceReferenceFinder
}

// NewDaemonSetHandler creates a new DaemonSetHandler
func NewDaemonSetHandler(client client.Client) *DaemonSetHandler {
	analyzer := core.NewReferenceAnalyzer(client)
	h := &DaemonSetHandler{
		Client:            client,
		referenceAnalyzer: analyzer,
	}
	
	h.finders = map[string]ResourceReferenceFinder{
		"Secret":    h.findSecretReferences,
		"ConfigMap": h.findConfigMapReferences,
	}
	
	return h
}

// FindReferences finds all DaemonSets that reference the given resource
func (h *DaemonSetHandler) FindReferences(ctx context.Context, c client.Client, resource client.Object, namespace string) ([]client.Object, error) {
	resourceKind := resource.GetObjectKind().GroupVersionKind().Kind
	resourceName := resource.GetName()
	
	finder, exists := h.finders[resourceKind]
	if !exists {
		return nil, fmt.Errorf("unsupported resource type: %s", resourceKind)
	}
	
	return finder(ctx, resourceName, namespace)
}

// GetResourceType returns the resource type this handler processes
func (h *DaemonSetHandler) GetResourceType() string {
	return "DaemonSet"
}

// findSecretReferences finds all DaemonSets that reference the given Secret
func (h *DaemonSetHandler) findSecretReferences(ctx context.Context, resourceName, namespace string) ([]client.Object, error) {
	return h.referenceAnalyzer.FindReferencesForSecret(ctx, resourceName, namespace, "DaemonSet")
}

// findConfigMapReferences finds all DaemonSets that reference the given ConfigMap
func (h *DaemonSetHandler) findConfigMapReferences(ctx context.Context, resourceName, namespace string) ([]client.Object, error) {
	return h.referenceAnalyzer.FindReferencesForConfigMap(ctx, resourceName, namespace, "DaemonSet")
}

