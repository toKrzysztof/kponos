# ServiceReferenceFinder Documentation

## Overview

The `ServiceReferenceFinder` is a component that analyzes Kubernetes Pod resources to find Services whose label selectors match the Pod. This implements the Pod-to-Service relationship where Services use label selectors to determine which Pods they route traffic to.

## Static Reference Types Analyzed

### Service References (via Label Selector)

The finder detects Services whose selectors match a Pod's labels:

1. **Service Selector Matching**
   - `spec.selector` - A map of labels in the Service that must match the Pod's labels
   - The finder uses `FindPodReferences` to find all Services in the namespace whose selectors match the given Pod's labels

## Notes

- The finder performs **static analysis** of Pod and Service resources. It uses label matching to find Services that would select the given Pod.
- All searches are scoped to a specific namespace (Services and Pods must be in the same namespace for the selector to work).
- The `FindPodReferences` method takes a Pod name and namespace, then returns all Services whose selectors match that Pod's labels.
