package internal

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ServiceAccountHandler handles finding references to Secrets and ConfigMaps in ServiceAccount resources
type ServiceAccountHandler struct {
	client.Client
}

// NewServiceAccountHandler creates a new ServiceAccountHandler
func NewServiceAccountHandler(client client.Client) *ServiceAccountHandler {
	return &ServiceAccountHandler{
		Client: client,
	}
}

// FindReferences finds all ServiceAccounts that reference the given Secret or ConfigMap
func (h *ServiceAccountHandler) FindReferences(ctx context.Context, c client.Client, secretName, configMapName string, namespace string) ([]client.Object, error) {
	var serviceAccounts []client.Object
	serviceAccountList := &corev1.ServiceAccountList{}
	if err := c.List(ctx, serviceAccountList, client.InNamespace(namespace)); err != nil {
		return serviceAccounts, err
	}
	
	// TODO: Filter serviceaccounts that reference the secret/configmap
	
	return serviceAccounts, nil
}

// GetResourceType returns the resource type this handler processes
func (h *ServiceAccountHandler) GetResourceType() string {
	return "ServiceAccount"
}

