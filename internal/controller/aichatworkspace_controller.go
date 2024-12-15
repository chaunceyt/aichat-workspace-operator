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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/chaunceyt/aichat-workspace-operator/api/v1alpha1"
	"github.com/chaunceyt/aichat-workspace-operator/internal/constants"

	"github.com/go-logr/logr"
)

const (
	ReconcileErrorInterval       = 10 * time.Second
	ReconcileSuccessInterval     = 30 * time.Second
	reconcileStarted             = "staring reconcile"
	aichatWorkspaceFinalizerName = "core.aichatworkspace.io/finalizer"
)

// AIChatWorkspaceReconciler reconciles a AIChatWorkspace object
type AIChatWorkspaceReconciler struct {
	client.Client
	kubeconfig *restclient.Config
	Scheme     *runtime.Scheme
	Recorder   record.EventRecorder
}

type AIChatWorkspaceInstance struct {
	r                     *AIChatWorkspaceReconciler
	req                   ctrl.Request
	ctx                   context.Context
	aichatWorkspaceConfig *appsv1alpha1.AIChatWorkspace
	logger                logr.Logger
}

type AIChatWorkspace interface {
	execute(*AIChatWorkspaceInstance) (ctrl.Result, error)
	setNext(AIChatWorkspace)
}

var (
	aichatWorkspaceControllerLog = ctrl.Log.WithName(constants.ManagedBy)
)

// +kubebuilder:rbac:groups=apps.aichatworkspaces.io,resources=aichatworkspaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.aichatworkspaces.io,resources=aichatworkspaces/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.aichatworkspaces.io,resources=aichatworkspaces/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments;statefulsets,verbs=*
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=*
// +kubebuilder:rbac:groups="",resources=namespaces;pods;services;persistentvolumeclaims;serviceaccounts;resourcequotas,verbs=*
// +kubebuilder:rbac:groups="",resources=events,verbs=create
// +kubebuilder:rbac:groups="metrics.k8s.io",resources=pods,verbs=get;watch;list
// +kubebuilder:rbac:groups="http.keda.sh",resources=httpscaledobjects,verbs=*

/**
 * Reconciles an AIChatWorkspace object by executing a series of steps in order.
 *
 * Reconcile is part of the main kubernetes reconciliation loop which aims to
 * move the current state of the cluster closer to the desired state.
 *
 * @param ctx The context for the reconciliation request.
 * @param req The request to reconcile, containing the namespace and name of the AIChatWorkspace object.
 * @return A ctrl.Result indicating whether the reconciliation was successful or if it should be retried.
 *
 * For more details, check Reconcile and its Result here:
 * - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
 */
func (r *AIChatWorkspaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	now := time.Now()
	logger := log.FromContext(ctx, "ns", req.NamespacedName.Namespace, "cr", req.NamespacedName.Name).WithValues("starting reconcile of aichatworkspace", time.Since(now))

	logger.Info("starting reconciling aichatworkspace")

	instance := AIChatWorkspaceInstance{
		r:      r,
		req:    req,
		ctx:    ctx,
		logger: aichatWorkspaceControllerLog,
	}

	initStep := InitAIChatWorkspaceStep{}
	finalizerStep := FinalizerStep{}
	createStep := CreateAIChatWorkspaceStep{}
	updateStep := UpdateAIChatWorkspaceStep{}
	deleteStep := DeleteAIChatWorkspaceStep{}

	initStep.setNext(&finalizerStep)
	finalizerStep.setNext(&createStep)
	createStep.setNext(&updateStep)
	updateStep.setNext(&deleteStep)

	return initStep.execute(&instance)
}

/**
 * Sets up the AIChatWorkspaceReconciler with the provided manager.
 *
 * This function is responsible for setting up the controller's dependencies and
 * registering it with the manager. It also defines the resources that the
 * controller owns, which includes Namespaces, ResourceQuotas, ServiceAccounts,
 * PersistentVolumeClaims, Services, Deployments, and StatefulSets.
 *
 * @param mgr The manager to set up the controller with.
 * @return An error if there is an issue setting up the controller, or nil otherwise.
 */
func (r *AIChatWorkspaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.AIChatWorkspace{}).
		Owns(&corev1.Namespace{}).
		Owns(&corev1.ResourceQuota{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Owns(&corev1.Service{}).
		Owns(&appsv1.Deployment{}).
		Owns(&appsv1.StatefulSet{}).
		Named(constants.AIChatWorkspaceName).
		Complete(r)
}

/**
 * Finishes the reconciliation process by determining if it was successful or not.
 *
 * If an error occurred during the reconciliation, this function will return a
 * ctrl.Result with Requeue set to true and RequeueAfter set to ReconcileErrorInterval.
 * If no error occurred, but requeueImmediate is true, it will return a ctrl.Result
 * with Requeue set to true and RequeueAfter set to 0. Otherwise, it will return a
 * ctrl.Result with Requeue set to true and RequeueAfter set to ReconcileSuccessInterval.
 *
 * @param err The error that occurred during reconciliation, if any.
 * @param requeueImmediate Whether or not the controller should be immediately requeued.
 * @return A ctrl.Result indicating whether the reconciliation was successful or if it should be retried.
 */
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

/**
 * Patches the status of an AIChatWorkspace object.
 *
 * This function retrieves the latest version of the AIChatWorkspace object from the Kubernetes API,
 * and then patches its status using the provided object. It returns an error if there is a problem
 * retrieving or patching the object.
 *
 * @param ctx The context for the request to the Kubernetes API.
 * @param aichat The AIChatWorkspace object to be patched.
 * @return An error if there was a problem patching the object, or nil otherwise.
 */
func (r *AIChatWorkspaceReconciler) patchStatus(ctx context.Context, aichat *appsv1alpha1.AIChatWorkspace) error {
	key := client.ObjectKeyFromObject(aichat)
	latest := &appsv1alpha1.AIChatWorkspace{}
	if err := r.Client.Get(ctx, key, latest); err != nil {
		return err
	}
	return r.Client.Status().Patch(ctx, aichat, client.MergeFrom(latest))
}

/**
 * Returns the default labels for an AIChatWorkspace object.
 *
 * The function generates a set of default labels based on the provided namespace and name,
 * as well as the component. These labels are used to identify the resources created by
 * the AIChatWorkspace controller.
 *
 * @param namespace The namespace where the AIChatWorkspace is running.
 * @param name The name of the AIChatWorkspace object.
 * @param component The component that these labels belong to (e.g., "aichat-workspace").
 * @return A map of default labels for an AIChatWorkspace object.
 */
func defaultLabels(namespace string, name string, component string) map[string]string {
	partOf := fmt.Sprintf("aichat-workspace-%s", namespace)

	return map[string]string{
		"app.kubernetes.io/name":       name,
		"app.kubernetes.io/part-of":    partOf,
		"app.kubernetes.io/component":  component,
		"app.kubernetes.io/version":    constants.Version,
		"app.kubernetes.io/managed-by": constants.ManagedBy,
		"aichatworkspace":              namespace,
	}
}

/**
 * Deletes an AIChatWorkspace object and its associated resources.
 *
 * This function takes a context, logger, and instance as parameters. It deletes the
 * namespace with the given name from the spec of the provided instance. If any error
 * occurs during this process, it returns that error; otherwise, it logs information about
 * deleting the AIChatWorkspace object and returns nil.
 *
 * @param ctx The context for the request to the Kubernetes API.
 * @param instance The AIChatWorkspace object to be deleted.
 * @return An error if there is a problem deleting the resources associated with the AIChatWorkspace object, or nil otherwise.
 */
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

/**
 * Generates a name by combining the workspace and name.
 *
 * This function takes two parameters: the workspace and the name. It combines these
 * two strings with a hyphen in between to generate a unique name.
 *
 * @param workspace The workspace string to be combined with the name.
 * @param name The name string to be combined with the workspace.
 * @return A generated name that is a combination of the workspace and name.
 */
func generateName(workspace, name string) string {
	return fmt.Sprintf("%s-%s", workspace, name)
}
