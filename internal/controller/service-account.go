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

func (r *AIChatWorkspaceReconciler) ensureServiceAccount(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, sa *corev1.ServiceAccount) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)
	found := &corev1.ServiceAccount{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      sa.Name,
		Namespace: instance.Spec.WorkspaceName,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		// Create the ServiceAccount
		logger.Info("Creating a new ServiceAccount", "ServiceAccount.Namespace", instance.Spec.WorkspaceName, "ServiceAccount.Name", sa.Name)
		controllerutil.SetControllerReference(instance, sa, r.Scheme)
		err = r.Create(context.TODO(), sa)

		if err != nil {
			// Creation failed
			logger.Error(err, "Failed to create new ServiceAccount", "ServiceAccount.Namespace", instance.Spec.WorkspaceName, "ServiceAccount.Name", sa.Name)
			return &ctrl.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the ServiceAccount not existing
		logger.Error(err, "Failed to get ServiceAccount")
		return &ctrl.Result{}, err
	}

	return nil, nil
}
