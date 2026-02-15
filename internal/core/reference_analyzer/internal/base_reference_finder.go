package internal

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// BaseReferenceFinder provides default (no-op) implementations for all ReferenceFinderStrategy methods.
// Concrete finders can embed this struct and only override the methods they need.
type BaseReferenceFinder struct {
	client.Client
	resourceType string
}

// FindSecretReferences provides a default no-op implementation.
// Override this method in concrete finders that need to find Secret references.
func (f *BaseReferenceFinder) FindSecretReferences(ctx context.Context, c client.Client, secretName, namespace string) ([]client.Object, error) {
	return nil, nil
}

// FindConfigMapReferences provides a default no-op implementation.
// Override this method in concrete finders that need to find ConfigMap references.
func (f *BaseReferenceFinder) FindConfigMapReferences(ctx context.Context, c client.Client, configMapName, namespace string) ([]client.Object, error) {
	return nil, nil
}

// FindServiceReferences provides a default no-op implementation.
// Override this method in concrete finders that need to find Service references.
func (f *BaseReferenceFinder) FindServiceReferences(ctx context.Context, c client.Client, serviceName, namespace string) ([]client.Object, error) {
	return nil, nil
}

// FindPodReferences provides a default no-op implementation.
// Override this method in concrete finders that need to find Pod references.
func (f *BaseReferenceFinder) FindPodReferences(ctx context.Context, c client.Client, podName, namespace string) ([]client.Object, error) {
	return nil, nil
}

// GetResourceType returns the Kubernetes resource type this strategy handles.
// Override this method in concrete finders to return the specific resource type.
func (f *BaseReferenceFinder) GetResourceType() string {
	return f.resourceType
}
