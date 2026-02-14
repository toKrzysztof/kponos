package internal

import (
	"context"

	networkingv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IngressReferenceFinder finds references to Secrets and ConfigMaps in Ingress resources
type IngressReferenceFinder struct {
	client.Client
}

// NewIngressReferenceFinder creates a new IngressReferenceFinder
func NewIngressReferenceFinder(c client.Client) *IngressReferenceFinder {
	return &IngressReferenceFinder{
		Client: c,
	}
}

// FindSecretReferences finds all Ingresses that reference the given Secret
func (f *IngressReferenceFinder) FindSecretReferences(ctx context.Context, c client.Client, secretName, namespace string) ([]client.Object, error) {
	var results []client.Object

	ingressList := &networkingv1.IngressList{}
	if err := c.List(ctx, ingressList, client.InNamespace(namespace)); err != nil {
		return nil, err
	}

	for i := range ingressList.Items {
		ingress := &ingressList.Items[i]
		if f.ingressReferencesSecret(ingress, secretName) {
			results = append(results, ingress)
		}
	}

	return results, nil
}

// ingressReferencesSecret checks if an Ingress references the given secret
func (f *IngressReferenceFinder) ingressReferencesSecret(ingress *networkingv1.Ingress, secretName string) bool {
	// Check spec.tls[].secretName (for TLS secrets)
	for _, tls := range ingress.Spec.TLS {
		if tls.SecretName == secretName {
			return true
		}
	}

	return false
}

// Ingress does not reference ConfigMaps. This method is implemented to satisfy the ReferenceFinderStrategy interface.
func (f *IngressReferenceFinder) FindConfigMapReferences(ctx context.Context, c client.Client, configMapName, namespace string) ([]client.Object, error) {
	return nil, nil
}

// GetResourceType returns the Kubernetes resource type this strategy handles
func (f *IngressReferenceFinder) GetResourceType() string {
	return "Ingress"
}
