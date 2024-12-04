package controller

import (
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type DeleteAIChatWorkspaceStep struct {
	next AIChatWorkspace
}

func (step *DeleteAIChatWorkspaceStep) execute(instance *AIChatWorkspaceInstance) (ctrl.Result, error) {
	var err error
	isCreated := instance.aichatWorkspaceConfig.Status.IsCreated
	pendingDeletion := instance.aichatWorkspaceConfig.ObjectMeta.DeletionTimestamp != nil

	if isCreated && pendingDeletion {
		instance.logger.Info("reconciling aichat", "aichat", instance.aichatWorkspaceConfig, "action", "delete")
		if controllerutil.ContainsFinalizer(instance.aichatWorkspaceConfig, aichatWorkspaceFinalizerName) {
			if err = instance.r.deleteAIChatWorkspace(instance.ctx, instance.aichatWorkspaceConfig); err != nil {
				return instance.r.finishReconcile(err, false)
			}
			instance.r.Recorder.Event(instance.aichatWorkspaceConfig, "Warning", "Deleting",
				fmt.Sprintf("aichatWorkspace %s is being deleted from the namespace %s",
					instance.aichatWorkspaceConfig.Name,
					instance.aichatWorkspaceConfig.Namespace))
			controllerutil.RemoveFinalizer(instance.aichatWorkspaceConfig, aichatWorkspaceFinalizerName)
			if err = instance.r.Update(instance.ctx, instance.aichatWorkspaceConfig); err != nil {
				return instance.r.finishReconcile(err, false)
			}
		}

	}

	if step.next == nil {
		return instance.r.finishReconcile(nil, false)
	}

	return step.next.execute(instance)
}

func (step *DeleteAIChatWorkspaceStep) setNext(next AIChatWorkspace) {
	step.next = next
}
