package resourceHandler

import (
	"context"

	core "github.com/toKrzysztof/kponos/internal/core/reference_analyzer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ServiceAccountHandler handles finding references to Secrets and ConfigMaps in ServiceAccount resources
type ServiceAccountHandler struct {
	client.Client
	referenceAnalyzer *core.ReferenceAnalyzer
}

// NewServiceAccountHandler creates a new ServiceAccountHandler
func NewServiceAccountHandler(client client.Client) *ServiceAccountHandler {
	return &ServiceAccountHandler{
		Client:            client,
		referenceAnalyzer: core.NewReferenceAnalyzer(client),
	}
}

// TODO: refactor to not use configMapName
// FindReferences finds all ServiceAccounts that reference the given Secret
func (h *ServiceAccountHandler) FindReferences(ctx context.Context, c client.Client, secretName, configMapName string, namespace string) ([]client.Object, error) {
	var serviceAccounts []client.Object

	secretRefs, err := h.referenceAnalyzer.FindReferencesForSecret(ctx, secretName, namespace, "ServiceAccount")
	if err != nil {
		return nil, err
	}
	serviceAccounts = append(serviceAccounts, secretRefs...)

	return serviceAccounts, nil
}

// GetResourceType returns the resource type this handler processes
func (h *ServiceAccountHandler) GetResourceType() string {
	return "ServiceAccount"
}

