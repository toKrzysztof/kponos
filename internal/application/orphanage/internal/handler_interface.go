package handler

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ResourceHandler defines the interface for resource handlers
type ResourceHandler interface {
	// FindReferences finds all resources that reference a given Secret or ConfigMap
	FindReferences(ctx context.Context, client client.Client, secretName, configMapName string, namespace string) ([]client.Object, error)
	
	// GetResourceType returns the resource type this handler processes
	GetResourceType() string
}

