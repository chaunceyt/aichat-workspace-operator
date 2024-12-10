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

func (r *AIChatWorkspaceReconciler) ensureNamespace(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, ns *corev1.Namespace) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)

	found := &corev1.Namespace{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name: ns.Name,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating the namespace", "instance.Spec.Namespace", instance.Spec.WorkspaceName)

		controllerutil.SetControllerReference(instance, ns, r.Scheme)
		err = r.Create(context.TODO(), ns)
		if err != nil {
			logger.Error(err, "Failed to create namespace", "instance.Spec.Namespace", instance.Spec.WorkspaceName)
			return &ctrl.Result{}, err
		}

		return nil, nil

	} else if err != nil {
		logger.Error(err, "Failed to get namespace")

		return &ctrl.Result{}, err
	}

	return nil, nil
}
