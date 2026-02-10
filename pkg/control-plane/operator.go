package controlplane

import (
	"fmt"
	"log/slog"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
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