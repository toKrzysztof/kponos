package presentation

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"

	orphanagev1alpha1 "github.com/toKrzysztof/kponos/api/v1alpha1"
)

// StatusWriter handles writing status updates to OrphanagePolicy resources
type StatusWriter struct {
	client.Client
}

// NewStatusWriter creates a new StatusWriter
func NewStatusWriter(client client.Client) *StatusWriter {
	return &StatusWriter{
		Client: client,
	}
}

// UpdateStatus updates the status of an OrphanagePolicy
func (s *StatusWriter) UpdateStatus(ctx context.Context, policy *orphanagev1alpha1.OrphanagePolicy) error {
	// TODO: Implement status update logic
	// Update policy.Status with current reconciliation state
	return s.Status().Update(ctx, policy)
}

// UpdateStatusWithOrphanCount updates the status with orphan count information
func (s *StatusWriter) UpdateStatusWithOrphanCount(ctx context.Context, policy *orphanagev1alpha1.OrphanagePolicy, secretOrphanCount, configMapOrphanCount int) error {
	// TODO: Implement status update with orphan counts
	// Update policy.Status with the counts
	// Then call UpdateStatus
	
	return s.UpdateStatus(ctx, policy)
}

