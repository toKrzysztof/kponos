package internal

import (
	"context"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ValidatingWebhookConfigurationReferenceFinder finds references to Services in ValidatingWebhookConfiguration resources
type ValidatingWebhookConfigurationReferenceFinder struct {
	BaseReferenceFinder
}

// NewValidatingWebhookConfigurationReferenceFinder creates a new ValidatingWebhookConfigurationReferenceFinder
func NewValidatingWebhookConfigurationReferenceFinder(c client.Client) *ValidatingWebhookConfigurationReferenceFinder {
	return &ValidatingWebhookConfigurationReferenceFinder{
		BaseReferenceFinder: BaseReferenceFinder{
			Client:       c,
			resourceType: "ValidatingWebhookConfiguration",
		},
	}
}

// FindServiceReferences finds all ValidatingWebhookConfigurations that reference the given Service
func (f *ValidatingWebhookConfigurationReferenceFinder) FindServiceReferences(ctx context.Context, c client.Client, serviceName, namespace string) ([]client.Object, error) {
	var results []client.Object

	webhookList := &admissionregistrationv1.ValidatingWebhookConfigurationList{}
	// ValidatingWebhookConfiguration is cluster-scoped, so we don't filter by namespace
	if err := c.List(ctx, webhookList); err != nil {
		return nil, err
	}

	for i := range webhookList.Items {
		webhookConfig := &webhookList.Items[i]
		if f.validatingWebhookConfigurationReferencesService(webhookConfig, serviceName, namespace) {
			results = append(results, webhookConfig)
		}
	}

	return results, nil
}

// validatingWebhookConfigurationReferencesService checks if a ValidatingWebhookConfiguration references the given service
func (f *ValidatingWebhookConfigurationReferenceFinder) validatingWebhookConfigurationReferencesService(webhookConfig *admissionregistrationv1.ValidatingWebhookConfiguration, serviceName, namespace string) bool {
	// Check webhooks[].clientConfig.service.name and webhooks[].clientConfig.service.namespace
	for _, webhook := range webhookConfig.Webhooks {
		if webhook.ClientConfig.Service != nil {
			if webhook.ClientConfig.Service.Name == serviceName && webhook.ClientConfig.Service.Namespace == namespace {
				return true
			}
		}
	}

	return false
}
