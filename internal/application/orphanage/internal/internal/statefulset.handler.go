package resourceHandler

import (
	"context"
	"fmt"

	core "github.com/toKrzysztof/kponos/internal/core/reference_analyzer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StatefulSetHandler handles finding references to Secrets and ConfigMaps in StatefulSet resources
type StatefulSetHandler struct {
	client.Client
	referenceAnalyzer *core.ReferenceAnalyzer
	finders           map[string]ResourceReferenceFinder
}

// NewStatefulSetHandler creates a new StatefulSetHandler
func NewStatefulSetHandler(client client.Client) *StatefulSetHandler {
	analyzer := core.NewReferenceAnalyzer(client)
	h := &StatefulSetHandler{
		Client:            client,
		referenceAnalyzer: analyzer,
	}
	
	h.finders = map[string]ResourceReferenceFinder{
		"Secret":    h.findSecretReferences,
		"ConfigMap": h.findConfigMapReferences,
	}
	
	return h
}

// FindReferences finds all StatefulSets that reference the given resource
func (h *StatefulSetHandler) FindReferences(ctx context.Context, c client.Client, resource client.Object, namespace string) ([]client.Object, error) {
	resourceKind := resource.GetObjectKind().GroupVersionKind().Kind
	resourceName := resource.GetName()
	
	finder, exists := h.finders[resourceKind]
	if !exists {
		return nil, fmt.Errorf("unsupported resource type: %s", resourceKind)
	}
	
	return finder(ctx, resourceName, namespace)
}

// GetResourceType returns the resource type this handler processes
func (h *StatefulSetHandler) GetResourceType() string {
	return "StatefulSet"
}

// findSecretReferences finds all StatefulSets that reference the given Secret
func (h *StatefulSetHandler) findSecretReferences(ctx context.Context, resourceName, namespace string) ([]client.Object, error) {
	return h.referenceAnalyzer.FindReferencesForSecret(ctx, resourceName, namespace, "StatefulSet")
}

// findConfigMapReferences finds all StatefulSets that reference the given ConfigMap
func (h *StatefulSetHandler) findConfigMapReferences(ctx context.Context, resourceName, namespace string) ([]client.Object, error) {
	return h.referenceAnalyzer.FindReferencesForConfigMap(ctx, resourceName, namespace, "StatefulSet")
}

