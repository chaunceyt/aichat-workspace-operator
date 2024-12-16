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

// ensurePVC ensures that a Persistent Volume Claim (PVC) exists for the given AIChatWorkspace instance.
//
// It checks if the PVC already exists in the specified namespace. If it does not exist, it creates a new PVC and sets
// the controller reference to the AIChatWorkspace instance. If an error occurs during this process, it logs the error
// and returns it.
func (r *AIChatWorkspaceReconciler) ensurePVC(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, pvc *corev1.PersistentVolumeClaim) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)

	found := &corev1.PersistentVolumeClaim{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      pvc.Name,
		Namespace: instance.Spec.WorkspaceName,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a new PVC", "PVC.Namespace", instance.Spec.WorkspaceName, "PVC.Name", pvc.Name)
		controllerutil.SetControllerReference(instance, pvc, r.Scheme)
		err = r.Create(context.TODO(), pvc)

		if err != nil {
			logger.Error(err, "Failed to create new PVC", "PVC.Namespace", instance.Spec.WorkspaceName, "PVC.Name", pvc.Name)
			return &ctrl.Result{}, err
		}

		return nil, nil

	} else if err != nil {
		logger.Error(err, "Failed to get PVC")

		return &ctrl.Result{}, err
	}

	return nil, nil
}
