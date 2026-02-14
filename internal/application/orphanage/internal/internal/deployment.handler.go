package resourceHandler

import (
	"context"
	"fmt"

	core "github.com/toKrzysztof/kponos/internal/core/reference_analyzer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DeploymentHandler handles finding references to Secrets and ConfigMaps in Deployment resources
type DeploymentHandler struct {
	client.Client
	referenceAnalyzer *core.ReferenceAnalyzer
	finders           map[string]ResourceReferenceFinder
}

// NewDeploymentHandler creates a new DeploymentHandler
func NewDeploymentHandler(client client.Client) *DeploymentHandler {
	analyzer := core.NewReferenceAnalyzer(client)
	h := &DeploymentHandler{
		Client:            client,
		referenceAnalyzer: analyzer,
	}
	
	h.finders = map[string]ResourceReferenceFinder{
		"Secret":    h.findSecretReferences,
		"ConfigMap": h.findConfigMapReferences,
	}
	
	return h
}

// FindReferences finds all Deployments that reference the given resource
func (h *DeploymentHandler) FindReferences(ctx context.Context, c client.Client, resource client.Object, namespace string) ([]client.Object, error) {
	resourceKind := resource.GetObjectKind().GroupVersionKind().Kind
	resourceName := resource.GetName()
	
	finder, exists := h.finders[resourceKind]
	if !exists {
		return nil, fmt.Errorf("unsupported resource type: %s", resourceKind)
	}
	
	return finder(ctx, resourceName, namespace)
}

// GetResourceType returns the resource type this handler processes
func (h *DeploymentHandler) GetResourceType() string {
	return "Deployment"
}

// findSecretReferences finds all Deployments that reference the given Secret
func (h *DeploymentHandler) findSecretReferences(ctx context.Context, resourceName, namespace string) ([]client.Object, error) {
	return h.referenceAnalyzer.FindReferencesForSecret(ctx, resourceName, namespace, "Deployment")
}

// findConfigMapReferences finds all Deployments that reference the given ConfigMap
func (h *DeploymentHandler) findConfigMapReferences(ctx context.Context, resourceName, namespace string) ([]client.Object, error) {
	return h.referenceAnalyzer.FindReferencesForConfigMap(ctx, resourceName, namespace, "Deployment")
}
