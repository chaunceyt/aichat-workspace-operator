package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/chaunceyt/aichat-workspace-operator/api/v1alpha1"
)

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

	return nil, nil
}
