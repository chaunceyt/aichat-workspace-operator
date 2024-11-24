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

	ollamaRunning := r.isOllamaUp(ctx, instance)
	if !ollamaRunning {
		// If Ollama isn't running yet, requeue the ctrl
		// to run again after a delay
		delay := time.Second * time.Duration(5)

		logger.Info(fmt.Sprintf("Ollama isn't running, waiting for %s", delay))

		return &ctrl.Result{RequeueAfter: delay}, nil
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
