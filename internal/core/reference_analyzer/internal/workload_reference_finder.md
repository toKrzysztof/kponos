# WorkloadReferenceFinder Documentation

## Overview

The `WorkloadReferenceFinder` is a component that analyzes Kubernetes workload resources to find static references to Secrets and ConfigMaps. It searches through Pod specifications (either directly in Pods or embedded in Pod templates within workload controllers) to identify all references to these resources.

## Supported Resource Types

The finder supports the following Kubernetes workload resource types:

- **Pod** - Direct Pod resources
- **Deployment** - Deployment resources (analyzes the Pod template)
- **StatefulSet** - StatefulSet resources (analyzes the Pod template)
- **DaemonSet** - DaemonSet resources (analyzes the Pod template)

## Static Reference Types Analyzed

### Secret References

The finder detects references to Secrets in the following PodSpec locations:

1. **Volume Mounts**

   - `spec.volumes[].secret.secretName` - Secret volumes mounted in the Pod
2. **Container Environment Variables (Regular Containers)**

   - `spec.containers[].envFrom[].secretRef.name` - Secrets loaded as environment variables via `envFrom`
   - `spec.containers[].env[].valueFrom.secretKeyRef.name` - Individual secret keys referenced in environment variables
3. **Init Container Environment Variables**

   - `spec.initContainers[].envFrom[].secretRef.name` - Secrets loaded as environment variables in init containers via `envFrom`
   - `spec.initContainers[].env[].valueFrom.secretKeyRef.name` - Individual secret keys referenced in init container environment variables
4. **Image Pull Secrets**

   - `spec.imagePullSecrets[].name` - Secrets used for pulling container images from private registries

### ConfigMap References

The finder detects references to ConfigMaps in the following PodSpec locations:

1. **Volume Mounts**

   - `spec.volumes[].configMap.name` - ConfigMap volumes mounted in the Pod
2. **Container Environment Variables (Regular Containers)**

   - `spec.containers[].envFrom[].configMapRef.name` - ConfigMaps loaded as environment variables via `envFrom`
   - `spec.containers[].env[].valueFrom.configMapKeyRef.name` - Individual ConfigMap keys referenced in environment variables
3. **Init Container Environment Variables**

   - `spec.initContainers[].envFrom[].configMapRef.name` - ConfigMaps loaded as environment variables in init containers via `envFrom`
   - `spec.initContainers[].env[].valueFrom.configMapKeyRef.name` - Individual ConfigMap keys referenced in init container environment variables

## Notes

- The finder performs **static analysis** of resource specifications. It does not detect dynamic references or references created at runtime.
- For Deployment, StatefulSet, and DaemonSet resources, the finder analyzes the Pod template (`spec.template.spec`) rather than the top-level resource specification.
- All searches are scoped to a specific namespace.
- The finder returns all matching resources that reference the given Secret or ConfigMap by name.
