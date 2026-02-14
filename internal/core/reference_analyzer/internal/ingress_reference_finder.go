package internal

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TODO: Implement logic to find Ingresses referencing the secret/configmap
// Check:
// - spec.tls[].secretName (for TLS secrets)

// IngressReferenceFinder finds references to Secrets and ConfigMaps in Ingress resources
type IngressReferenceFinder struct {
	client.Client
}

// NewIngressReferenceFinder creates a new IngressReferenceFinder
func NewIngressReferenceFinder(client client.Client) *IngressReferenceFinder {
	return &IngressReferenceFinder{
		Client: client,
	}
}

// FindSecretReferences finds all Ingresses that reference the given Secret
func (f *IngressReferenceFinder) FindSecretReferences(ctx context.Context, client client.Client, secretName, namespace string) ([]client.Object, error) {
	// TODO: Implement logic to find Ingresses referencing the secret
	return nil, nil
}

// Ingress does not reference ConfigMaps. This method is implemented to satisfy the ReferenceFinderStrategy interface.
func (f *IngressReferenceFinder) FindConfigMapReferences(ctx context.Context, client client.Client, configMapName, namespace string) ([]client.Object, error) {
	return nil, nil
}

// GetResourceType returns the Kubernetes resource type this strategy handles
func (f *IngressReferenceFinder) GetResourceType() string {
	return "Ingress"
}

