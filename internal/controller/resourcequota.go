package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/chaunceyt/aichat-workspace-operator/api/v1alpha1"
)

func (r *AIChatWorkspaceReconciler) ensureResourceQuota(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, rq *corev1.ResourceQuota) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)

	found := &corev1.ResourceQuota{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      rq.Name,
		Namespace: instance.Spec.WorkspaceName,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a resource quota", "ResourceQuota.Namespace", instance.Spec.WorkspaceName, "ResourceQuota.Name", rq.Name)

		controllerutil.SetControllerReference(instance, rq, r.Scheme)
		err = r.Create(context.TODO(), rq)

		if err != nil {
			logger.Error(err, "Failed to create resource quota", "ResourceQuota.Namespace", instance.Spec.WorkspaceName, "ResourceQuota.Name", rq.Name)
			return &ctrl.Result{}, err
		}

		return nil, nil

	} else if err != nil {
		logger.Error(err, "Failed to get resource quota")

		return &ctrl.Result{}, err
	}

	return nil, nil
}
