package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	kedahttpv1alpha1 "github.com/kedacore/http-add-on/operator/apis/http/v1alpha1"

	appsv1alpha1 "github.com/chaunceyt/aichat-workspace-operator/api/v1alpha1"
)

func (r *AIChatWorkspaceReconciler) ensureHTTPScaledObject(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, httpso *kedahttpv1alpha1.HTTPScaledObject) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)
	scheme := runtime.NewScheme()
	_ = kedahttpv1alpha1.AddToScheme(scheme)
	found := &kedahttpv1alpha1.HTTPScaledObject{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      httpso.Name,
		Namespace: instance.Spec.WorkspaceName,
	}, found)

	if err != nil && errors.IsNotFound(err) {
		logger.Info("Creating a HTTPScaledObject", "HTTPScaledObject.Namespace", instance.Spec.WorkspaceName, "HTTPScaledObject.Name", httpso.Name)
		controllerutil.SetControllerReference(instance, httpso, r.Scheme)
		err = r.Create(context.TODO(), httpso)

		if err != nil {
			logger.Error(err, "Failed to createHTTPScaledObject", "HTTPScaledObject.Namespace", instance.Spec.WorkspaceName, "HTTPScaledObject.Name", httpso.Name)
			return &ctrl.Result{}, err
		}

		return nil, nil

	} else if err != nil {
		logger.Error(err, "Failed to get HTTPScaledObject")

		return &ctrl.Result{}, err
	}

	return nil, nil
}
