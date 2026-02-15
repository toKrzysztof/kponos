package internal

import (
	"context"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CustomResourceDefinitionReferenceFinder finds references to Services in CustomResourceDefinition resources
type CustomResourceDefinitionReferenceFinder struct {
	client.Client
}

// NewCustomResourceDefinitionReferenceFinder creates a new CustomResourceDefinitionReferenceFinder
func NewCustomResourceDefinitionReferenceFinder(c client.Client) *CustomResourceDefinitionReferenceFinder {
	return &CustomResourceDefinitionReferenceFinder{
		Client: c,
	}
}

// CustomResourceDefinition does not reference Secrets. This method is implemented to satisfy the ReferenceFinderStrategy interface.
func (f *CustomResourceDefinitionReferenceFinder) FindSecretReferences(ctx context.Context, c client.Client, secretName, namespace string) ([]client.Object, error) {
	return nil, nil
}

// CustomResourceDefinition does not reference ConfigMaps. This method is implemented to satisfy the ReferenceFinderStrategy interface.
func (f *CustomResourceDefinitionReferenceFinder) FindConfigMapReferences(ctx context.Context, c client.Client, configMapName, namespace string) ([]client.Object, error) {
	return nil, nil
}

// FindServiceReferences finds all CustomResourceDefinitions that reference the given Service
func (f *CustomResourceDefinitionReferenceFinder) FindServiceReferences(ctx context.Context, c client.Client, serviceName, namespace string) ([]client.Object, error) {
	var results []client.Object

	crdList := &apiextensionsv1.CustomResourceDefinitionList{}
	// CustomResourceDefinition is cluster-scoped, so we don't filter by namespace
	if err := c.List(ctx, crdList); err != nil {
		return nil, err
	}

	for i := range crdList.Items {
		crd := &crdList.Items[i]
		if f.customResourceDefinitionReferencesService(crd, serviceName, namespace) {
			results = append(results, crd)
		}
	}

	return results, nil
}

// customResourceDefinitionReferencesService checks if a CustomResourceDefinition references the given service
func (f *CustomResourceDefinitionReferenceFinder) customResourceDefinitionReferencesService(crd *apiextensionsv1.CustomResourceDefinition, serviceName, namespace string) bool {
	// Check spec.conversion.webhook.clientConfig.service.name and spec.conversion.webhook.clientConfig.service.namespace
	if crd.Spec.Conversion != nil && crd.Spec.Conversion.Webhook != nil && crd.Spec.Conversion.Webhook.ClientConfig.Service != nil {
		if crd.Spec.Conversion.Webhook.ClientConfig.Service.Name == serviceName && crd.Spec.Conversion.Webhook.ClientConfig.Service.Namespace == namespace {
			return true
		}
	}

	return false
}

// GetResourceType returns the Kubernetes resource type this strategy handles
func (f *CustomResourceDefinitionReferenceFinder) GetResourceType() string {
	return "CustomResourceDefinition"
}
