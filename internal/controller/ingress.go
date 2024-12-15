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

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appsv1alpha1 "github.com/chaunceyt/aichat-workspace-operator/api/v1alpha1"
)

// ensureIngress ensures that the specified ingress resource exists in the cluster.
// If it does not exist, it will be created. If an error occurs during this process,
// it will be logged and returned.
func (r *AIChatWorkspaceReconciler) ensureIngress(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, ing *networkingv1.Ingress) (*reconcile.Result, error) {
	logger := log.FromContext(ctx)
	found := &networkingv1.Ingress{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      ing.Name,
		Namespace: instance.Spec.WorkspaceName,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the ingress
		logger.Info("Creating a new Ingress", "Ingress.Namespace", instance.Spec.WorkspaceName, "Ingress.Name", ing.Name)
		err = r.Create(context.TODO(), ing)

		if err != nil {
			// Creation failed
			logger.Error(err, "Failed to create new Ingress", "Ingress.Namespace", instance.Spec.WorkspaceName, "Ingress.Name", ing.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the ingress not existing
		logger.Error(err, "Failed to get Ingress")
		return &reconcile.Result{}, err
	}

	return nil, nil
}
