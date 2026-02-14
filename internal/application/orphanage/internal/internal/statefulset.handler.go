package resourceHandler

import (
	"context"

	core "github.com/toKrzysztof/kponos/internal/core/reference_analyzer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StatefulSetHandler handles finding references to Secrets and ConfigMaps in StatefulSet resources
type StatefulSetHandler struct {
	client.Client
	referenceAnalyzer *core.ReferenceAnalyzer
}

// NewStatefulSetHandler creates a new StatefulSetHandler
func NewStatefulSetHandler(client client.Client) *StatefulSetHandler {
	return &StatefulSetHandler{
		Client:            client,
		referenceAnalyzer: core.NewReferenceAnalyzer(client),
	}
}

// FindReferences finds all StatefulSets that reference the given Secret or ConfigMap
func (h *StatefulSetHandler) FindReferences(ctx context.Context, c client.Client, secretName, configMapName string, namespace string) ([]client.Object, error) {
	var statefulSets []client.Object

	// Find Secret references if secretName is provided
	if secretName != "" {
		secretRefs, err := h.referenceAnalyzer.FindReferencesForSecret(ctx, secretName, namespace, "StatefulSet")
		if err != nil {
			return nil, err
		}
		statefulSets = append(statefulSets, secretRefs...)
	}

	// Find ConfigMap references if configMapName is provided
	if configMapName != "" {
		configMapRefs, err := h.referenceAnalyzer.FindReferencesForConfigMap(ctx, configMapName, namespace, "StatefulSet")
		if err != nil {
			return nil, err
		}
		statefulSets = append(statefulSets, configMapRefs...)
	}

	return statefulSets, nil
}

// GetResourceType returns the resource type this handler processes
func (h *StatefulSetHandler) GetResourceType() string {
	return "StatefulSet"
}

