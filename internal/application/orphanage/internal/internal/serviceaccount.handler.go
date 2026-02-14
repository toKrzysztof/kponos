package resourceHandler

import (
	"context"
	"fmt"

	core "github.com/toKrzysztof/kponos/internal/core/reference_analyzer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ServiceAccountHandler handles finding references to Secrets and ConfigMaps in ServiceAccount resources
type ServiceAccountHandler struct {
	client.Client
	referenceAnalyzer *core.ReferenceAnalyzer
	finders           map[string]ResourceReferenceFinder
}

// NewServiceAccountHandler creates a new ServiceAccountHandler
func NewServiceAccountHandler(c client.Client) *ServiceAccountHandler {
	analyzer := core.NewReferenceAnalyzer(c)
	h := &ServiceAccountHandler{
		Client:            c,
		referenceAnalyzer: analyzer,
	}

	h.finders = map[string]ResourceReferenceFinder{
		"Secret": h.findSecretReferences,
	}

	return h
}

// FindReferences finds all ServiceAccounts that reference the given resource
func (h *ServiceAccountHandler) FindReferences(ctx context.Context, c client.Client, resource client.Object, namespace string) ([]client.Object, error) {
	resourceKind := resource.GetObjectKind().GroupVersionKind().Kind
	resourceName := resource.GetName()

	finder, exists := h.finders[resourceKind]
	if !exists {
		return nil, fmt.Errorf("unsupported resource type: %s", resourceKind)
	}

	return finder(ctx, resourceName, namespace)
}

// GetResourceType returns the resource type this handler processes
func (h *ServiceAccountHandler) GetResourceType() string {
	return "ServiceAccount"
}

// findSecretReferences finds all ServiceAccounts that reference the given Secret
func (h *ServiceAccountHandler) findSecretReferences(ctx context.Context, resourceName, namespace string) ([]client.Object, error) {
	return h.referenceAnalyzer.FindReferencesForSecret(ctx, resourceName, namespace, "ServiceAccount")
}
