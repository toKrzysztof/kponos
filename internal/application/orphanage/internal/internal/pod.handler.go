package resourceHandler

import (
	"context"
	"fmt"

	core "github.com/toKrzysztof/kponos/internal/core/reference_analyzer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PodHandler handles finding references to Secrets and ConfigMaps in Pod resources
type PodHandler struct {
	client.Client
	referenceAnalyzer *core.ReferenceAnalyzer
	finders           map[string]ResourceReferenceFinder
}

// NewPodHandler creates a new PodHandler
func NewPodHandler(client client.Client) *PodHandler {
	analyzer := core.NewReferenceAnalyzer(client)
	h := &PodHandler{
		Client:            client,
		referenceAnalyzer: analyzer,
	}
	
	h.finders = map[string]ResourceReferenceFinder{
		"Secret":    h.findSecretReferences,
		"ConfigMap": h.findConfigMapReferences,
	}
	
	return h
}

// FindReferences finds all Pods that reference the given resource
func (h *PodHandler) FindReferences(ctx context.Context, c client.Client, resource client.Object, namespace string) ([]client.Object, error) {
	resourceKind := resource.GetObjectKind().GroupVersionKind().Kind
	resourceName := resource.GetName()
	
	finder, exists := h.finders[resourceKind]
	if !exists {
		return nil, fmt.Errorf("unsupported resource type: %s", resourceKind)
	}
	
	return finder(ctx, resourceName, namespace)
}

// GetResourceType returns the resource type this handler processes
func (h *PodHandler) GetResourceType() string {
	return "Pod"
}

// findSecretReferences finds all Pods that reference the given Secret
func (h *PodHandler) findSecretReferences(ctx context.Context, resourceName, namespace string) ([]client.Object, error) {
	return h.referenceAnalyzer.FindReferencesForSecret(ctx, resourceName, namespace, "Pod")
}

// findConfigMapReferences finds all Pods that reference the given ConfigMap
func (h *PodHandler) findConfigMapReferences(ctx context.Context, resourceName, namespace string) ([]client.Object, error) {
	return h.referenceAnalyzer.FindReferencesForConfigMap(ctx, resourceName, namespace, "Pod")
}

