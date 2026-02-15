package core

import (
	"context"
	"fmt"

	"github.com/toKrzysztof/kponos/internal/core/reference_analyzer/internal"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ReferenceFinderStrategy defines the interface for finding references in a specific resource type
type ReferenceFinderStrategy interface {
	// FindSecretReferences finds all resources of this type that reference the given Secret
	FindSecretReferences(ctx context.Context, c client.Client, secretName, namespace string) ([]client.Object, error)

	// FindConfigMapReferences finds all resources of this type that reference the given ConfigMap
	FindConfigMapReferences(ctx context.Context, c client.Client, configMapName, namespace string) ([]client.Object, error)

	// FindServiceReferences finds all resources of this type that reference the given Service
	FindServiceReferences(ctx context.Context, c client.Client, serviceName, namespace string) ([]client.Object, error)

	// FindPodReferences finds all resources of this type that reference the given Pod
	FindPodReferences(ctx context.Context, c client.Client, podName, namespace string) ([]client.Object, error)

	// GetResourceType returns the Kubernetes resource type this strategy handles
	GetResourceType() string
}

// ReferenceAnalyzer finds resources that reference Secrets or ConfigMaps
type ReferenceAnalyzer struct {
	client.Client
	strategies map[string]ReferenceFinderStrategy
}

// NewReferenceAnalyzer creates a new ReferenceAnalyzer with all strategies initialized
func NewReferenceAnalyzer(c client.Client) *ReferenceAnalyzer {
	strategies := map[string]ReferenceFinderStrategy{
		"Pod":                            internal.NewWorkloadReferenceFinder(c, internal.WorkloadResourceTypePod),
		"Deployment":                     internal.NewWorkloadReferenceFinder(c, internal.WorkloadResourceTypeDeployment),
		"StatefulSet":                    internal.NewWorkloadReferenceFinder(c, internal.WorkloadResourceTypeStatefulSet),
		"DaemonSet":                      internal.NewWorkloadReferenceFinder(c, internal.WorkloadResourceTypeDaemonSet),
		"Ingress":                        internal.NewIngressReferenceFinder(c),
		"ServiceAccount":                 internal.NewServiceAccountReferenceFinder(c),
		"ValidatingWebhookConfiguration": internal.NewValidatingWebhookConfigurationReferenceFinder(c),
		"MutatingWebhookConfiguration":   internal.NewMutatingWebhookConfigurationReferenceFinder(c),
		"APIService":                     internal.NewAPIServiceReferenceFinder(c),
		"CustomResourceDefinition":       internal.NewCustomResourceDefinitionReferenceFinder(c),
		"Service":                        internal.NewServiceReferenceFinder(c),
	}

	return &ReferenceAnalyzer{
		Client:     c,
		strategies: strategies,
	}
}

// FindReferencesForSecret finds all resources that reference the given Secret
// If resourceType is provided, only searches that resource type; otherwise searches all types
func (s *ReferenceAnalyzer) FindReferencesForSecret(ctx context.Context, secretName, namespace string, resourceType string) ([]client.Object, error) {
	strategy := s.strategies[resourceType]
	if strategy == nil {
		return nil, fmt.Errorf("unknown resource type: %s", resourceType)
	}

	return strategy.FindSecretReferences(ctx, s.Client, secretName, namespace)
}

// FindReferencesForConfigMap finds all resources that reference the given ConfigMap
// If resourceType is provided, only searches that resource type; otherwise searches all types
func (s *ReferenceAnalyzer) FindReferencesForConfigMap(ctx context.Context, configMapName, namespace string, resourceType string) ([]client.Object, error) {
	strategy := s.strategies[resourceType]
	if strategy == nil {
		return nil, fmt.Errorf("unknown resource type: %s", resourceType)
	}

	return strategy.FindConfigMapReferences(ctx, s.Client, configMapName, namespace)
}

func (s *ReferenceAnalyzer) FindReferencesForService(ctx context.Context, serviceName, namespace string, resourceType string) ([]client.Object, error) {
	strategy := s.strategies[resourceType]
	if strategy == nil {
		return nil, fmt.Errorf("unknown resource type: %s", resourceType)
	}

	return strategy.FindServiceReferences(ctx, s.Client, serviceName, namespace)
}
