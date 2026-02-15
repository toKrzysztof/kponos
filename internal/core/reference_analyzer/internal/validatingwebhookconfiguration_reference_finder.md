# ValidatingWebhookConfigurationReferenceFinder Documentation

## Overview

The `ValidatingWebhookConfigurationReferenceFinder` is a component that analyzes Kubernetes ValidatingWebhookConfiguration resources to find static references to Services. ValidatingWebhookConfigurations reference Services to configure admission webhooks that validate Kubernetes resources before they are persisted.

## Static Reference Types Analyzed

### Service References

The finder detects references to Services in the following ValidatingWebhookConfiguration specification locations:

1. **Webhook Client Configuration**
   - `webhooks[].clientConfig.service.name` - The name of the Service that hosts the validating webhook
   - `webhooks[].clientConfig.service.namespace` - The namespace of the Service (must match the target Service namespace)

## Notes

- The finder performs **static analysis** of ValidatingWebhookConfiguration resource specifications. It does not detect dynamic references or references created at runtime.
- ValidatingWebhookConfiguration is a **cluster-scoped** resource, so searches are performed across all namespaces, but Service references include namespace information.
- The finder returns all matching ValidatingWebhookConfiguration resources that reference the given Service by name and namespace.

