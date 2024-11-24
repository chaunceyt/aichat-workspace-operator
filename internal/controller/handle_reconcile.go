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
	  AIChat Workspace
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

	// ensureNamespace - create the "aichatworkspace" namespace that contains all the components required
	// to run the AIChat Workspace.
	namespaceDefaultLabels := defaultLabels(aichat.Spec.WorkspaceName, aichat.Spec.WorkspaceName, "aichatworkspace")
	result, err = r.ensureNamespace(ctx, aichat, k8s.NewNamespace(aichat.Spec.WorkspaceName, namespaceDefaultLabels))
	if result != nil {
		return result, err
	}

	// ensurePVC - ensure the persistentvolumeclaim for web files folder is managed.
	pvcName := generateName(aichat.Spec.WorkspaceName, "openwebui")
	openwebuiPVCLabels := defaultLabels(aichat.Spec.WorkspaceName, pvcName, "sa")
	result, err = r.ensurePVC(ctx, aichat, k8s.NewPersistentVolumeClaim(pvcName, aichat.Spec.WorkspaceName, "2Gi", openwebuiPVCLabels))
	if result != nil {
		return result, err
	}

	// serviceAccount for the Open WebUI workload.
	serviceAccountForOpenWebUIName := generateName(aichat.Spec.WorkspaceName, "openwebui")
	openwebuiDefaultLabels := defaultLabels(aichat.Spec.WorkspaceName, serviceAccountForOpenWebUIName, "sa")
	result, err = r.ensureServiceAccount(ctx, aichat, k8s.NewServiceAccount(serviceAccountForOpenWebUIName, aichat.Spec.WorkspaceName, openwebuiDefaultLabels))
	if result != nil {
		return result, err
	}

	// serviceAccout for the Ollama workload
	serviceAccountForOllamaName := generateName(aichat.Spec.WorkspaceName, "ollama")
	ollamaDefaultLabels := defaultLabels(aichat.Spec.WorkspaceName, serviceAccountForOllamaName, "sa")
	result, err = r.ensureServiceAccount(ctx, aichat, k8s.NewServiceAccount(serviceAccountForOllamaName, aichat.Spec.WorkspaceName, ollamaDefaultLabels))
	if result != nil {
		return result, err
	}

	// ensureStatefulSet - creating the StatefulSet used to run the Ollama API
	ollamaName := generateName(aichat.Spec.WorkspaceName, "ollama")
	ollamaContainerPort := int32(11434)
	ollamaVolumeSize := int32(20)
	result, err = r.ensureStatefulSet(ctx, aichat, k8s.NewStatefulSet(aichat.Spec.WorkspaceName, ollamaName, ollamaContainerPort, ollamaVolumeSize))
	if result != nil {
		return result, err
	}

	// ensureService - creating the Service used to route traffic to the Ollama API pod.
	ollamaServiceDefaultLabels := defaultLabels(aichat.Spec.WorkspaceName, ollamaName, "svc")
	result, err = r.ensureService(ctx, aichat, k8s.NewService(aichat.Spec.WorkspaceName, ollamaName, ollamaContainerPort, ollamaServiceDefaultLabels))
	if result != nil {
		return result, err
	}

	// ensureDeployment - creating the Deployment used to the Open WebUI workload.
	openwebuiName := generateName(aichat.Spec.WorkspaceName, "openwebui")
	openwebuiContainerPort := int32(8080)
	result, err = r.ensureDeployment(ctx, aichat, k8s.NewDeployment(aichat.Spec.WorkspaceName, openwebuiName, openwebuiContainerPort))
	if result != nil {
		return result, err
	}

	// ensureService - creating the Service used to route traffic to the Open WebUI pod.
	openwebuiServiceDefaultLabels := defaultLabels(aichat.Spec.WorkspaceName, ollamaName, "svc")
	result, err = r.ensureService(ctx, aichat, k8s.NewService(aichat.Spec.WorkspaceName, openwebuiName, openwebuiContainerPort, openwebuiServiceDefaultLabels))
	if result != nil {
		return result, err
	}

	return result, nil
}
