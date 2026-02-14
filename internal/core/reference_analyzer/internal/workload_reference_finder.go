package internal

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// WorkloadResourceType represents a valid workload resource type
type WorkloadResourceType string

const (
	WorkloadResourceTypePod        WorkloadResourceType = "Pod"
	WorkloadResourceTypeDeployment WorkloadResourceType = "Deployment"
	WorkloadResourceTypeStatefulSet WorkloadResourceType = "StatefulSet"
	WorkloadResourceTypeDaemonSet   WorkloadResourceType = "DaemonSet"
)

// WorkloadReferenceFinder finds references to Secrets and ConfigMaps in workload resources
// that are Pods or create Pods (Deployment, StatefulSet, DaemonSet)
type WorkloadReferenceFinder struct {
	client.Client
	resourceType WorkloadResourceType
}

// NewWorkloadReferenceFinder creates a new WorkloadReferenceFinder for the given resource type
func NewWorkloadReferenceFinder(client client.Client, resourceType WorkloadResourceType) *WorkloadReferenceFinder {
	return &WorkloadReferenceFinder{
		Client:       client,
		resourceType: resourceType,
	}
}

// FindSecretReferences finds all resources that reference the given Secret
func (f *WorkloadReferenceFinder) FindSecretReferences(ctx context.Context, c client.Client, secretName, namespace string) ([]client.Object, error) {
	var results []client.Object

	switch f.resourceType {
	case WorkloadResourceTypePod:
		podList := &corev1.PodList{}
		if err := c.List(ctx, podList, client.InNamespace(namespace)); err != nil {
			return nil, err
		}
		for i := range podList.Items {
			pod := &podList.Items[i]
			if f.podSpecReferencesSecret(&pod.Spec, secretName) {
				results = append(results, pod)
			}
		}

	case WorkloadResourceTypeDeployment:
		deploymentList := &appsv1.DeploymentList{}
		if err := c.List(ctx, deploymentList, client.InNamespace(namespace)); err != nil {
			return nil, err
		}
		for i := range deploymentList.Items {
			deployment := &deploymentList.Items[i]
			if f.podSpecReferencesSecret(&deployment.Spec.Template.Spec, secretName) {
				results = append(results, deployment)
			}
		}

	case WorkloadResourceTypeStatefulSet:
		statefulSetList := &appsv1.StatefulSetList{}
		if err := c.List(ctx, statefulSetList, client.InNamespace(namespace)); err != nil {
			return nil, err
		}
		for i := range statefulSetList.Items {
			statefulSet := &statefulSetList.Items[i]
			if f.podSpecReferencesSecret(&statefulSet.Spec.Template.Spec, secretName) {
				results = append(results, statefulSet)
			}
		}

	case WorkloadResourceTypeDaemonSet:
		daemonSetList := &appsv1.DaemonSetList{}
		if err := c.List(ctx, daemonSetList, client.InNamespace(namespace)); err != nil {
			return nil, err
		}
		for i := range daemonSetList.Items {
			daemonSet := &daemonSetList.Items[i]
			if f.podSpecReferencesSecret(&daemonSet.Spec.Template.Spec, secretName) {
				results = append(results, daemonSet)
			}
		}
	}

	return results, nil
}

// podSpecReferencesSecret checks if a PodSpec references the given secret
func (f *WorkloadReferenceFinder) podSpecReferencesSecret(podSpec *corev1.PodSpec, secretName string) bool {
	// Check volumes[].secret.secretName
	for _, volume := range podSpec.Volumes {
		if volume.Secret != nil && volume.Secret.SecretName == secretName {
			return true
		}
	}

	// Check containers[].envFrom[].secretRef.name and containers[].env[].valueFrom.secretKeyRef.name
	for _, container := range podSpec.Containers {
		for _, envFrom := range container.EnvFrom {
			if envFrom.SecretRef != nil && envFrom.SecretRef.Name == secretName {
				return true
			}
		}
		for _, env := range container.Env {
			if env.ValueFrom != nil && env.ValueFrom.SecretKeyRef != nil && env.ValueFrom.SecretKeyRef.Name == secretName {
				return true
			}
		}
	}

	// Check initContainers[].envFrom[].secretRef.name and initContainers[].env[].valueFrom.secretKeyRef.name
	for _, container := range podSpec.InitContainers {
		for _, envFrom := range container.EnvFrom {
			if envFrom.SecretRef != nil && envFrom.SecretRef.Name == secretName {
				return true
			}
		}
		for _, env := range container.Env {
			if env.ValueFrom != nil && env.ValueFrom.SecretKeyRef != nil && env.ValueFrom.SecretKeyRef.Name == secretName {
				return true
			}
		}
	}

	// Check imagePullSecrets[].name
	for _, imagePullSecret := range podSpec.ImagePullSecrets {
		if imagePullSecret.Name == secretName {
			return true
		}
	}

	return false
}

// FindConfigMapReferences finds all resources that reference the given ConfigMap
func (f *WorkloadReferenceFinder) FindConfigMapReferences(ctx context.Context, c client.Client, configMapName, namespace string) ([]client.Object, error) {
	var results []client.Object

	switch f.resourceType {
	case WorkloadResourceTypePod:
		podList := &corev1.PodList{}
		if err := c.List(ctx, podList, client.InNamespace(namespace)); err != nil {
			return nil, err
		}
		for i := range podList.Items {
			pod := &podList.Items[i]
			if f.podSpecReferencesConfigMap(&pod.Spec, configMapName) {
				results = append(results, pod)
			}
		}

	case WorkloadResourceTypeDeployment:
		deploymentList := &appsv1.DeploymentList{}
		if err := c.List(ctx, deploymentList, client.InNamespace(namespace)); err != nil {
			return nil, err
		}
		for i := range deploymentList.Items {
			deployment := &deploymentList.Items[i]
			if f.podSpecReferencesConfigMap(&deployment.Spec.Template.Spec, configMapName) {
				results = append(results, deployment)
			}
		}

	case WorkloadResourceTypeStatefulSet:
		statefulSetList := &appsv1.StatefulSetList{}
		if err := c.List(ctx, statefulSetList, client.InNamespace(namespace)); err != nil {
			return nil, err
		}
		for i := range statefulSetList.Items {
			statefulSet := &statefulSetList.Items[i]
			if f.podSpecReferencesConfigMap(&statefulSet.Spec.Template.Spec, configMapName) {
				results = append(results, statefulSet)
			}
		}

	case WorkloadResourceTypeDaemonSet:
		daemonSetList := &appsv1.DaemonSetList{}
		if err := c.List(ctx, daemonSetList, client.InNamespace(namespace)); err != nil {
			return nil, err
		}
		for i := range daemonSetList.Items {
			daemonSet := &daemonSetList.Items[i]
			if f.podSpecReferencesConfigMap(&daemonSet.Spec.Template.Spec, configMapName) {
				results = append(results, daemonSet)
			}
		}
	}

	return results, nil
}

// podSpecReferencesConfigMap checks if a PodSpec references the given configmap
func (f *WorkloadReferenceFinder) podSpecReferencesConfigMap(podSpec *corev1.PodSpec, configMapName string) bool {
	// Check volumes[].configMap.name
	for _, volume := range podSpec.Volumes {
		if volume.ConfigMap != nil && volume.ConfigMap.Name == configMapName {
			return true
		}
	}

	// Check containers[].envFrom[].configMapRef.name and containers[].env[].valueFrom.configMapKeyRef.name
	for _, container := range podSpec.Containers {
		for _, envFrom := range container.EnvFrom {
			if envFrom.ConfigMapRef != nil && envFrom.ConfigMapRef.Name == configMapName {
				return true
			}
		}
		for _, env := range container.Env {
			if env.ValueFrom != nil && env.ValueFrom.ConfigMapKeyRef != nil && env.ValueFrom.ConfigMapKeyRef.Name == configMapName {
				return true
			}
		}
	}

	// Check initContainers[].envFrom[].configMapRef.name and initContainers[].env[].valueFrom.configMapKeyRef.name
	for _, container := range podSpec.InitContainers {
		for _, envFrom := range container.EnvFrom {
			if envFrom.ConfigMapRef != nil && envFrom.ConfigMapRef.Name == configMapName {
				return true
			}
		}
		for _, env := range container.Env {
			if env.ValueFrom != nil && env.ValueFrom.ConfigMapKeyRef != nil && env.ValueFrom.ConfigMapKeyRef.Name == configMapName {
				return true
			}
		}
	}

	return false
}

// GetResourceType returns the Kubernetes resource type this strategy handles
func (f *WorkloadReferenceFinder) GetResourceType() string {
	return string(f.resourceType)
}