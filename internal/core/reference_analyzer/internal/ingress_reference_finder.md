# IngressReferenceFinder Documentation

## Overview

The `IngressReferenceFinder` is a component that analyzes Kubernetes Ingress resources to find static references to Secrets and Services. Ingresses reference Secrets for TLS/SSL certificate configuration to enable HTTPS traffic, and reference Services as backends to route traffic to.

## Static Reference Types Analyzed

### Secret References

The finder detects references to Secrets in the following Ingress specification locations:

1. **TLS Configuration**
   - `spec.tls[].secretName` - Secrets containing TLS certificates and keys used for HTTPS termination on the Ingress

### Service References

The finder detects references to Services in the following Ingress specification locations:

1. **Default Backend**
   - `spec.defaultBackend.service.name` - The Service name used as the default backend when no rules match

2. **Path Backends**
   - `spec.rules[].http.paths[].backend.service.name` - The Service name used as the backend for specific path rules

### ConfigMap References

Ingresses do not reference ConfigMaps. The `FindConfigMapReferences` method is implemented to satisfy the `ReferenceFinderStrategy` interface but always returns an empty result.

## Notes

- The finder performs **static analysis** of Ingress resource specifications. It does not detect dynamic references or references created at runtime.
- All searches are scoped to a specific namespace.
- The finder returns all matching Ingress resources that reference the given Secret or Service by name.