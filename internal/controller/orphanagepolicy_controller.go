/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	orphanagev1alpha1 "github.com/toKrzysztof/kponos/api/v1alpha1"
)

var log = logf.Log.WithName("controller_orphanagepolicy")

// OrphanagePolicyReconciler reconciles a OrphanagePolicy object
type OrphanagePolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=orphanage.kponos.io,resources=orphanagepolicies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=orphanage.kponos.io,resources=orphanagepolicies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=orphanage.kponos.io,resources=orphanagepolicies/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *OrphanagePolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.WithValues("orphanagepolicy", req.NamespacedName)
	logger.Info("Reconciling OrphanagePolicy")

	// Fetch the OrphanagePolicy instance
	policy := &orphanagev1alpha1.OrphanagePolicy{}
	if err := r.Get(ctx, req.NamespacedName, policy); err != nil {
		logger.Error(err, "unable to fetch OrphanagePolicy")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Monitor only the resource types specified in the spec
	for _, resourceType := range policy.Spec.ResourceTypes {
		switch resourceType {
		case orphanagev1alpha1.ResourceTypeSecret:
			secretList := &corev1.SecretList{}
			if err := r.List(ctx, secretList, client.InNamespace("default")); err != nil {
				logger.Error(err, "unable to list Secrets")
				continue
			}
			logger.Info("Found Secrets", "count", len(secretList.Items))
			for _, secret := range secretList.Items {
				logger.Info("Secret", "name", secret.Name, "type", secret.Type)
			}

		case orphanagev1alpha1.ResourceTypeConfigMap:
			configMapList := &corev1.ConfigMapList{}
			if err := r.List(ctx, configMapList, client.InNamespace("default")); err != nil {
				logger.Error(err, "unable to list ConfigMaps")
				continue
			}
			logger.Info("Found ConfigMaps", "count", len(configMapList.Items))
			for _, cm := range configMapList.Items {
				logger.Info("ConfigMap", "name", cm.Name, "dataKeys", len(cm.Data))
			}

		default:
			logger.Info("Unknown resource type", "type", resourceType)
		}
	}

	return ctrl.Result{}, nil
}

// mapToOrphanagePolicy maps Secret/ConfigMap events to reconcile all OrphanagePolicy objects
func (r *OrphanagePolicyReconciler) mapToOrphanagePolicy(ctx context.Context, obj client.Object) []reconcile.Request {
	policyList := &orphanagev1alpha1.OrphanagePolicyList{}
	if err := r.List(ctx, policyList); err != nil {
		log.Error(err, "unable to list OrphanagePolicy objects")
		return []reconcile.Request{}
	}

	requests := make([]reconcile.Request, 0, len(policyList.Items))

	// Determine the resource type of the object that triggered this
	var resourceType orphanagev1alpha1.ResourceType
	switch obj.(type) {
	case *corev1.Secret:
		resourceType = orphanagev1alpha1.ResourceTypeSecret
	case *corev1.ConfigMap:
		resourceType = orphanagev1alpha1.ResourceTypeConfigMap
	default:
		return []reconcile.Request{}
	}

	// Only enqueue policies that are watching this resource type
	for _, policy := range policyList.Items {
		for _, rt := range policy.Spec.ResourceTypes {
			if rt == resourceType {
				requests = append(requests, reconcile.Request{
					NamespacedName: types.NamespacedName{
						Name:      policy.Name,
						Namespace: policy.Namespace,
					},
				})
				break
			}
		}
	}

	return requests
}

// SetupWithManager sets up the controller with the Manager.
func (r *OrphanagePolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&orphanagev1alpha1.OrphanagePolicy{}).
		Named("orphanagepolicy").
		Watches(&corev1.Secret{}, handler.EnqueueRequestsFromMapFunc(r.mapToOrphanagePolicy)).
		Watches(&corev1.ConfigMap{}, handler.EnqueueRequestsFromMapFunc(r.mapToOrphanagePolicy)).
		Complete(r)
}
