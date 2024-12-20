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
	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/ollama"
)

// ensureStatefulSet ensures the Ollama service is created and running as a StatefulSet.
/**
 * This function checks if the given StatefulSet exists in the cluster.
 * If it does not, it creates a new one with the provided instance and returns nil.
 * If an error occurs during this process, it logs the error and returns a Result.
 */
func (r *AIChatWorkspaceReconciler) ensureStatefulSet(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, sts *appsv1.StatefulSet) (*ctrl.Result, error) {
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
		logger.Info("Creating a new StatefulSet", "StatefulSet.Namespace", instance.Spec.WorkspaceName, "StatefulSet.Name", sts.Name)

		// Set the controller reference for the StatefulSet
		controllerutil.SetControllerReference(instance, sts, r.Scheme)
		err = r.Create(context.TODO(), sts)
		if err != nil {
			logger.Error(err, "Failed to create new StatefulSet", "StatefulSet.Namespace", instance.Spec.WorkspaceName, "StatefulSet.Name", sts.Name)

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

	// ensure ollama is running.
	// it needs to be running in order to pull in the instance.Spec.Models
	ollamaRunning := r.isOllamaUp(ctx, instance)
	if !ollamaRunning {
		delay := time.Second * time.Duration(5)
		logger.Info(fmt.Sprintf("Ollama isn't running, waiting for %s", delay))

		return &ctrl.Result{RequeueAfter: delay}, nil
	}

	// ensure the instance.Spec.Models are available.
	serviceName := fmt.Sprintf("%s-ollama", instance.Spec.WorkspaceName)
	ollamaPort := int64(11434)
	ollamaServerURI := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", serviceName, instance.Spec.WorkspaceName, ollamaPort)

	for _, llm := range instance.Spec.Models {
		ok, err := ollama.DoesModelExist(llm, ollamaServerURI)
		if err != nil {
			return &ctrl.Result{}, err
		}
		if !ok {
			fmt.Printf("The %s LLM does not exist. Starting the ollama pull ...\n", llm)
			err = ollama.PullModel(llm, ollamaServerURI)
			if err != nil {
				logger.Error(err, "Failed to pull Model", "ModelName", llm, "StatefulSet.Namespace", instance.Spec.WorkspaceName, "StatefulSet.Name", sts.Name)
				return &ctrl.Result{}, err
			}
			ollama.CreateFromModelFile(llm, ollamaServerURI, instance.Spec.Patterns)
		}
	}

	models, err := ollama.ListRunningModels(ollamaServerURI)
	if err != nil {
		return &ctrl.Result{}, err
	}
	fmt.Println("Models Running: ", models)

	return nil, nil
}

// Returns whether or not the ollama StatefulSet is running
func (r *AIChatWorkspaceReconciler) isOllamaUp(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace) bool {
	logger := log.FromContext(ctx)
	sts := &appsv1.StatefulSet{}
	ollamaName := generateName(instance.Spec.WorkspaceName, "ollama")

	err := r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      ollamaName,
		Namespace: instance.Spec.WorkspaceName,
	}, sts)

	if err != nil {
		logger.Error(err, "StatefulSet for Ollama not found")
		return false
	}

	if sts.Status.ReadyReplicas == 1 {
		return true
	}

	return false
}
