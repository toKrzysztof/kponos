package resourceHandler

import (
	"context"

	core "github.com/toKrzysztof/kponos/internal/core/reference_analyzer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IngressHandler handles finding references to Secrets and ConfigMaps in Ingress resources
type IngressHandler struct {
	client.Client
	referenceAnalyzer *core.ReferenceAnalyzer
}

// NewIngressHandler creates a new IngressHandler
func NewIngressHandler(client client.Client) *IngressHandler {
	return &IngressHandler{
		Client:            client,
		referenceAnalyzer: core.NewReferenceAnalyzer(client),
	}
}

// FindReferences finds all Ingresses that reference the given Secret
func (h *IngressHandler) FindReferences(ctx context.Context, c client.Client, secretName string, configMapName string, namespace string) ([]client.Object, error) {
	var ingresses []client.Object

	// Find Secret references if secretName is provided
	secretRefs, err := h.referenceAnalyzer.FindReferencesForSecret(ctx, secretName, namespace, "Ingress")
	if err != nil {
		return nil, err
	}
	ingresses = append(ingresses, secretRefs...)


	return ingresses, nil
}

// GetResourceType returns the resource type this handler processes
func (h *IngressHandler) GetResourceType() string {
	return "Ingress"
}

