package internal

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// APIServiceReferenceFinder finds references to Services in APIService resources
type APIServiceReferenceFinder struct {
	BaseReferenceFinder
}

// NewAPIServiceReferenceFinder creates a new APIServiceReferenceFinder
func NewAPIServiceReferenceFinder(c client.Client) *APIServiceReferenceFinder {
	return &APIServiceReferenceFinder{
		BaseReferenceFinder: BaseReferenceFinder{
			Client:       c,
			resourceType: "APIService",
		},
	}
}

// FindServiceReferences finds all APIServices that reference the given Service
func (f *APIServiceReferenceFinder) FindServiceReferences(ctx context.Context, c client.Client, serviceName, namespace string) ([]client.Object, error) {
	var results []client.Object

	apiServiceList := &unstructured.UnstructuredList{}
	apiServiceList.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "apiregistration.k8s.io",
		Version: "v1",
		Kind:    "APIServiceList",
	})

	// APIService is cluster-scoped, so we don't filter by namespace
	if err := c.List(ctx, apiServiceList, &client.ListOptions{}); err != nil {
		return nil, err
	}

	for i := range apiServiceList.Items {
		apiService := &apiServiceList.Items[i]
		if f.apiServiceReferencesService(apiService, serviceName, namespace) {
			results = append(results, apiService)
		}
	}

	return results, nil
}

// apiServiceReferencesService checks if an APIService references the given service
func (f *APIServiceReferenceFinder) apiServiceReferencesService(apiService *unstructured.Unstructured, serviceName, namespace string) bool {
	// Check spec.service.name and spec.service.namespace
	spec, found, err := unstructured.NestedMap(apiService.Object, "spec")
	if !found || err != nil {
		return false
	}

	service, found, err := unstructured.NestedMap(spec, "service")
	if !found || err != nil {
		return false
	}

	svcName, found, _ := unstructured.NestedString(service, "name")
	svcNamespace, foundNS, _ := unstructured.NestedString(service, "namespace")

	if found && svcName == serviceName {
		// If namespace is specified, it must match; if not specified, it defaults to the target namespace
		if !foundNS || svcNamespace == namespace {
			return true
		}
	}

	return false
}
