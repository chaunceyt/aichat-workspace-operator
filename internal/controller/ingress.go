package controller

import (
	"context"

	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	appsv1alpha1 "github.com/chaunceyt/aichat-workspace-operator/api/v1alpha1"
)

func (r *AIChatWorkspaceReconciler) ensureIngress(ctx context.Context, instance *appsv1alpha1.AIChatWorkspace, ing *networkingv1beta1.Ingress) (*reconcile.Result, error) {
	logger := log.FromContext(ctx)
	found := &networkingv1beta1.Ingress{}

	err := r.Get(context.TODO(), types.NamespacedName{
		Name:      ing.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the ingress
		logger.Info("Creating a new Ingress", "Ingress.Namespace", ing.Namespace, "Ingress.Name", ing.Name)
		err = r.Create(context.TODO(), ing)

		if err != nil {
			// Creation failed
			logger.Error(err, "Failed to create new Ingress", "Ingress.Namespace", ing.Namespace, "Ingress.Name", ing.Name)
			return &reconcile.Result{}, err
		}
		// Creation was successful
		return nil, nil

	} else if err != nil {
		// Error that isn't due to the ingress not existing
		logger.Error(err, "Failed to get Ingress")
		return &reconcile.Result{}, err
	}

	return nil, nil
}
