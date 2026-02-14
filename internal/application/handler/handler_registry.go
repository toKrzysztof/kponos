package application

import (
	"github.com/toKrzysztof/kponos/internal/application/handler/internal"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Registry holds all resource handlers
type HandlerRegistry struct {
	handlers map[string]ResourceHandler
}

// NewHandlerRegistry creates a new handler registry with all handlers initialized
func NewHandlerRegistry(client client.Client) *HandlerRegistry {
	return &HandlerRegistry{
		handlers: map[string]ResourceHandler{
			"Pod":            internal.NewPodHandler(client),
			"Deployment":    internal.NewDeploymentHandler(client),
			"StatefulSet":   internal.NewStatefulSetHandler(client),
			"DaemonSet":     internal.NewDaemonSetHandler(client),
			"Service":       internal.NewServiceHandler(client),
			"Ingress":       internal.NewIngressHandler(client),
			"ServiceAccount": internal.NewServiceAccountHandler(client),
		},
	}
}

// GetHandler returns a handler for the given resource type
func (r *HandlerRegistry) GetHandler(resourceType string) ResourceHandler {
	return r.handlers[resourceType]
}
