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

	if step.next == nil {
		return instance.r.finishReconcile(nil, false)
	}

	return step.next.execute(instance)

}

func (step *InitAIChatWorkspaceStep) setNext(next AIChatWorkspace) {
	step.next = next
}
