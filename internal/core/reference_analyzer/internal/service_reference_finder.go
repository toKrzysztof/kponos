package internal

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TODO: Implement logic to find Services referencing the secret/configmap
// Check:
// - spec.tls.secretName (for TLS secrets)

// ServiceReferenceFinder finds references to Secrets and ConfigMaps in Service resources
type ServiceReferenceFinder struct {
	client.Client
}

// NewServiceReferenceFinder creates a new ServiceReferenceFinder
func NewServiceReferenceFinder(client client.Client) *ServiceReferenceFinder {
	return &ServiceReferenceFinder{
		Client: client,
	}
}

// FindSecretReferences finds all Services that reference the given Secret
func (f *ServiceReferenceFinder) FindSecretReferences(ctx context.Context, client client.Client, secretName, namespace string) ([]client.Object, error) {
	// TODO: Implement logic to find Services referencing the secret
	return nil, nil
}

// Service does not reference ConfigMaps. This method is implemented to satisfy the ReferenceFinderStrategy interface.
func (f *ServiceReferenceFinder) FindConfigMapReferences(ctx context.Context, client client.Client, configMapName, namespace string) ([]client.Object, error) {
	return nil, nil
}

// GetResourceType returns the Kubernetes resource type this strategy handles
func (f *ServiceReferenceFinder) GetResourceType() string {
	return "Service"
}

