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
	"github.com/chaunceyt/aichat-workspace-operator/internal/config"
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
	now := time.Now()
	logger := log.FromContext(ctx, "ns", req.NamespacedName.Namespace, "cr", req.NamespacedName.Name).WithValues("starting reconcile of aichatworkspace", time.Since(now))

	logger.Info("starting reconciling aichatworkspace")

	restConfig := ctrl.GetConfigOrDie()
	ctrlClient, err := client.New(restConfig, client.Options{})
	if err != nil {
		r.finishReconcile(err, false)
	}

	config, err := config.GetConfigFromConfigMap(ctx, ctrlClient, "aichat-workspace-operator-system")
	if err != nil {
		r.finishReconcile(err, false)
	}
	fmt.Println("config: ", config)

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

// SetupWithManager sets up the controller with the Manager.
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

// deleteAIChatWorkspace is responsible for cleaning up the resources created for the aichat workspace.
// deleting the namespace ensures each of the objects created are deleted.
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
