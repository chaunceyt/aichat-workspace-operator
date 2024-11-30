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

func (r *AIChatWorkspaceReconciler) ensureService(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, svc *corev1.Service) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)

	found := &corev1.Service{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      svc.Name,
		Namespace: instance.Spec.WorkspaceName,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a Service", "Service.Namespace", instance.Spec.WorkspaceName, "Service.Name", svc.Name)

		controllerutil.SetControllerReference(instance, svc, r.Scheme)
		err = r.Create(context.TODO(), svc)

		if err != nil {
			logger.Error(err, "Failed to create Service", "Service.Namespace", instance.Spec.WorkspaceName, "Service.Name", svc.Name)

			return &ctrl.Result{}, err
		}

		return nil, nil

	} else if err != nil {
		logger.Error(err, "Failed to get Service")

		return &ctrl.Result{}, err
	}

	return nil, nil
}
