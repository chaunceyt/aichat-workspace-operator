package controller

import (
	ctrl "sigs.k8s.io/controller-runtime"
)

type UpdateAIChatWorkspaceStep struct {
	next AIChatWorkspace
}

func (step *UpdateAIChatWorkspaceStep) execute(instance *AIChatWorkspaceInstance) (ctrl.Result, error) {
	var err error
	var result *ctrl.Result
	isCreated := instance.aichatWorkspaceConfig.Status.IsCreated
	pendingDeletion := instance.aichatWorkspaceConfig.ObjectMeta.DeletionTimestamp != nil

	if isCreated && !pendingDeletion {
		instance.logger.Info("reconciling aichat", "aichat", instance.aichatWorkspaceConfig, "action", "update")

		result, err = instance.r.handleReconcile(instance.ctx, result, instance.aichatWorkspaceConfig)
		if result != nil {
			return instance.r.finishReconcile(err, true)
		}

		return step.next.execute(instance)
	}
	return step.next.execute(instance)

}

func (step *UpdateAIChatWorkspaceStep) setNext(next AIChatWorkspace) {
	step.next = next
}
