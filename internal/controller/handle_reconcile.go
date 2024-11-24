package controller

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/chaunceyt/aichat-workspace-operator/api/v1alpha1"
	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/k8s"
)

func (r *AIChatWorkspaceReconciler) handleReconcile(ctx context.Context, result *ctrl.Result, aichat *appsv1alpha1.AIChatWorkspace) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var err error
	logger.Info("reconciling aichatworkspace")

	/**
	  Create the following:
	  some derived from https://github.com/open-webui/open-webui/tree/main/kubernetes/manifest/base
	  - namespace
	  - resourcequota
	  - serviceaccount for ollama api
	  - serviceaccount for openwebui
	  - statefulset for Ollama API
	  - service for Ollama API
	  - deployment for Open WebUI
	  - service for Open WebUI
	  - ingress for Open WebUI
	**/

	result, err = r.ensureNamespace(ctx, aichat, k8s.NewNamespace(aichat.Spec.WorkspaceName))
	if result != nil {
		return result, err
	}

	//ensurePVC - ensure the persistentvolumeclaim for web files folder is managed.
	pvcName := generateName(aichat.Spec.WorkspaceName, "openwebui")
	result, err = r.ensurePVC(ctx, aichat, k8s.NewPersistentVolumeClaim(pvcName, aichat.Spec.WorkspaceName, "2Gi"))
	if result != nil {
		return result, err
	}

	serviceAccountForOpenWebUIName := generateName(aichat.Spec.WorkspaceName, "openwebui")
	result, err = r.ensureServiceAccount(ctx, aichat, k8s.NewServiceAccount(serviceAccountForOpenWebUIName, aichat.Spec.WorkspaceName))
	if result != nil {
		return result, err
	}

	serviceAccountForOllamaName := generateName(aichat.Spec.WorkspaceName, "ollama")
	result, err = r.ensureServiceAccount(ctx, aichat, k8s.NewServiceAccount(serviceAccountForOllamaName, aichat.Spec.WorkspaceName))
	if result != nil {
		return result, err
	}

	ollamaName := generateName(aichat.Spec.WorkspaceName, "ollama")
	result, err = r.ensureStatefulSet(ctx, aichat, k8s.NewStatefulSet(aichat.Spec.WorkspaceName, ollamaName, 11434, 20))
	if result != nil {
		return result, err
	}

	result, err = r.ensureService(ctx, aichat, k8s.NewService(aichat.Spec.WorkspaceName, ollamaName, int32(11434)))
	if result != nil {
		return result, err
	}

	openwebuiName := generateName(aichat.Spec.WorkspaceName, "openwebui")
	result, err = r.ensureDeployment(ctx, aichat, k8s.NewDeployment(aichat.Spec.WorkspaceName, openwebuiName, 8080))
	if result != nil {
		return result, err
	}

	result, err = r.ensureService(ctx, aichat, k8s.NewService(aichat.Spec.WorkspaceName, openwebuiName, int32(8080)))
	if result != nil {
		return result, err
	}

	return result, nil
}
