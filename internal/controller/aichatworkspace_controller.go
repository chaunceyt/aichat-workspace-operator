/*
Copyright 2024.

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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/chaunceyt/aichat-workspace-operator/api/v1alpha1"
)

const (
	ReconcileErrorInterval   = 10 * time.Second
	ReconcileSuccessInterval = 30 * time.Second
)

// AIChatWorkspaceReconciler reconciles a AIChatWorkspace object
type AIChatWorkspaceReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=apps.aichatworkspaces.io,resources=aichatworkspaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.aichatworkspaces.io,resources=aichatworkspaces/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.aichatworkspaces.io,resources=aichatworkspaces/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments;statefulsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=namespaces;pods;services;persistentvolumeclaims;serviceaccounts;events,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AIChatWorkspace object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *AIChatWorkspaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var err error

	logger.Info("starting reconciling aichatworkspace")

	aichat := &appsv1alpha1.AIChatWorkspace{}
	err = r.Get(context.TODO(), req.NamespacedName, aichat)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	aichatWorkspaceFinalizerName := "core.aichatworkspace.io/finalizer"

	isCreated := aichat.Status.IsCreated
	pendingDeletion := aichat.ObjectMeta.DeletionTimestamp != nil
	hasFinalizer := controllerutil.ContainsFinalizer(aichat, aichatWorkspaceFinalizerName)

	switch {
	case !hasFinalizer && !pendingDeletion:
		controllerutil.AddFinalizer(aichat, aichatWorkspaceFinalizerName)
		if err = r.Update(ctx, aichat); err != nil {
			return r.finishReconcile(err, false)
		}
	case !isCreated && !pendingDeletion:
		logger.Info("reconciling aichat", "aichat", aichat, "action", "create")

		// handleReconcile - creates the aichatworkspace.
		var result *ctrl.Result
		result, err = r.handleReconcile(ctx, result, aichat)
		if result != nil {
			return r.finishReconcile(err, true)
		}

		aichat.Status.IsCreated = true

		if err = r.Status().Update(ctx, aichat); err != nil {
			apimeta.SetStatusCondition(&aichat.Status.Conditions, metav1.Condition{
				Status:             metav1.ConditionFalse,
				Reason:             appsv1alpha1.ReconciliationFailedReason,
				Message:            err.Error(),
				Type:               appsv1alpha1.ConditionTypeReady,
				ObservedGeneration: aichat.GetGeneration(),
			})
			if err = r.patchStatus(ctx, aichat); err != nil {
				err = fmt.Errorf("unable to patch status after progressing: %w", err)
				return ctrl.Result{Requeue: true}, err
			}
		}
		apimeta.SetStatusCondition(&aichat.Status.Conditions, metav1.Condition{
			Status:             metav1.ConditionFalse,
			Reason:             appsv1alpha1.ProgressingReason,
			Message:            "Reconciliation progressing",
			Type:               appsv1alpha1.ConditionTypeReady,
			ObservedGeneration: aichat.GetGeneration(),
		})
		apimeta.SetStatusCondition(&aichat.Status.Conditions, metav1.Condition{
			Status:             metav1.ConditionTrue,
			Reason:             appsv1alpha1.ReconciliationSucceededReason,
			Message:            "AIChatWorkspace reconciled",
			Type:               appsv1alpha1.ConditionTypeReady,
			ObservedGeneration: aichat.GetGeneration(),
		})

		if err = r.patchStatus(ctx, aichat); err != nil {
			err = fmt.Errorf("unable to patch status after progressing: %w", err)
			return ctrl.Result{Requeue: true}, err
		}
		r.Recorder.Event(aichat, "Normal", "Created",
			fmt.Sprintf("aichatWorkspace %s was created in namespace %s",
				aichat.Name,
				aichat.Namespace))
	case !isCreated && pendingDeletion:
		logger.Info("reconciling aichat", "aichat", aichat, "action", "no-op")
	case isCreated && !pendingDeletion:
		logger.Info("reconciling aichat", "aichat", aichat, "action", "update")

		// handleReconcile - reconciles changes to the aichatworkspace.
		var result *ctrl.Result
		result, err = r.handleReconcile(ctx, result, aichat)
		if result != nil {
			return r.finishReconcile(err, true)
		}

	case isCreated && pendingDeletion:
		logger.Info("reconciling aichat", "aichat", aichat, "action", "delete")
		if controllerutil.ContainsFinalizer(aichat, aichatWorkspaceFinalizerName) {
			if err = r.deleteAIChatWorkspace(ctx, aichat); err != nil {
				return ctrl.Result{}, err
			}
			r.Recorder.Event(aichat, "Warning", "Deleting",
				fmt.Sprintf("aichatWorkspace %s is being deleted from the namespace %s",
					aichat.Name,
					aichat.Namespace))
			controllerutil.RemoveFinalizer(aichat, aichatWorkspaceFinalizerName)
			if err = r.Update(ctx, aichat); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return r.finishReconcile(err, false)
}

// SetupWithManager sets up the controller with the Manager.
func (r *AIChatWorkspaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.AIChatWorkspace{}).
		Owns(&corev1.Namespace{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Owns(&appsv1.Deployment{}).
		Owns(&appsv1.StatefulSet{}).
		Named("aichatworkspace").
		Complete(r)
}

func (r *AIChatWorkspaceReconciler) finishReconcile(err error, requeueImmediate bool) (ctrl.Result, error) {
	if err != nil {
		interval := ReconcileErrorInterval
		if requeueImmediate {
			interval = 0
		}
		return ctrl.Result{Requeue: true, RequeueAfter: interval}, err
	}
	interval := ReconcileSuccessInterval
	if requeueImmediate {
		interval = 0
	}
	return ctrl.Result{Requeue: true, RequeueAfter: interval}, nil
}

func (r *AIChatWorkspaceReconciler) patchStatus(ctx context.Context, aichat *appsv1alpha1.AIChatWorkspace) error {
	key := client.ObjectKeyFromObject(aichat)
	latest := &appsv1alpha1.AIChatWorkspace{}
	if err := r.Client.Get(ctx, key, latest); err != nil {
		return err
	}
	return r.Client.Status().Patch(ctx, aichat, client.MergeFrom(latest))
}

func createInt32(x int32) *int32 {
	return &x
}

func defaultLabels(cr *appsv1alpha1.AIChatWorkspace, component string) map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":      cr.Name,
		"app.kubernetes.io/part-of":   "part-of",
		"app.kubernetes.io/component": component,
		"app.kubernetes.io/version":   "version",
		"release":                     "release",
		"provider":                    "aichatworkspace-operator",
	}
}

func workloadName(cr *appsv1alpha1.AIChatWorkspace, workloadType string) string {
	return cr.Spec.WorkspaceName + "-" + workloadType
}

// deleteAIChatWorkspace is responsible for cleaning up the resources created for the aichat workspace.
// deleting the namespace ensure each of the objects created was deleted.
//   - resourcequota
//   - serviceaccount for ollama api
//   - serviceaccount for openwebui
//   - statefulset for Ollama API
//   - service for Ollama API
//   - deployment for Open WebUI
//   - service for Open WebUI
//   - ingress for Open WebUI
func (r *AIChatWorkspaceReconciler) deleteAIChatWorkspace(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace) error {
	logger := log.FromContext(ctx)

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: instance.Spec.WorkspaceName,
		},
	}
	err := r.Delete(context.TODO(), namespace)
	if err != nil {
		return err
	}

	logger.Info("deleted aichatworkspace", "aichatworkspace", instance.Spec.WorkspaceName, "action", "deleted")

	return nil
}

func generateName(workspace, name string) string {
	return fmt.Sprintf("%s-%s", workspace, name)
}
