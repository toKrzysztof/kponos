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

// OrphanagePolicyReconciler reconciles an OrphanagePolicy object
type OrphanagePolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *OrphanagePolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.WithValues("orphanagepolicy", req.NamespacedName)
	logger.Info("Reconciling OrphanagePolicy")

	policy := &orphanagev1alpha1.OrphanagePolicy{}
	if err := r.Get(ctx, req.NamespacedName, policy); err != nil {
		logger.Error(err, "unable to fetch OrphanagePolicy")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	secretList := &corev1.SecretList{}
	if err := r.List(ctx, secretList, client.InNamespace(req.Namespace)); err != nil {
		logger.Error(err, "unable to list Secrets")
	}

	configMapList := &corev1.ConfigMapList{}
	if err := r.List(ctx, configMapList, client.InNamespace(req.Namespace)); err != nil {
		logger.Error(err, "unable to list ConfigMaps")
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
