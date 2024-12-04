package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1alpha1 "github.com/chaunceyt/aichat-workspace-operator/api/v1alpha1"
)

type InitAIChatWorkspaceStep struct {
	next AIChatWorkspace
}

func (step *InitAIChatWorkspaceStep) execute(instance *AIChatWorkspaceInstance) (ctrl.Result, error) {
	instance.logger.Info("starting InitStep")
	aichatWorkspaceConfig := &appsv1alpha1.AIChatWorkspace{}
	err := instance.r.Get(instance.ctx, instance.req.NamespacedName, aichatWorkspaceConfig)
	if err != nil {
		// Error reading the object - requeue the request.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	instance.aichatWorkspaceConfig = aichatWorkspaceConfig

	instance.logger.Info("ending InitStep")

	return step.next.execute(instance)

}

func (step *InitAIChatWorkspaceStep) setNext(next AIChatWorkspace) {
	step.next = next
}
