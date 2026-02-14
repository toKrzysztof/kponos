package resourceHandler

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ResourceReferenceFinder func(context.Context, string, string) ([]client.Object, error)
