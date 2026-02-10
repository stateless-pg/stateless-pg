package controlplane

import (
	"context"
	"fmt"
	"log/slog"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
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
	tenant := &corev1alpha1.Tenant{}
	if err := o.nclient.Get(ctx, client.ObjectKey{
		Name:      name,
		Namespace: namespace,
	}, tenant); err != nil {
		if apierrors.IsNotFound(err) {
			// Tenant resource not found, could have been deleted after reconcile request.
			// Return and don't requeue
			return nil
		}
		return err
	}

	tenant = tenant.DeepCopy()

	key := fmt.Sprintf("%s/%s", namespace, name)
	logger := o.logger.With("key", key)

	logger.Info("Sync tenant")

	// TODO: Add tenant reconciliation logic here

	logger.Info("Tenant synced successfully")
	return nil
}