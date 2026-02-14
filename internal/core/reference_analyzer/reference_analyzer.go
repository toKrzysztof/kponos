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
	FindSecretReferences(ctx context.Context, client client.Client, secretName, namespace string) ([]client.Object, error)
	
	// FindConfigMapReferences finds all resources of this type that reference the given ConfigMap
	FindConfigMapReferences(ctx context.Context, client client.Client, configMapName, namespace string) ([]client.Object, error)
	
	// GetResourceType returns the Kubernetes resource type this strategy handles
	GetResourceType() string
}

// ReferenceAnalyzer finds resources that reference Secrets or ConfigMaps
type ReferenceAnalyzer struct {
	client.Client
	strategies map[string]ReferenceFinderStrategy
}

// NewReferenceAnalyzer creates a new ReferenceAnalyzer with all strategies initialized
func NewReferenceAnalyzer(client client.Client) *ReferenceAnalyzer {
	strategies := map[string]ReferenceFinderStrategy{
		"Pod":            internal.NewWorkloadReferenceFinder(client, internal.WorkloadResourceTypePod),
		"Deployment":    internal.NewWorkloadReferenceFinder(client, internal.WorkloadResourceTypeDeployment),
		"StatefulSet":   internal.NewWorkloadReferenceFinder(client, internal.WorkloadResourceTypeStatefulSet),
		"DaemonSet":     internal.NewWorkloadReferenceFinder(client, internal.WorkloadResourceTypeDaemonSet),
		"Ingress":       internal.NewIngressReferenceFinder(client),
		"ServiceAccount": internal.NewServiceAccountReferenceFinder(client),
	}
	
	return &ReferenceAnalyzer{
		Client:     client,
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

// IsOrphaned checks if a Secret or ConfigMap is orphaned (has no references)
func (s *ReferenceAnalyzer) IsOrphaned(ctx context.Context, secretName string, configMapName string, namespace string, resourceType string) (bool, error) {
	if secretName != "" {
		refs, err := s.FindReferencesForSecret(ctx, secretName, namespace, resourceType)
		if err != nil {
			return false, err
		}
		return len(refs) == 0, nil
	}
	
	if configMapName != "" {
		refs, err := s.FindReferencesForConfigMap(ctx, configMapName, namespace, resourceType)
		if err != nil {
			return false, err
		}
		return len(refs) == 0, nil
	}
	
	return false, nil
}