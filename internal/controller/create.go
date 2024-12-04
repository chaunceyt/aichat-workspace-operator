package controller

import (
	"fmt"

	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	appsv1alpha1 "github.com/chaunceyt/aichat-workspace-operator/api/v1alpha1"
)

type CreateAIChatWorkspaceStep struct {
	next AIChatWorkspace
}

func (step *CreateAIChatWorkspaceStep) execute(instance *AIChatWorkspaceInstance) (ctrl.Result, error) {

	var result *ctrl.Result
	var err error
	isCreated := instance.aichatWorkspaceConfig.Status.IsCreated
	pendingDeletion := instance.aichatWorkspaceConfig.ObjectMeta.DeletionTimestamp != nil

	if !isCreated && !pendingDeletion {
		result, err = instance.r.handleReconcile(instance.ctx, result, instance.aichatWorkspaceConfig)
		if result != nil {
			return instance.r.finishReconcile(err, false)
		}

		instance.aichatWorkspaceConfig.Status.IsCreated = true

		if err = instance.r.Status().Update(instance.ctx, instance.aichatWorkspaceConfig); err != nil {
			apimeta.SetStatusCondition(&instance.aichatWorkspaceConfig.Status.Conditions, metav1.Condition{
				Status:             metav1.ConditionFalse,
				Reason:             appsv1alpha1.ReconciliationFailedReason,
				Message:            err.Error(),
				Type:               appsv1alpha1.ConditionTypeReady,
				ObservedGeneration: instance.aichatWorkspaceConfig.GetGeneration(),
			})
			if err = instance.r.patchStatus(instance.ctx, instance.aichatWorkspaceConfig); err != nil {
				err = fmt.Errorf("unable to patch status after progressing: %w", err)
				return instance.r.finishReconcile(err, false)
			}
		}
		apimeta.SetStatusCondition(&instance.aichatWorkspaceConfig.Status.Conditions, metav1.Condition{
			Status:             metav1.ConditionFalse,
			Reason:             appsv1alpha1.ProgressingReason,
			Message:            "Reconciliation progressing",
			Type:               appsv1alpha1.ConditionTypeReady,
			ObservedGeneration: instance.aichatWorkspaceConfig.GetGeneration(),
		})
		apimeta.SetStatusCondition(&instance.aichatWorkspaceConfig.Status.Conditions, metav1.Condition{
			Status:             metav1.ConditionTrue,
			Reason:             appsv1alpha1.ReconciliationSucceededReason,
			Message:            "AIChatWorkspace reconciled",
			Type:               appsv1alpha1.ConditionTypeReady,
			ObservedGeneration: instance.aichatWorkspaceConfig.GetGeneration(),
		})

		if err = instance.r.patchStatus(instance.ctx, instance.aichatWorkspaceConfig); err != nil {
			err = fmt.Errorf("unable to patch status after progressing: %w", err)
			return instance.r.finishReconcile(err, false)
		}
		instance.r.Recorder.Event(instance.aichatWorkspaceConfig, "Normal", "Created",
			fmt.Sprintf("aichatWorkspace %s was created in namespace %s",
				instance.aichatWorkspaceConfig.Name,
				instance.aichatWorkspaceConfig.Namespace))

	}

	if step.next == nil {
		return instance.r.finishReconcile(nil, false)
	}

	return step.next.execute(instance)
}

func (step *CreateAIChatWorkspaceStep) setNext(next AIChatWorkspace) {
	step.next = next
}
