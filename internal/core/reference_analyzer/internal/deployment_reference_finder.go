package internal

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// TODO: Implement logic to find Deployments referencing the secret/configmap
// Check in spec.template.spec (Pod template):
// - spec.template.spec.volumes[].secret.secretName
// - spec.template.spec.volumes[].configMap.name
// - spec.template.spec.containers[].envFrom[].secretRef.name
// - spec.template.spec.containers[].envFrom[].configMapRef.name
// - spec.template.spec.containers[].env[].valueFrom.secretKeyRef.name
// - spec.template.spec.containers[].env[].valueFrom.configMapKeyRef.name
// - spec.template.spec.imagePullSecrets[].name

// DeploymentReferenceFinder finds references to Secrets and ConfigMaps in Deployment resources
type DeploymentReferenceFinder struct {
	client.Client
}

// NewDeploymentReferenceFinder creates a new DeploymentReferenceFinder
func NewDeploymentReferenceFinder(client client.Client) *DeploymentReferenceFinder {
	return &DeploymentReferenceFinder{
		Client: client,
	}
}

// FindSecretReferences finds all Deployments that reference the given Secret
func (f *DeploymentReferenceFinder) FindSecretReferences(ctx context.Context, client client.Client, secretName, namespace string) ([]client.Object, error) {
	// TODO: Implement logic to find Deployments referencing the secret
	return nil, nil
}

// FindConfigMapReferences finds all Deployments that reference the given ConfigMap
func (f *DeploymentReferenceFinder) FindConfigMapReferences(ctx context.Context, client client.Client, configMapName, namespace string) ([]client.Object, error) {
	// TODO: Implement logic to find Deployments referencing the configmap
	return nil, nil
}

// GetResourceType returns the Kubernetes resource type this strategy handles
func (f *DeploymentReferenceFinder) GetResourceType() string {
	return "Deployment"
}

