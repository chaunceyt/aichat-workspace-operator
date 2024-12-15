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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kedahttpv1alpha1 "github.com/kedacore/http-add-on/operator/apis/http/v1alpha1"

	appsv1alpha1 "github.com/chaunceyt/aichat-workspace-operator/api/v1alpha1"
)

/**
 * Ensures that a HTTPScaledObject exists for the given AIChatWorkspace instance.
 *
 * If the HTTPScaledObject does not exist, it will be created. If an error occurs
 * while trying to create or retrieve the HTTPScaledObject, the function will return
 * an error.
 *
 * @param ctx The context in which the function is being executed.
 * @param instance The AIChatWorkspace instance for which to ensure a HTTPScaledObject exists.
 * @param httpso The desired state of the HTTPScaledObject.
 * @return A ctrl.Result and an error, or nil if no further reconciliation is needed.
 */
func (r *AIChatWorkspaceReconciler) ensureHTTPScaledObject(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, httpso *kedahttpv1alpha1.HTTPScaledObject) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)
	scheme := runtime.NewScheme()
	_ = kedahttpv1alpha1.AddToScheme(scheme)
	found := &kedahttpv1alpha1.HTTPScaledObject{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      httpso.Name,
		Namespace: instance.Spec.WorkspaceName,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a HTTPScaledObject", "HTTPScaledObject.Namespace", instance.Spec.WorkspaceName, "HTTPScaledObject.Name", httpso.Name)
		controllerutil.SetControllerReference(instance, httpso, r.Scheme)
		err = r.Create(context.TODO(), httpso)

		if err != nil {
			logger.Error(err, "Failed to createHTTPScaledObject", "HTTPScaledObject.Namespace", instance.Spec.WorkspaceName, "HTTPScaledObject.Name", httpso.Name)
			return &ctrl.Result{}, err
		}

		return nil, nil

	} else if err != nil {
		logger.Error(err, "Failed to get HTTPScaledObject")

		return &ctrl.Result{}, err
	}

	return nil, nil
}
