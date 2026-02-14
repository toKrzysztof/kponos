package application

import (
	"context"
	"fmt"

	handlerRegistry "github.com/toKrzysztof/kponos/internal/application/orphanage/internal"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// OrphanFinder is a function that finds orphaned resources of a specific type
type OrphanFinder func(context.Context, string) ([]client.Object, error)

// Orphanage handles finding orphaned resources in a namespace
type Orphanage struct {
	client          client.Client
	handlerRegistry *handlerRegistry.HandlerRegistry
	finders         map[string]OrphanFinder
}

// NewOrphanage creates a new Orphanage instance
func NewOrphanage(c client.Client) *Orphanage {
	o := &Orphanage{
		client:          c,
		handlerRegistry: handlerRegistry.NewHandlerRegistry(c),
	}

	o.finders = map[string]OrphanFinder{
		"Secret":    o.findOrphanedSecrets,
		"ConfigMap": o.findOrphanedConfigMaps,
	}

	return o
}

// FindOrphans finds all orphaned Secrets and ConfigMaps in a namespace.
// An orphan is a Secret or ConfigMap that is not referenced by any other resources.
func (o *Orphanage) FindOrphans(ctx context.Context, resourceType string, namespace string) ([]client.Object, error) {
	finder, exists := o.finders[resourceType]
	if !exists {
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}

	return finder(ctx, namespace)
}

// findOrphanedSecrets finds all orphaned Secrets in the given namespace
func (o *Orphanage) findOrphanedSecrets(ctx context.Context, namespace string) ([]client.Object, error) {
	var orphanedSecrets []client.Object

	secretList := &corev1.SecretList{}
	if err := o.client.List(ctx, secretList, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("unable to list Secrets: %w", err)
	}

	for i := range secretList.Items {
		secret := &secretList.Items[i]
		if isOrphaned, err := o.isOrphaned(ctx, secret, namespace); err != nil {
			return nil, fmt.Errorf("error checking if Secret %s is orphaned: %w", secret.Name, err)
		} else if isOrphaned {
			orphanedSecrets = append(orphanedSecrets, secret)
		}
	}

	return orphanedSecrets, nil
}

// findOrphanedConfigMaps finds all orphaned ConfigMaps in the given namespace
func (o *Orphanage) findOrphanedConfigMaps(ctx context.Context, namespace string) ([]client.Object, error) {
	var orphanedConfigMaps []client.Object

	configMapList := &corev1.ConfigMapList{}
	if err := o.client.List(ctx, configMapList, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("unable to list ConfigMaps: %w", err)
	}

	for i := range configMapList.Items {
		configMap := &configMapList.Items[i]
		if isOrphaned, err := o.isOrphaned(ctx, configMap, namespace); err != nil {
			return nil, fmt.Errorf("error checking if ConfigMap %s is orphaned: %w", configMap.Name, err)
		} else if isOrphaned {
			orphanedConfigMaps = append(orphanedConfigMaps, configMap)
		}
	}

	return orphanedConfigMaps, nil
}

// isOrphaned checks if a Secret or ConfigMap is orphaned (not referenced by any resources).
func (o *Orphanage) isOrphaned(ctx context.Context, resource client.Object, namespace string) (bool, error) {
	// Check if any of these resources reference the secret/configmap
	resourceTypes := []string{
		"DaemonSet",
		"Deployment",
		"Ingress",
		"Pod",
		"Service",
		"ServiceAccount",
		"StatefulSet",
	}

	for _, resourceType := range resourceTypes {
		handler := o.handlerRegistry.GetHandler(resourceType)
		if handler == nil {
			return false, fmt.Errorf("no handler found for resource type: %s", resourceType)
		}

		references, err := handler.FindReferences(ctx, o.client, resource, namespace)
		if err != nil {
			return false, fmt.Errorf("error finding references for %s: %w", resourceType, err)
		}

		if len(references) > 0 {
			return false, nil
		}
	}

	return true, nil
}
