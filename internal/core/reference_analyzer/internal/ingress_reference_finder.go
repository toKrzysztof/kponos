package internal

import (
	"context"

	networkingv1 "k8s.io/api/networking/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IngressReferenceFinder finds references to Secrets and ConfigMaps in Ingress resources
type IngressReferenceFinder struct {
	BaseReferenceFinder
}

// NewIngressReferenceFinder creates a new IngressReferenceFinder
func NewIngressReferenceFinder(c client.Client) *IngressReferenceFinder {
	return &IngressReferenceFinder{
		BaseReferenceFinder: BaseReferenceFinder{
			Client:       c,
			resourceType: "Ingress",
		},
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

// FindServiceReferences finds all Ingresses that reference the given Service
func (f *IngressReferenceFinder) FindServiceReferences(ctx context.Context, c client.Client, serviceName, namespace string) ([]client.Object, error) {
	var results []client.Object

	ingressList := &networkingv1.IngressList{}
	if err := c.List(ctx, ingressList, client.InNamespace(namespace)); err != nil {
		return nil, err
	}

	for i := range ingressList.Items {
		ingress := &ingressList.Items[i]
		if f.ingressReferencesService(ingress, serviceName) {
			results = append(results, ingress)
		}
	}

	return results, nil
}

// ingressReferencesService checks if an Ingress references the given service
func (f *IngressReferenceFinder) ingressReferencesService(ingress *networkingv1.Ingress, serviceName string) bool {
	// Check spec.defaultBackend.service.name
	if ingress.Spec.DefaultBackend != nil && ingress.Spec.DefaultBackend.Service != nil {
		if ingress.Spec.DefaultBackend.Service.Name == serviceName {
			return true
		}
	}

	// Check spec.rules[].http.paths[].backend.service.name
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP == nil {
			continue
		}
		for _, path := range rule.HTTP.Paths {
			if path.Backend.Service != nil && path.Backend.Service.Name == serviceName {
				return true
			}
		}
	}

	return false
}
