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
	"strconv"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/chaunceyt/aichat-workspace-operator/api/v1alpha1"
	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/k8s"
	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/ollama"
	"github.com/chaunceyt/aichat-workspace-operator/internal/constants"
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
		logger.Info("Creating a StatefulSet", "Name", sts.Name, "WorkspaceName", instance.Spec.WorkspaceName)

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
			logger.Error(
				err,
				"Checking to see Model exists.",
				"Model Name", llm,
				"Workspace Name", instance.Spec.WorkspaceName,
			)
			return &ctrl.Result{}, err
		}
		if !ok {
			fmt.Printf("The %s LLM does not exist. Starting the ollama pull ...\n", llm)
			err = ollama.PullModel(llm, ollamaServerURI)
			if err != nil {
				logger.Error(
					err,
					"Pulling down model",
					"ModelName", llm,
					"WorkspaceName", instance.Spec.WorkspaceName,
					"Name", sts.Name,
				)
				return &ctrl.Result{}, err
			}
			ollama.CreateFromModelFile(llm, ollamaServerURI, instance.Spec.Patterns)
		}
	}

	models, err := ollama.ListRunningModels(ollamaServerURI)
	if err != nil {
		return &ctrl.Result{}, err
	}

	ollamaWorkloadInfo, err := r.getPodInfo(ctx, instance, "ollama")
	if err != nil {
		return &ctrl.Result{}, err
	}

	openwebuiWorkloadInfo, err := r.getPodInfo(ctx, instance, "openwebui")
	if err != nil {
		return &ctrl.Result{}, err
	}

	instance.Status.OllamaWorkload = ollamaWorkloadInfo
	instance.Status.OpenWebUIWorkload = openwebuiWorkloadInfo
	instance.Status.ModelsInUse = models
	err = r.Status().Update(context.TODO(), instance)
	if err != nil {
		logger.Error(
			err,
			"Updating Status for OllamaPod",
			"WorkspaceName", instance.Spec.WorkspaceName,
		)
		return &ctrl.Result{}, err
	}

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

func (r *AIChatWorkspaceReconciler) getPodInfo(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, workload string) (string, error) {
	logger := log.FromContext(ctx)
	var err error
	var podInfo string
	var podName string
	var containerName string
	var cpu int64
	var memory int64
	var diskusage string

	ollamaPod := fmt.Sprintf("%s-%s", instance.Spec.WorkspaceName, workload)
	lbs := map[string]string{
		"app.kubernetes.io/name": ollamaPod,
	}
	labelSelector := labels.SelectorFromSet(lbs)
	listOps := &client.ListOptions{Namespace: instance.Spec.WorkspaceName, LabelSelector: labelSelector}

	podList := &corev1.PodList{}
	if err = r.List(context.TODO(), podList, listOps); err != nil {
		return podInfo, err
	}

	var duCommand string
	var mountPath string
	var defaultStorage string

	switch workload {
	case "ollama":
		duCommand = duSHCommand(constants.OllamaVolumeMountPath)
		mountPath = constants.OllamaVolumeMountPath
		defaultStorage = fmt.Sprintf("used of %s allocated", constants.OllamaDefaultVolumeSize)
	case "openwebui":
		duCommand = duSHCommand(constants.OpenwebuiVolumeMountPath)
		mountPath = constants.OpenwebuiVolumeMountPath
		defaultStorage = fmt.Sprintf("used of %s allocated", constants.OpenwebuiDefaultVolumeSize)
	}

	for _, pod := range podList.Items {
		podName = pod.Name
		duResults, _, err := k8s.ExecuteRemoteCommand(&pod, duCommand)
		if err != nil {
			return podInfo, err
		}

		duResults = strings.Replace(duResults, "\t", " ", -1)
		duResults = strings.Replace(duResults, "\r\n", " ", -1)
		duResults = strings.Replace(duResults, mountPath, defaultStorage, -1)

		diskusage = duResults
	}

	// get metrics
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		logger.Error(
			err,
			"Setting up restConfig",
		)
		return podInfo, err
	}

	metricsClient, err := metricsv.NewForConfig(restConfig)
	if err != nil {
		logger.Error(
			err,
			"Setting up metrics client",
		)
		return podInfo, err
	}

	podMetricsList, err := metricsClient.MetricsV1beta1().PodMetricses(instance.Spec.WorkspaceName).List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=" + ollamaPod,
	})
	if err != nil {
		logger.Error(
			err,
			"Getting Pod Metriceses",
			"WorkspaceName", instance.Spec.WorkspaceName,
		)
		return podInfo, err
	}

	for _, podMetrics := range podMetricsList.Items {
		for _, container := range podMetrics.Containers {
			containerName += container.Name
			cpu += container.Usage.Cpu().MilliValue()
			memory += container.Usage.Memory().MilliValue() / 1000000000
		}
	}

	// PodInfo that will be available when describing a aichatworkspace
	podInfo = fmt.Sprintf("Name: %s\n\tContainer: %s\n\tStorage: %s \n\tCPU (millicores): %sm\n\tMemory (bytes): %sMi", podName, containerName, diskusage, strconv.FormatInt(cpu, 10), strconv.FormatInt(memory, 10))

	return podInfo, nil
}

func duSHCommand(mountPath string) string {
	return fmt.Sprintf("du -sh %s", mountPath)
}
