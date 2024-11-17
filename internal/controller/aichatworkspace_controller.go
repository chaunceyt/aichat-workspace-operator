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

	"github.com/pingcap/errors"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/api/resource"
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
		var result *ctrl.Result
		// ensureNamespace - ensure the namepace for the aichatworkspace is created.
		// should match the workspacename
		result, err = r.ensureNamespace(ctx, aichat, r.namespaceForAIChatWorkspace(aichat))
		if result != nil {
			return *result, err
		}
		// ensurePVC - ensure the persistentvolumeclaim for web files folder is managed.
		result, err = r.ensurePVC(ctx, aichat, r.pvcForOpenWebUI(aichat))
		if result != nil {
			return *result, err
		}

		/**
		  Create the following:
		  some derived from https://github.com/open-webui/open-webui/tree/main/kubernetes/manifest/base
		  - namespace
		  - resourcequota
		  - serviceaccount for ollama api
		  - serviceaccount for openwebui
		  - statefulset for Ollama API
		  - service for Ollama API
		  - deployment for Open WebUI
		  - service for Open WebUI
		  - ingress for Open WebUI
		**/

		// run ensure<Component> func(s) to manage the creation and reconcile of each component
		// namespace, serviceaccount, statefulset, deployment, service, ingress, pvc
		namespace := r.namespaceForAIChatWorkspace(aichat)
		resourcequota := r.resourceQuotaForAIChatWorkspace(aichat)
		fmt.Println(namespace, resourcequota)

		ollamaServiceAccount := r.serviceAccountForOllama(aichat)
		ollamaAPI := r.statefulsetForOllama(aichat)
		ollamaService := r.serviceForOllama(aichat)
		fmt.Println(ollamaServiceAccount, ollamaAPI, ollamaService)

		openwebuiServiceAccount := r.serviceAccountForOpenWebUI(aichat)
		openwebuiDeployment := r.deploymentForOpenWebUI(aichat)
		openwebuiService := r.serviceForOpenWebUI(aichat)
		openwebuiIngress := r.ingressForOpenWebUI(aichat)
		fmt.Println(openwebuiServiceAccount, openwebuiDeployment, openwebuiService, openwebuiIngress)

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
		// run ensure<Component> func(s) to manage the re-creation and reconcile of each component
		// namespace, serviceaccount, statefulset, deployment, service, ingress, pvc
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

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AIChatWorkspaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.AIChatWorkspace{}).
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

func (r *AIChatWorkspaceReconciler) namespaceForAIChatWorkspace(cr *appsv1alpha1.AIChatWorkspace) *corev1.Namespace {
	namespace := &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Namespace",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   cr.Spec.WorkspaceName,
			Labels: defaultLabels(cr, "namespace"),
		},
	}
	return namespace
}

func (r *AIChatWorkspaceReconciler) resourceQuotaForAIChatWorkspace(cr *appsv1alpha1.AIChatWorkspace) *corev1.ResourceQuota {
	resourcequota := &corev1.ResourceQuota{

		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.WorkspaceName,
			Namespace: cr.Spec.WorkspaceName,
			Labels:    defaultLabels(cr, "resourcequota"),
		},
	}
	return resourcequota
}

func (r *AIChatWorkspaceReconciler) serviceAccountForOllama(cr *appsv1alpha1.AIChatWorkspace) *corev1.ServiceAccount {
	sa := &corev1.ServiceAccount{

		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "ollama-sa"),
			Namespace: cr.Spec.WorkspaceName,
			Labels:    defaultLabels(cr, "sa"),
		},
	}
	return sa
}

func (r *AIChatWorkspaceReconciler) serviceAccountForOpenWebUI(cr *appsv1alpha1.AIChatWorkspace) *corev1.ServiceAccount {
	sa := &corev1.ServiceAccount{

		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "openwebui-sa"),
			Namespace: cr.Spec.WorkspaceName,
			Labels:    defaultLabels(cr, "sa"),
		},
	}
	return sa
}

func (r *AIChatWorkspaceReconciler) deploymentForOpenWebUI(cr *appsv1alpha1.AIChatWorkspace) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "openwebui"),
			Namespace: cr.Spec.WorkspaceName,
			Labels:    defaultLabels(cr, "deployment"),
		},
	}
	return deployment

}

func (r *AIChatWorkspaceReconciler) serviceForOpenWebUI(cr *appsv1alpha1.AIChatWorkspace) *corev1.Service {
	service := &corev1.Service{

		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "openwebui"),
			Namespace: cr.Spec.WorkspaceName,
			Labels:    defaultLabels(cr, "service"),
		},
	}
	return service
}

func (r *AIChatWorkspaceReconciler) serviceForOllama(cr *appsv1alpha1.AIChatWorkspace) *corev1.Service {
	service := &corev1.Service{

		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "ollama"),
			Namespace: cr.Spec.WorkspaceName,
			Labels:    defaultLabels(cr, "srv"),
		},
	}
	return service
}

func (r *AIChatWorkspaceReconciler) statefulsetForOllama(cr *appsv1alpha1.AIChatWorkspace) *appsv1.StatefulSet {
	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "ollama"),
			Namespace: cr.Spec.WorkspaceName,
			Labels:    defaultLabels(cr, "sts"),
		},
	}

	return sts
}

func (r *AIChatWorkspaceReconciler) pvcForOpenWebUI(cr *appsv1alpha1.AIChatWorkspace) *corev1.PersistentVolumeClaim {
	var storageSize = "2Gi"
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "openwebui-pvc"),
			Namespace: cr.Spec.WorkspaceName,
			Labels:    defaultLabels(cr, "pvc"),
		},
		Spec: corev1.PersistentVolumeClaimSpec{

			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},

			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(storageSize),
				},
			},
		},
	}
	controllerutil.SetControllerReference(cr, pvc, r.Scheme)
	return pvc
}

func (r *AIChatWorkspaceReconciler) ingressForOpenWebUI(cr *appsv1alpha1.AIChatWorkspace) *networkingv1.Ingress {
	ing := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName(cr, "openwebui"),
			Namespace: cr.Namespace,
			Labels:    defaultLabels(cr, "ingress"),
		},
	}

	return ing
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
