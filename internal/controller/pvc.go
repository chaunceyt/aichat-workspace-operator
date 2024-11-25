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

func (r *AIChatWorkspaceReconciler) ensurePVC(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, pvc *corev1.PersistentVolumeClaim) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)
	found := &corev1.PersistentVolumeClaim{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      pvc.Name,
		Namespace: instance.Spec.WorkspaceName,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		// Create the PVC
		logger.Info("Creating a new PVC", "PVC.Namespace", instance.Spec.WorkspaceName, "PVC.Name", pvc.Name)
		controllerutil.SetControllerReference(instance, pvc, r.Scheme)
		err = r.Create(context.TODO(), pvc)

		if err != nil {
			// Creation failed
			logger.Error(err, "Failed to create new PVC", "PVC.Namespace", instance.Spec.WorkspaceName, "PVC.Name", pvc.Name)
			return &ctrl.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the pvc not existing
		logger.Error(err, "Failed to get PVC")
		return &ctrl.Result{}, err
	}

	return nil, nil
}
