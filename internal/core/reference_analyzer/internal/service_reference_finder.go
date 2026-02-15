package internal

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ServiceReferenceFinder finds references to Pods in Service resources
type ServiceReferenceFinder struct {
	BaseReferenceFinder
}

// NewServiceReferenceFinder creates a new ServiceReferenceFinder
func NewServiceReferenceFinder(c client.Client) *ServiceReferenceFinder {
	return &ServiceReferenceFinder{
		BaseReferenceFinder: BaseReferenceFinder{
			Client:       c,
			resourceType: "Service",
		},
	}
}

// FindPodReferences finds all Services whose selectors match the given Pod
func (f *ServiceReferenceFinder) FindPodReferences(ctx context.Context, c client.Client, podName, namespace string) ([]client.Object, error) {
	// First, get the Pod to extract its labels
	pod := &corev1.Pod{}
	if err := c.Get(ctx, client.ObjectKey{Name: podName, Namespace: namespace}, pod); err != nil {
		return nil, err
	}

	// Get the Pod's labels
	podLabels := pod.Labels

	// List all Services in the namespace
	serviceList := &corev1.ServiceList{}
	if err := c.List(ctx, serviceList, client.InNamespace(namespace)); err != nil {
		return nil, err
	}

	// Filter Services whose selectors match the Pod's labels
	var results []client.Object
	for i := range serviceList.Items {
		service := &serviceList.Items[i]

		// Skip Services without selectors (they don't match any Pods)
		if len(service.Spec.Selector) == 0 {
			continue
		}

		selector := labels.SelectorFromSet(service.Spec.Selector)

		// Check if the selector matches the Pod's labels
		if selector.Matches(labels.Set(podLabels)) {
			results = append(results, service)
		}
	}

	return results, nil
}
