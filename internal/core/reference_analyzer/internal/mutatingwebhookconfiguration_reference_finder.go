package internal

import (
	"context"

	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MutatingWebhookConfigurationReferenceFinder finds references to Services in MutatingWebhookConfiguration resources
type MutatingWebhookConfigurationReferenceFinder struct {
	BaseReferenceFinder
}

// NewMutatingWebhookConfigurationReferenceFinder creates a new MutatingWebhookConfigurationReferenceFinder
func NewMutatingWebhookConfigurationReferenceFinder(c client.Client) *MutatingWebhookConfigurationReferenceFinder {
	return &MutatingWebhookConfigurationReferenceFinder{
		BaseReferenceFinder: BaseReferenceFinder{
			Client:       c,
			resourceType: "MutatingWebhookConfiguration",
		},
	}
}

// FindServiceReferences finds all MutatingWebhookConfigurations that reference the given Service
func (f *MutatingWebhookConfigurationReferenceFinder) FindServiceReferences(ctx context.Context, c client.Client, serviceName, namespace string) ([]client.Object, error) {
	var results []client.Object

	webhookList := &admissionregistrationv1.MutatingWebhookConfigurationList{}
	// MutatingWebhookConfiguration is cluster-scoped, so we don't filter by namespace
	if err := c.List(ctx, webhookList); err != nil {
		return nil, err
	}

	for i := range webhookList.Items {
		webhookConfig := &webhookList.Items[i]
		if f.mutatingWebhookConfigurationReferencesService(webhookConfig, serviceName, namespace) {
			results = append(results, webhookConfig)
		}
	}

	return results, nil
}

// mutatingWebhookConfigurationReferencesService checks if a MutatingWebhookConfiguration references the given service
func (f *MutatingWebhookConfigurationReferenceFinder) mutatingWebhookConfigurationReferencesService(webhookConfig *admissionregistrationv1.MutatingWebhookConfiguration, serviceName, namespace string) bool {
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
