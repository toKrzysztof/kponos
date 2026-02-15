# APIServiceReferenceFinder Documentation

## Overview

The `APIServiceReferenceFinder` is a component that analyzes Kubernetes APIService resources to find static references to Services. APIServices reference Services to configure extension API servers that extend the Kubernetes API.

## Static Reference Types Analyzed

### Service References

The finder detects references to Services in the following APIService specification locations:

1. **Service Configuration**
   - `spec.service.name` - The name of the Service that hosts the extension API server
   - `spec.service.namespace` - The namespace of the Service (must match the target Service namespace)

### Secret References

APIServices do not reference Secrets. The `FindSecretReferences` method is implemented to satisfy the `ReferenceFinderStrategy` interface but always returns an empty result.

### ConfigMap References

APIServices do not reference ConfigMaps. The `FindConfigMapReferences` method is implemented to satisfy the `ReferenceFinderStrategy` interface but always returns an empty result.

## Notes

- The finder performs **static analysis** of APIService resource specifications. It does not detect dynamic references or references created at runtime.
- APIService is a **cluster-scoped** resource, so searches are performed across all namespaces, but Service references include namespace information.
- The finder returns all matching APIService resources that reference the given Service by name and namespace.

