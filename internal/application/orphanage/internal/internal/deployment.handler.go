package resourceHandler

import (
	"context"

	core "github.com/toKrzysztof/kponos/internal/core/reference_analyzer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DeploymentHandler handles finding references to Secrets and ConfigMaps in Deployment resources
type DeploymentHandler struct {
	client.Client
	referenceAnalyzer *core.ReferenceAnalyzer
}

// NewDeploymentHandler creates a new DeploymentHandler
func NewDeploymentHandler(client client.Client) *DeploymentHandler {
	return &DeploymentHandler{
		Client:            client,
		referenceAnalyzer: core.NewReferenceAnalyzer(client),
	}
}

// FindReferences finds all Deployments that reference the given Secret or ConfigMap
func (h *DeploymentHandler) FindReferences(ctx context.Context, c client.Client, secretName, configMapName string, namespace string) ([]client.Object, error) {
	var deployments []client.Object

	// Find Secret references if secretName is provided
	if secretName != "" {
		secretRefs, err := h.referenceAnalyzer.FindReferencesForSecret(ctx, secretName, namespace, "Deployment")
		if err != nil {
			return nil, err
		}
		deployments = append(deployments, secretRefs...)
	}

	// Find ConfigMap references if configMapName is provided
	if configMapName != "" {
		configMapRefs, err := h.referenceAnalyzer.FindReferencesForConfigMap(ctx, configMapName, namespace, "Deployment")
		if err != nil {
			return nil, err
		}
		deployments = append(deployments, configMapRefs...)
	}

	return deployments, nil
}

// GetResourceType returns the resource type this handler processes
func (h *DeploymentHandler) GetResourceType() string {
	return "Deployment"
}

