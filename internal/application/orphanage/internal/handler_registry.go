package handlerRegistry

import (
	"context"

	resourceHandler "github.com/toKrzysztof/kponos/internal/application/orphanage/internal/internal"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ResourceHandler interface {
	// FindReferences finds all resources that reference a given Secret or ConfigMap
	FindReferences(ctx context.Context, c client.Client, resource client.Object, namespace string) ([]client.Object, error)

	// GetResourceType returns the resource type this handler processes
	GetResourceType() string
}

// Registry holds all resource handlers
type HandlerRegistry struct {
	handlers map[string]ResourceHandler
}

// NewHandlerRegistry creates a new handler registry with all handlers initialized
func NewHandlerRegistry(c client.Client) *HandlerRegistry {
	return &HandlerRegistry{
		// TODO: replace strings with strictly typed enums
		handlers: map[string]ResourceHandler{
			"Pod":            resourceHandler.NewPodHandler(c),
			"Deployment":     resourceHandler.NewDeploymentHandler(c),
			"StatefulSet":    resourceHandler.NewStatefulSetHandler(c),
			"DaemonSet":      resourceHandler.NewDaemonSetHandler(c),
			"Ingress":        resourceHandler.NewIngressHandler(c),
			"ServiceAccount": resourceHandler.NewServiceAccountHandler(c),
		},
	}
}

// GetHandler returns a handler for the given resource type
func (r *HandlerRegistry) GetHandler(resourceType string) ResourceHandler {
	return r.handlers[resourceType]
}
