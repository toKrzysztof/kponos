package internal

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TODO: Implement logic to find ServiceAccounts referencing the secret/configmap
// Check:
// - secrets[].name (for image pull secrets and mounted secrets)
// - imagePullSecrets[].name (for image pull secrets)

// ServiceAccountReferenceFinder finds references to Secrets and ConfigMaps in ServiceAccount resources
type ServiceAccountReferenceFinder struct {
	client.Client
}

// NewServiceAccountReferenceFinder creates a new ServiceAccountReferenceFinder
func NewServiceAccountReferenceFinder(client client.Client) *ServiceAccountReferenceFinder {
	return &ServiceAccountReferenceFinder{
		Client: client,
	}
}

// FindSecretReferences finds all ServiceAccounts that reference the given Secret
func (f *ServiceAccountReferenceFinder) FindSecretReferences(ctx context.Context, client client.Client, secretName, namespace string) ([]client.Object, error) {
	// TODO: Implement logic to find ServiceAccounts referencing the secret
	return nil, nil
}

// ServiceAccount does not reference ConfigMaps. This method is implemented to satisfy the ReferenceFinderStrategy interface.
func (f *ServiceAccountReferenceFinder) FindConfigMapReferences(ctx context.Context, client client.Client, configMapName, namespace string) ([]client.Object, error) {
	return nil, nil
}

// GetResourceType returns the Kubernetes resource type this strategy handles
func (f *ServiceAccountReferenceFinder) GetResourceType() string {
	return "ServiceAccount"
}

