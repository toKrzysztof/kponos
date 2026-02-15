package internal

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ServiceAccountReferenceFinder finds references to Secrets and ConfigMaps in ServiceAccount resources
type ServiceAccountReferenceFinder struct {
	BaseReferenceFinder
}

// NewServiceAccountReferenceFinder creates a new ServiceAccountReferenceFinder
func NewServiceAccountReferenceFinder(c client.Client) *ServiceAccountReferenceFinder {
	return &ServiceAccountReferenceFinder{
		BaseReferenceFinder: BaseReferenceFinder{
			Client:       c,
			resourceType: "ServiceAccount",
		},
	}
}

// FindSecretReferences finds all ServiceAccounts that reference the given Secret
func (f *ServiceAccountReferenceFinder) FindSecretReferences(ctx context.Context, c client.Client, secretName, namespace string) ([]client.Object, error) {
	var results []client.Object

	serviceAccountList := &corev1.ServiceAccountList{}
	if err := c.List(ctx, serviceAccountList, client.InNamespace(namespace)); err != nil {
		return nil, err
	}

	for i := range serviceAccountList.Items {
		serviceAccount := &serviceAccountList.Items[i]
		if f.serviceAccountReferencesSecret(serviceAccount, secretName) {
			results = append(results, serviceAccount)
		}
	}

	return results, nil
}

// serviceAccountReferencesSecret checks if a ServiceAccount references the given secret
func (f *ServiceAccountReferenceFinder) serviceAccountReferencesSecret(serviceAccount *corev1.ServiceAccount, secretName string) bool {
	// Check secrets[].name (for image pull secrets and mounted secrets)
	for _, secret := range serviceAccount.Secrets {
		if secret.Name == secretName {
			return true
		}
	}

	// Check imagePullSecrets[].name (for image pull secrets)
	for _, imagePullSecret := range serviceAccount.ImagePullSecrets {
		if imagePullSecret.Name == secretName {
			return true
		}
	}

	return false
}
