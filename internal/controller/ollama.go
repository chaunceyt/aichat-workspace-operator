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
func (r *AIChatWorkspaceReconciler) ensureStatefulSet(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, sts *appsv1.StatefulSet) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)
	found := &appsv1.StatefulSet{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      sts.Name,
		Namespace: instance.Spec.WorkspaceName,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		// Create the StatefulSet
		logger.Info("Creating a new StatefulSet", "StatefulSet.Namespace", instance.Spec.WorkspaceName, "StatefulSet.Name", sts.Name)
		controllerutil.SetControllerReference(instance, sts, r.Scheme)
		err = r.Create(context.TODO(), sts)

		if err != nil {
			// Creation failed
			logger.Error(err, "Failed to create new StatefulSet", "StatefulSet.Namespace", instance.Spec.WorkspaceName, "StatefulSet.Name", sts.Name)
			return &ctrl.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the StatefulSet not existing
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

	for _, llm := range instance.Spec.Models {
		ollamaServerURI := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", serviceName, instance.Spec.WorkspaceName, ollamaPort)
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
		}
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
