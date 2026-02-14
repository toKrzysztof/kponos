package resourceHandler

import (
	"context"

	core "github.com/toKrzysztof/kponos/internal/core/reference_analyzer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DaemonSetHandler handles finding references to Secrets and ConfigMaps in DaemonSet resources
type DaemonSetHandler struct {
	client.Client
	referenceAnalyzer *core.ReferenceAnalyzer
}

// NewDaemonSetHandler creates a new DaemonSetHandler
func NewDaemonSetHandler(client client.Client) *DaemonSetHandler {
	return &DaemonSetHandler{
		Client:            client,
		referenceAnalyzer: core.NewReferenceAnalyzer(client),
	}
}

// FindReferences finds all DaemonSets that reference the given Secret or ConfigMap
func (h *DaemonSetHandler) FindReferences(ctx context.Context, c client.Client, secretName, configMapName string, namespace string) ([]client.Object, error) {
	var daemonSets []client.Object

	// Find Secret references if secretName is provided
	if secretName != "" {
		secretRefs, err := h.referenceAnalyzer.FindReferencesForSecret(ctx, secretName, namespace, "DaemonSet")
		if err != nil {
			return nil, err
		}
		daemonSets = append(daemonSets, secretRefs...)
	}

	// Find ConfigMap references if configMapName is provided
	if configMapName != "" {
		configMapRefs, err := h.referenceAnalyzer.FindReferencesForConfigMap(ctx, configMapName, namespace, "DaemonSet")
		if err != nil {
			return nil, err
		}
		daemonSets = append(daemonSets, configMapRefs...)
	}

	return daemonSets, nil
}

// GetResourceType returns the resource type this handler processes
func (h *DaemonSetHandler) GetResourceType() string {
	return "DaemonSet"
}

