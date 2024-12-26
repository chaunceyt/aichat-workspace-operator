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
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/chaunceyt/aichat-workspace-operator/api/v1alpha1"
)

/**
 * This function checks if the given StatefulSet exists in the cluster.
 * If it does not, it creates a new one with the provided instance and returns nil.
 * If an error occurs during this process, it logs the error and returns a Result.
 */
func (r *AIChatWorkspaceReconciler) ensurePostgres(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, sts *appsv1.StatefulSet) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)

	found := &appsv1.StatefulSet{}

	// Check if the StatefulSet already exists
	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      sts.Name,
		Namespace: instance.Spec.WorkspaceName,
	}, found)
	if err != nil && errors.IsNotFound(err) {
		/**
		 * If the StatefulSet does not exist, create a new one.
		 */
		logger.Info("Creating a Ollama StatefulSet", "Name", sts.Name, "WorkspaceName", instance.Spec.WorkspaceName)

		// Set the controller reference for the StatefulSet
		controllerutil.SetControllerReference(instance, sts, r.Scheme)
		err = r.Create(context.TODO(), sts)
		if err != nil {
			logger.Error(
				err,
				"Failed to create StatefulSet",
				"Name", sts.Name,
				"WorkspaceName", instance.Spec.WorkspaceName,
			)

			return &ctrl.Result{}, err
		}

		return nil, nil
	} else if err != nil {
		/**
		 * If an error occurs during the process of checking or creating the StatefulSet, log it and return a Result.
		 */
		logger.Error(err, "Failed to get StatefulSet")

		return &ctrl.Result{}, err
	}

	// ensure postgres is running.
	postgresRunning := r.isPostgresUp(ctx, instance)
	if !postgresRunning {
		delay := time.Second * time.Duration(5)
		logger.Info(fmt.Sprintf("Postgres isn't running, waiting for %s", delay))

		return &ctrl.Result{RequeueAfter: delay}, nil
	}

	return nil, nil
}

// Returns whether or not the ollama StatefulSet is running
func (r *AIChatWorkspaceReconciler) isPostgresUp(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace) bool {
	logger := log.FromContext(ctx)
	sts := &appsv1.StatefulSet{}
	postgresName := "postgres"

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      postgresName,
		Namespace: instance.Spec.WorkspaceName,
	}, sts)

	if err != nil {
		logger.Error(err, "StatefulSet for Postgres not found")
		return false
	}

	if sts.Status.ReadyReplicas == 1 {
		return true
	}

	return false
}
