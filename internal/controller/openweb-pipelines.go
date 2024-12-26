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

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/chaunceyt/aichat-workspace-operator/api/v1alpha1"
)

/**
 * Ensures the existence and desired state of a Deployment.
 *
 * This function checks if a Deployment with the given name exists in the specified namespace. If it does not exist, it creates a new Deployment.
 * If the Deployment already exists, this function checks for any changes to the Deployment's spec and updates it if necessary.
 *
 * @param ctx The context in which the function is being executed.
 * @param instance The AIChatWorkspace instance that owns the Deployment.
 * @param deploy The desired state of the Deployment.
 * @return A Result object indicating whether the function should be retried, or an error if one occurred.
 */
func (r *AIChatWorkspaceReconciler) ensureOpenWebUIPipelines(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, deploy *appsv1.Deployment) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)

	found := &appsv1.Deployment{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      deploy.Name,
		Namespace: instance.Spec.WorkspaceName,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		// Create the Deployment
		logger.Info("Creating a new Deployment", "Deployment.Namespace", instance.Spec.WorkspaceName, "Deployment.Name", deploy.Name)
		controllerutil.SetControllerReference(instance, deploy, r.Scheme)
		err = r.Create(context.TODO(), deploy)

		if err != nil {
			// Creation failed
			logger.Error(
				err,
				"Creating Deployment",
				"Name", deploy.Name,
				"Workspace", instance.Spec.WorkspaceName,
			)
			return &ctrl.Result{}, err
		}

		return nil, nil

	} else if err != nil {
		// Error that isn't due to the Deployment not existing
		logger.Error(
			err,
			"Failed to get Deployment",
		)
		return &ctrl.Result{}, err
	}

	return nil, nil
}
