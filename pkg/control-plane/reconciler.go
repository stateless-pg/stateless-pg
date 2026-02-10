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

package controlplane

import (
	"context"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"

	corev1alpha1 "github.com/stateless-pg/stateless-pg/pkg/api/v1alpha1"
)

// +kubebuilder:rbac:groups=stateless-pg.com,resources=tenants,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=stateless-pg.com,resources=tenants/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=stateless-pg.com,resources=tenants/finalizers,verbs=update

// Reconcile handles the reconciliation loop for Tenant resources
func (r *Operator) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	err := r.sync(ctx, req.Name, req.Namespace)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("Failed to sync tenant %s/%s: %w", req.Namespace, req.Name, err)
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Operator) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.Tenant{}).
		Named("tenant").
		Complete(r)
}
