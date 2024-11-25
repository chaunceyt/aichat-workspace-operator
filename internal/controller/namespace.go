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

func (r *AIChatWorkspaceReconciler) ensureNamespace(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, ns *corev1.Namespace) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)
	found := &corev1.Namespace{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name: ns.Name,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		// Create the namespace
		logger.Info("Creating a new Namespace", "instance.Spec.Namespace", instance.Spec.WorkspaceName)
		controllerutil.SetControllerReference(instance, ns, r.Scheme)
		err = r.Create(context.TODO(), ns)

		if err != nil {
			// Creation failed
			logger.Error(err, "Failed to create new Namespace", "instance.Spec.Namespace", instance.Spec.WorkspaceName)
			return &ctrl.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the namespace not existing
		logger.Error(err, "Failed to get namespace")
		return &ctrl.Result{}, err
	}

	return nil, nil
}
