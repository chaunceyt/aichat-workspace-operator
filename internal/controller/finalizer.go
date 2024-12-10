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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type FinalizerStep struct {
	next AIChatWorkspace
}

func (step *FinalizerStep) execute(instance *AIChatWorkspaceInstance) (ctrl.Result, error) {
	var err error
	pendingDeletion := instance.aichatWorkspaceConfig.ObjectMeta.DeletionTimestamp != nil
	hasFinalizer := controllerutil.ContainsFinalizer(instance.aichatWorkspaceConfig, aichatWorkspaceFinalizerName)

	if !hasFinalizer && !pendingDeletion {
		instance.logger.Info("reconciling aichat", "aichat", instance.aichatWorkspaceConfig, "action", "add finalizer")
		controllerutil.AddFinalizer(instance.aichatWorkspaceConfig, aichatWorkspaceFinalizerName)
		if err = instance.r.Update(instance.ctx, instance.aichatWorkspaceConfig); err != nil {
			return instance.r.finishReconcile(err, false)
		}
	}

	if step.next == nil {
		return instance.r.finishReconcile(nil, false)
	}

	return step.next.execute(instance)
}

func (step *FinalizerStep) setNext(next AIChatWorkspace) {
	step.next = next
}
