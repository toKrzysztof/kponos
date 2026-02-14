package internal

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// WorkloadResourceType represents a valid workload resource type
type WorkloadResourceType string

const (
	WorkloadResourceTypePod        WorkloadResourceType = "Pod"
	WorkloadResourceTypeDeployment WorkloadResourceType = "Deployment"
	WorkloadResourceTypeStatefulSet WorkloadResourceType = "StatefulSet"
	WorkloadResourceTypeDaemonSet   WorkloadResourceType = "DaemonSet"
)

// WorkloadReferenceFinder finds references to Secrets and ConfigMaps in workload resources
// that are Pods or create Pods (Deployment, StatefulSet, DaemonSet)
type WorkloadReferenceFinder struct {
	client.Client
	resourceType WorkloadResourceType
}

// NewWorkloadReferenceFinder creates a new WorkloadReferenceFinder for the given resource type
func NewWorkloadReferenceFinder(client client.Client, resourceType WorkloadResourceType) *WorkloadReferenceFinder {
	return &WorkloadReferenceFinder{
		Client:       client,
		resourceType: resourceType,
	}
}

// getPathPrefix returns the path prefix to use when searching for references
// Pods check spec directly, while workload resources check spec.template.spec
func (f *WorkloadReferenceFinder) getPathPrefix() string {
	if f.resourceType == WorkloadResourceTypePod {
		return ""
	}
	return "spec.template.spec"
}

// FindSecretReferences finds all resources that reference the given Secret
func (f *WorkloadReferenceFinder) FindSecretReferences(ctx context.Context, client client.Client, secretName, namespace string) ([]client.Object, error) {
	// TODO: Implement logic to find resources referencing the secret
	// Use f.getPathPrefix() to determine where to search:
	// - {prefix}.volumes[].secret.secretName
	// - {prefix}.containers[].envFrom[].secretRef.name
	// - {prefix}.containers[].env[].valueFrom.secretKeyRef.name
	// - {prefix}.imagePullSecrets[].name
	return nil, nil
}

// FindConfigMapReferences finds all resources that reference the given ConfigMap
func (f *WorkloadReferenceFinder) FindConfigMapReferences(ctx context.Context, client client.Client, configMapName, namespace string) ([]client.Object, error) {
	// TODO: Implement logic to find resources referencing the configmap
	// Use f.getPathPrefix() to determine where to search:
	// - {prefix}.volumes[].configMap.name
	// - {prefix}.containers[].envFrom[].configMapRef.name
	// - {prefix}.containers[].env[].valueFrom.configMapKeyRef.name
	return nil, nil
}

// GetResourceType returns the Kubernetes resource type this strategy handles
func (f *WorkloadReferenceFinder) GetResourceType() string {
	return string(f.resourceType)
}