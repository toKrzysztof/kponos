package application

import (
	"context"
	"fmt"

	handler "github.com/toKrzysztof/kponos/internal/application/orphanage/internal"
	corev1 "k8s.io/api/core/v1"
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

// FindOrphans finds all orphaned Secrets and ConfigMaps in a namespace.
// An orphan is a Secret or ConfigMap that is not referenced by any other resources.
func (o *Orphanage) FindOrphans(ctx context.Context, namespace string) (map[string][]client.Object, error) {
	results := make(map[string][]client.Object)
	var orphanedSecrets []client.Object
	var orphanedConfigMaps []client.Object

	secretList := &corev1.SecretList{}
	if err := o.client.List(ctx, secretList, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("unable to list Secrets: %w", err)
	}

	for i := range secretList.Items {
		secret := &secretList.Items[i]
		if isOrphaned, err := o.isOrphaned(ctx, secret.Name, "", namespace); err != nil {
			return nil, fmt.Errorf("error checking if Secret %s is orphaned: %w", secret.Name, err)
		} else if isOrphaned {
			orphanedSecrets = append(orphanedSecrets, secret)
		}
	}

	configMapList := &corev1.ConfigMapList{}
	if err := o.client.List(ctx, configMapList, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("unable to list ConfigMaps: %w", err)
	}

	for i := range configMapList.Items {
		configMap := &configMapList.Items[i]
		if isOrphaned, err := o.isOrphaned(ctx, "", configMap.Name, namespace); err != nil {
			return nil, fmt.Errorf("error checking if ConfigMap %s is orphaned: %w", configMap.Name, err)
		} else if isOrphaned {
			orphanedConfigMaps = append(orphanedConfigMaps, configMap)
		}
	}

	results["Secret"] = orphanedSecrets
	results["ConfigMap"] = orphanedConfigMaps

	return results, nil
}

// isOrphaned checks if a Secret or ConfigMap is orphaned (not referenced by any resources).
func (o *Orphanage) isOrphaned(ctx context.Context, secretName, configMapName, namespace string) (bool, error) {
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

		references, err := handler.FindReferences(ctx, o.client, secretName, configMapName, namespace)
		if err != nil {
			return false, fmt.Errorf("error finding references for %s: %w", resourceType, err)
		}

		if len(references) > 0 {
			return false, nil
		}
	}

	return true, nil
}