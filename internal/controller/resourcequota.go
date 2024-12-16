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
 * Ensures a resource quota exists for the given AIChatWorkspace instance.
 *
 * If the resource quota does not exist, it will be created. If it already exists,
 * this function will do nothing.
 *
 * @param ctx The context in which to perform the operation.
 * @param instance The AIChatWorkspace instance for which to ensure a resource quota.
 * @param rq The desired resource quota to create or check.
 * @return A ctrl.Result indicating whether the reconciliation was successful, and an error if one occurred.
 */
func (r *AIChatWorkspaceReconciler) ensureResourceQuota(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, rq *corev1.ResourceQuota) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)

	found := &corev1.ResourceQuota{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      rq.Name,
		Namespace: instance.Spec.WorkspaceName,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a resource quota", "ResourceQuota.Namespace", instance.Spec.WorkspaceName, "ResourceQuota.Name", rq.Name)

		controllerutil.SetControllerReference(instance, rq, r.Scheme)
		err = r.Create(context.TODO(), rq)

		if err != nil {
			logger.Error(err, "Failed to create resource quota", "ResourceQuota.Namespace", instance.Spec.WorkspaceName, "ResourceQuota.Name", rq.Name)
			return &ctrl.Result{}, err
		}

		return nil, nil

	} else if err != nil {
		logger.Error(err, "Failed to get resource quota")

		return &ctrl.Result{}, err
	}

	return nil, nil
}
