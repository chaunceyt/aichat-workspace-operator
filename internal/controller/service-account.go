/*
Copyright 2024 AIChatWorkspace Contributors.

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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/chaunceyt/aichat-workspace-operator/api/v1alpha1"
)

/**
 * Ensures the existence of a ServiceAccount in the specified namespace.
 *
 * If the ServiceAccount does not exist, it will be created. If it already exists,
 * this function will return without modifying it.
 *
 * @param ctx The context for the request.
 * @param instance The AIChatWorkspace instance that owns the ServiceAccount.
 * @param sa The ServiceAccount to ensure existence of.
 * @return A ctrl.Result and an error, or nil if no action is required.
 */
func (r *AIChatWorkspaceReconciler) ensureServiceAccount(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, sa *corev1.ServiceAccount) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)

	found := &corev1.ServiceAccount{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      sa.Name,
		Namespace: instance.Spec.WorkspaceName,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a ServiceAccount", "ServiceAccount.Namespace", instance.Spec.WorkspaceName, "ServiceAccount.Name", sa.Name)

		controllerutil.SetControllerReference(instance, sa, r.Scheme)
		err = r.Create(context.TODO(), sa)

		if err != nil {
			logger.Error(err, "Failed to create ServiceAccount", "ServiceAccount.Namespace", instance.Spec.WorkspaceName, "ServiceAccount.Name", sa.Name)

			return &ctrl.Result{}, err
		}

		return nil, nil

	} else if err != nil {
		logger.Error(err, "Failed to get ServiceAccount")

		return &ctrl.Result{}, err
	}

	return nil, nil
}
