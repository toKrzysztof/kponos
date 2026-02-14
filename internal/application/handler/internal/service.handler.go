package internal

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ServiceHandler handles finding references to Secrets and ConfigMaps in Service resources
type ServiceHandler struct {
	client.Client
}

// NewServiceHandler creates a new ServiceHandler
func NewServiceHandler(client client.Client) *ServiceHandler {
	return &ServiceHandler{
		Client: client,
	}
}

// FindReferences finds all Services that reference the given Secret or ConfigMap
func (h *ServiceHandler) FindReferences(ctx context.Context, c client.Client, secretName, configMapName string, namespace string) ([]client.Object, error) {
	var services []client.Object
	serviceList := &corev1.ServiceList{}
	if err := c.List(ctx, serviceList, client.InNamespace(namespace)); err != nil {
		return services, err
	}
	
	// TODO: Filter services that reference the secret/configmap
	
	return services, nil
}

// GetResourceType returns the resource type this handler processes
func (h *ServiceHandler) GetResourceType() string {
	return "Service"
}

