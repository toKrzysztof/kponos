# MutatingWebhookConfigurationReferenceFinder Documentation

## Overview

The `MutatingWebhookConfigurationReferenceFinder` is a component that analyzes Kubernetes MutatingWebhookConfiguration resources to find static references to Services. MutatingWebhookConfigurations reference Services to configure admission webhooks that modify Kubernetes resources before they are persisted.

## Static Reference Types Analyzed

### Service References

The finder detects references to Services in the following MutatingWebhookConfiguration specification locations:

1. **Webhook Client Configuration**
   - `webhooks[].clientConfig.service.name` - The name of the Service that hosts the mutating webhook
   - `webhooks[].clientConfig.service.namespace` - The namespace of the Service (must match the target Service namespace)

### Secret References

MutatingWebhookConfigurations do not reference Secrets. The `FindSecretReferences` method is implemented to satisfy the `ReferenceFinderStrategy` interface but always returns an empty result.

### ConfigMap References

MutatingWebhookConfigurations do not reference ConfigMaps. The `FindConfigMapReferences` method is implemented to satisfy the `ReferenceFinderStrategy` interface but always returns an empty result.

## Notes

- The finder performs **static analysis** of MutatingWebhookConfiguration resource specifications. It does not detect dynamic references or references created at runtime.
- MutatingWebhookConfiguration is a **cluster-scoped** resource, so searches are performed across all namespaces, but Service references include namespace information.
- The finder returns all matching MutatingWebhookConfiguration resources that reference the given Service by name and namespace.

