package resourceHandler

import (
	"context"

	core "github.com/toKrzysztof/kponos/internal/core/reference_analyzer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PodHandler handles finding references to Secrets and ConfigMaps in Pod resources
type PodHandler struct {
	client.Client
	referenceAnalyzer *core.ReferenceAnalyzer
}

// NewPodHandler creates a new PodHandler
func NewPodHandler(client client.Client) *PodHandler {
	return &PodHandler{
		Client:            client,
		referenceAnalyzer: core.NewReferenceAnalyzer(client),
	}
}

// FindReferences finds all Pods that reference the given Secret or ConfigMap
func (h *PodHandler) FindReferences(ctx context.Context, c client.Client, secretName, configMapName string, namespace string) ([]client.Object, error) {
	var pods []client.Object

	// Find Secret references if secretName is provided
	if secretName != "" {
		secretRefs, err := h.referenceAnalyzer.FindReferencesForSecret(ctx, secretName, namespace, "Pod")
		if err != nil {
			return nil, err
		}
		pods = append(pods, secretRefs...)
	}

	// Find ConfigMap references if configMapName is provided
	if configMapName != "" {
		configMapRefs, err := h.referenceAnalyzer.FindReferencesForConfigMap(ctx, configMapName, namespace, "Pod")
		if err != nil {
			return nil, err
		}
		pods = append(pods, configMapRefs...)
	}

	return pods, nil
}

// GetResourceType returns the resource type this handler processes
func (h *PodHandler) GetResourceType() string {
	return "Pod"
}

