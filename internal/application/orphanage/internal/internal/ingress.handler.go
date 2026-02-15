package resourceHandler

import (
	"context"
	"fmt"

	core "github.com/toKrzysztof/kponos/internal/core/reference_analyzer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IngressHandler handles finding references to Secrets and ConfigMaps in Ingress resources
type IngressHandler struct {
	client.Client
	referenceAnalyzer *core.ReferenceAnalyzer
	finders           map[string]ResourceReferenceFinder
}

// NewIngressHandler creates a new IngressHandler
func NewIngressHandler(c client.Client) *IngressHandler {
	analyzer := core.NewReferenceAnalyzer(c)
	h := &IngressHandler{
		Client:            c,
		referenceAnalyzer: analyzer,
	}

	h.finders = map[string]ResourceReferenceFinder{
		"Secret":    h.findSecretReferences,
		"ConfigMap": h.findConfigMapReferences,
	}

	return h
}

// FindReferences finds all Ingresses that reference the given resource
func (h *IngressHandler) FindReferences(ctx context.Context, c client.Client, resource client.Object, namespace string) ([]client.Object, error) {
	resourceKind := resource.GetObjectKind().GroupVersionKind().Kind
	resourceName := resource.GetName()

	finder, exists := h.finders[resourceKind]
	if !exists {
		return nil, fmt.Errorf("unsupported resource type: %s", resourceKind)
	}

	return finder(ctx, resourceName, namespace)
}

// GetResourceType returns the resource type this handler processes
func (h *IngressHandler) GetResourceType() string {
	return "Ingress"
}

// findSecretReferences finds all Ingresses that reference the given Secret
func (h *IngressHandler) findSecretReferences(ctx context.Context, resourceName, namespace string) ([]client.Object, error) {
	return h.referenceAnalyzer.FindReferencesForSecret(ctx, resourceName, namespace, "Ingress")
}

func (h *IngressHandler) findConfigMapReferences(ctx context.Context, resourceName, namespace string) ([]client.Object, error) {
	return []client.Object{}, nil
}
