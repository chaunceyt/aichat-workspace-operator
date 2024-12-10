/*
Copyright 2024 AIChatWorkspace Contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
