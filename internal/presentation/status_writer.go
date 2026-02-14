package presentation

import (
	"context"
	"time"

	orphanagev1alpha1 "github.com/toKrzysztof/kponos/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// StatusWriter handles writing status updates to OrphanagePolicy resources
type StatusWriter struct {
	client.Client
}

// NewStatusWriter creates a new StatusWriter
func NewStatusWriter(c client.Client) *StatusWriter {
	return &StatusWriter{
		Client: c,
	}
}

// UpdateStatus updates the status of an OrphanagePolicy
func (s *StatusWriter) UpdateStatus(ctx context.Context, policy *orphanagev1alpha1.OrphanagePolicy, orphans []client.Object) error {
	now := time.Now()

	policy.Status.OrphanCount = len(orphans)
	policy.Status.LastChanged = metav1.NewTime(now)
	policy.Status.Orphans = make([]orphanagev1alpha1.Orphan, len(orphans))
	for i, orphan := range orphans {
		policy.Status.Orphans[i].Kind = orphan.GetObjectKind().GroupVersionKind().Kind
		policy.Status.Orphans[i].Name = orphan.GetName()
	}

	return s.Status().Update(ctx, policy)
}