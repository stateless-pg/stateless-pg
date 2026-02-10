package controlplane

import (
	"context"
	"fmt"
	"log/slog"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1alpha1 "github.com/stateless-pg/stateless-pg/pkg/api/v1alpha1"
)

const (
	controllerName = "control-plane"
)

// Operator manages lifecycle for control-plane resources.
type Operator struct {
	nclient client.Client
	kclient kubernetes.Interface
	scheme  *runtime.Scheme
	logger  *slog.Logger
}

// New creates a new Controller.
func New(client client.Client, scheme *runtime.Scheme, logger *slog.Logger, config *rest.Config) (*Operator, error) {
	logger = logger.With("component", controllerName)

	// Create kubernetes clientset for direct client-go operations
	kclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	return &Operator{
		logger:  logger,
		nclient: client,
		kclient: kclient,
		scheme:  scheme,
	}, nil
}

// sync performs the actual reconciliation logic for a tenant
func (o *Operator) sync(ctx context.Context, name, namespace string) error {
	logger := o.logger.With("tenant", name, "namespace", namespace)
	logger.Info("Syncing tenant")

	// Fetch the Tenant resource
	tenant := &corev1alpha1.Tenant{}
	key := client.ObjectKey{Name: name, Namespace: namespace}
	if err := o.nclient.Get(ctx, key, tenant); err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	logger.Info("Tenant synced successfully", "generation", tenant.Spec.Generation, "placementPolicy", tenant.Spec.PlacementPolicy)
	return nil
}