# CustomResourceDefinitionReferenceFinder Documentation

## Overview

The `CustomResourceDefinitionReferenceFinder` is a component that analyzes Kubernetes CustomResourceDefinition resources to find static references to Services. CustomResourceDefinitions can reference Services when using webhook-based conversion to convert between different versions of custom resources.

## Static Reference Types Analyzed

### Service References

The finder detects references to Services in the following CustomResourceDefinition specification locations:

1. **Webhook Conversion Configuration**
   - `spec.conversion.webhook.clientConfig.service.name` - The name of the Service that hosts the conversion webhook
   - `spec.conversion.webhook.clientConfig.service.namespace` - The namespace of the Service (must match the target Service namespace)

## Notes

- The finder performs **static analysis** of CustomResourceDefinition resource specifications. It does not detect dynamic references or references created at runtime.
- CustomResourceDefinition is a **cluster-scoped** resource, so searches are performed across all namespaces, but Service references include namespace information.
- The finder returns all matching CustomResourceDefinition resources that reference the given Service by name and namespace.

