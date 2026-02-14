package application

import (
	"context"
	"fmt"

	handler "github.com/toKrzysztof/kponos/internal/application/orphanage/internal"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Orphanage handles finding orphaned resources in a namespace
type Orphanage struct {
	client          client.Client
	handlerRegistry *handler.HandlerRegistry
}

// NewOrphanage creates a new Orphanage instance
func NewOrphanage(client client.Client) *Orphanage {
	return &Orphanage{
		client:          client,
		handlerRegistry: handler.NewHandlerRegistry(client),
	}
}

// FindOrphans lists all daemonsets, deployments, ingresses, pods, services, serviceaccounts and statefulsets
// in a namespace and for each invokes an appropriate handler.FindReferences
func (o *Orphanage) FindOrphans(ctx context.Context, namespace, secretName, configMapName string) (map[string][]client.Object, error) {
	results := make(map[string][]client.Object)

	// Resource types to check
	resourceTypes := []string{
		"DaemonSet",
		"Deployment",
		"Ingress",
		"Pod",
		"Service",
		"ServiceAccount",
		"StatefulSet",
	}

	// For each resource type, get the handler and call FindReferences
	for _, resourceType := range resourceTypes {
		handler := o.handlerRegistry.GetHandler(resourceType)
		if handler == nil {
			return nil, fmt.Errorf("no handler found for resource type: %s", resourceType)
		}

		references, err := handler.FindReferences(ctx, o.client, secretName, configMapName, namespace)
		if err != nil {
			return nil, fmt.Errorf("error finding references for %s: %w", resourceType, err)
		}

		results[resourceType] = references
	}

	return results, nil
}
