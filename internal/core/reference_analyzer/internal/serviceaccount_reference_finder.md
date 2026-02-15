# ServiceAccountReferenceFinder Documentation

## Overview

The `ServiceAccountReferenceFinder` is a component that analyzes Kubernetes ServiceAccount resources to find static references to Secrets. ServiceAccounts can reference Secrets for image pull authentication and for mounting secrets into Pods that use the ServiceAccount.

## Static Reference Types Analyzed

### Secret References

The finder detects references to Secrets in the following ServiceAccount specification locations:

1. **Secrets List**
   - `secrets[].name` - Secrets that are automatically mounted into Pods using this ServiceAccount, or used for image pull authentication

2. **Image Pull Secrets**
   - `imagePullSecrets[].name` - Secrets used for pulling container images from private registries when Pods use this ServiceAccount

## Notes

- The finder performs **static analysis** of ServiceAccount resource specifications. It does not detect dynamic references or references created at runtime.
- All searches are scoped to a specific namespace.
- The finder returns all matching ServiceAccount resources that reference the given Secret by name.