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

func (r *AIChatWorkspaceReconciler) ensureDeployment(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, deploy *appsv1.Deployment) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)
	found := &appsv1.Deployment{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      deploy.Name,
		Namespace: instance.Spec.WorkspaceName,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		// Create the Deployment
		logger.Info("Creating a new Deployment", "Deployment.Namespace", instance.Spec.WorkspaceName, "Deployment.Name", deploy.Name)
		controllerutil.SetControllerReference(instance, deploy, r.Scheme)
		err = r.Create(context.TODO(), deploy)

		if err != nil {
			// Creation failed
			logger.Error(err, "Failed to create new Deployment", "Deployment.Namespace", instance.Spec.WorkspaceName, "Deployment.Name", deploy.Name)
			return &ctrl.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the Deployment not existing
		logger.Error(err, "Failed to get Deployment")
		return &ctrl.Result{}, err
	}

	return nil, nil
}
