package controller

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/chaunceyt/aichat-workspace-operator/api/v1alpha1"
	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/k8s"
)

func (r *AIChatWorkspaceReconciler) handlerReconcile(ctx context.Context, aichat *appsv1alpha1.AIChatWorkspace) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var err error
	logger.Info("reconciling aichatworkspace")

	var result *ctrl.Result
	// ensureNamespace - ensure the namepace for the aichatworkspace is created.
	// should match the workspacename
	result, err = r.ensureNamespace(ctx, aichat, k8s.NewNamespace(aichat.Spec.WorkspaceName))
	if result != nil {
		return *result, err
	}

	//ensurePVC - ensure the persistentvolumeclaim for web files folder is managed.
	pvcName := generateName(aichat.Spec.WorkspaceName, "openwebui")
	result, err = r.ensurePVC(ctx, aichat, k8s.NewPersistentVolumeClaim(pvcName, aichat.Spec.WorkspaceName, "2Gi"))
	if result != nil {
		return *result, err
	}

	serviceAccountForOpenWebUIName := generateName(aichat.Spec.WorkspaceName, "openwebui")
	result, err = r.ensureServiceAccount(ctx, aichat, k8s.NewServiceAccount(serviceAccountForOpenWebUIName, aichat.Spec.WorkspaceName))
	if result != nil {
		return *result, err
	}

	serviceAccountForOllamaName := generateName(aichat.Spec.WorkspaceName, "ollama")
	result, err = r.ensureServiceAccount(ctx, aichat, k8s.NewServiceAccount(serviceAccountForOllamaName, aichat.Spec.WorkspaceName))
	if result != nil {
		return *result, err
	}

	return *result, nil
}
