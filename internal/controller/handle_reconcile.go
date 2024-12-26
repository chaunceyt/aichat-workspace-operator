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
	"context"
	"fmt"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/chaunceyt/aichat-workspace-operator/api/v1alpha1"
	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/k8s"
	"github.com/chaunceyt/aichat-workspace-operator/internal/config"
	"github.com/chaunceyt/aichat-workspace-operator/internal/constants"
)

func (r *AIChatWorkspaceReconciler) handleReconcile(ctx context.Context, result *ctrl.Result, aichat *appsv1alpha1.AIChatWorkspace) (*ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var err error

	// get config
	config, err := config.GetConfig()
	if err != nil {
		return result, err
	}

	logger.Info("reconciling aichatworkspace")

	// ensureNamespace - create the "aichatworkspace" namespace that contains all the components required
	// to run the AIChat Workspace.
	namespaceDefaultLabels := defaultLabels(aichat.Spec.WorkspaceName, aichat.Spec.WorkspaceName, constants.AIChatWorkspaceName)
	result, err = r.ensureNamespace(ctx, aichat, k8s.NewNamespace(aichat.Spec.WorkspaceName, namespaceDefaultLabels))
	if result != nil {
		logger.Error(
			err,
			"Creating AIChat Workspace (namespace)",
			"Name", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	// ensureResourceQuota - creates the ResourceQuota object that limits resources that can be ran in the namespace.
	// resourceQuotaName := generateName(aichat.Spec.WorkspaceName, constants.ResourceQuotaName)
	// resourceQuotaDefaultLabels := defaultLabels(aichat.Spec.WorkspaceName, aichat.Spec.WorkspaceName, constants.ResourceQuotaLabelName)
	// result, err = r.ensureResourceQuota(ctx, aichat, k8s.NewResourceQuota(aichat.Spec.WorkspaceName, resourceQuotaName, resourceQuotaDefaultLabels))
	// if result != nil {
	// 	logger.Error(
	// 		err,
	// 		"Creating ResourceQuota",
	// 		"Name", resourceQuotaName,
	// 		"Workspace", aichat.Spec.WorkspaceName,
	// 	)
	// 	return result, err
	// }

	result, err = r.ensureInitConfigMap(ctx, aichat, k8s.NewInitConfigMap(aichat.Spec.WorkspaceName))
	if result != nil {
		logger.Error(
			err,
			"Creating InitConfigMap",
			"Name", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	// ensurePVC - ensure the persistentvolumeclaim for Open WebUI is managed.
	pvcName := generateName(aichat.Spec.WorkspaceName, constants.OpenwebuiName)
	openwebuiPVCLabels := defaultLabels(aichat.Spec.WorkspaceName, pvcName, constants.PVCLabelName)
	result, err = r.ensurePVC(ctx, aichat, k8s.NewPersistentVolumeClaim(pvcName, aichat.Spec.WorkspaceName, constants.OpenwebuiDefaultVolumeSize, openwebuiPVCLabels))
	if result != nil {
		logger.Error(
			err,
			"Creating PVC",
			"Name", pvcName,
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	// serviceAccount for the Open WebUI workload.
	serviceAccountForOpenWebUIName := generateName(aichat.Spec.WorkspaceName, constants.OpenwebuiName)
	openwebuiDefaultLabels := defaultLabels(aichat.Spec.WorkspaceName, serviceAccountForOpenWebUIName, constants.ServiceAccountLabelName)
	result, err = r.ensureServiceAccount(ctx, aichat, k8s.NewServiceAccount(serviceAccountForOpenWebUIName, aichat.Spec.WorkspaceName, openwebuiDefaultLabels))
	if result != nil {
		logger.Error(
			err,
			"Creating ServiceAccount",
			"Name", serviceAccountForOpenWebUIName,
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	// serviceAccout for the Ollama workload
	serviceAccountForOllamaName := generateName(aichat.Spec.WorkspaceName, constants.OllamaName)
	ollamaDefaultLabels := defaultLabels(aichat.Spec.WorkspaceName, serviceAccountForOllamaName, constants.ServiceAccountLabelName)
	result, err = r.ensureServiceAccount(ctx, aichat, k8s.NewServiceAccount(serviceAccountForOllamaName, aichat.Spec.WorkspaceName, ollamaDefaultLabels))
	if result != nil {
		logger.Error(
			err,
			"Creating ServiceAccount",
			"Name", serviceAccountForOllamaName,
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	// serviceAccout for the Postgres workload
	serviceAccountForPostgresName := generateName(aichat.Spec.WorkspaceName, "postgres")
	postgresDefaultLabels := defaultLabels(aichat.Spec.WorkspaceName, serviceAccountForPostgresName, constants.ServiceAccountLabelName)
	result, err = r.ensureServiceAccount(ctx, aichat, k8s.NewServiceAccount(serviceAccountForPostgresName, aichat.Spec.WorkspaceName, postgresDefaultLabels))
	if result != nil {
		logger.Error(
			err,
			"Creating ServiceAccount",
			"Name", serviceAccountForPostgresName,
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	// serviceAccout for the n8n workload
	serviceAccountForN8NName := generateName(aichat.Spec.WorkspaceName, "n8n")
	n8nDefaultLabels := defaultLabels(aichat.Spec.WorkspaceName, serviceAccountForN8NName, constants.ServiceAccountLabelName)
	result, err = r.ensureServiceAccount(ctx, aichat, k8s.NewServiceAccount(serviceAccountForN8NName, aichat.Spec.WorkspaceName, n8nDefaultLabels))
	if result != nil {
		logger.Error(
			err,
			"Creating ServiceAccount",
			"Name", serviceAccountForN8NName,
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	// serviceAccout for the qdrant workload
	serviceAccountForQdrantName := generateName(aichat.Spec.WorkspaceName, "qdrant")
	qdrantDefaultLabels := defaultLabels(aichat.Spec.WorkspaceName, serviceAccountForQdrantName, constants.ServiceAccountLabelName)
	result, err = r.ensureServiceAccount(ctx, aichat, k8s.NewServiceAccount(serviceAccountForQdrantName, aichat.Spec.WorkspaceName, qdrantDefaultLabels))
	if result != nil {
		logger.Error(
			err,
			"Creating ServiceAccount",
			"Name", serviceAccountForQdrantName,
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	// ensureStatefulSet - creating the StatefulSet used to run the Ollama API
	ollamaName := generateName(aichat.Spec.WorkspaceName, constants.OllamaName)
	result, err = r.ensureOllama(ctx, aichat, k8s.NewOllama(aichat.Spec.WorkspaceName, ollamaName, constants.OllamaPort, constants.OllamaDefaultVolumeSize, config.OllamaImageTag))
	if result != nil {
		logger.Error(
			err,
			"Creating StatefulSet",
			"Name", ollamaName,
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	// ensureService - creating the Service used to route traffic to the Ollama API pod.
	ollamaServiceDefaultLabels := defaultLabels(aichat.Spec.WorkspaceName, ollamaName, constants.ServiceLabelName)
	result, err = r.ensureService(ctx, aichat, k8s.NewService(aichat.Spec.WorkspaceName, ollamaName, constants.OllamaPort, ollamaServiceDefaultLabels))
	if result != nil {
		logger.Error(
			err,
			"Creating Service",
			"Workload", "ollama",
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	// ensureStatefulSet - creating the StatefulSet used to run the Postgres
	postgresName := "postgres"
	pgPort := int32(5432)
	result, err = r.ensurePostgres(ctx, aichat, k8s.NewPostgres(aichat.Spec.WorkspaceName, postgresName, pgPort, "pg17"))
	if result != nil {
		logger.Error(
			err,
			"Creating StatefulSet",
			"Name", postgresName,
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	// ensureService - creating the Service used to route traffic to the Postgres pod.
	postgresServiceDefaultLabels := defaultLabels(aichat.Spec.WorkspaceName, postgresName, constants.ServiceLabelName)
	result, err = r.ensureService(ctx, aichat, k8s.NewService(aichat.Spec.WorkspaceName, postgresName, pgPort, postgresServiceDefaultLabels))
	if result != nil {
		logger.Error(
			err,
			"Creating Service",
			"Workload", "ollama",
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	// ensureDeployment - creating the Deployment used to deploy the Open WebUI workload.
	openwebuiPipelinesName := generateName(aichat.Spec.WorkspaceName, "openwebui-pipelines")
	pipelinesPort := int32(9099)
	result, err = r.ensureOpenWebUIPipelines(ctx, aichat, k8s.NewOpenWebUIPipelines(aichat.Spec.WorkspaceName, openwebuiPipelinesName, pipelinesPort, "main"))
	if result != nil {
		logger.Error(
			err,
			"Creating Deployment",
			"Name", openwebuiPipelinesName,
			"Workload", "openwebui",
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	// ensureService - creating the Service used to route traffic to the Open WebUI Pipelines pod.
	openwebuiPipelinesServiceDefaultLabels := defaultLabels(aichat.Spec.WorkspaceName, openwebuiPipelinesName, constants.ServiceLabelName)
	result, err = r.ensureService(ctx, aichat, k8s.NewService(aichat.Spec.WorkspaceName, openwebuiPipelinesName, pipelinesPort, openwebuiPipelinesServiceDefaultLabels))
	if result != nil {
		logger.Error(
			err,
			"Creating Service",
			"Workload", "openwebui-pipelines",
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	// ensureDeployment - creating the Deployment used to deploy the Open WebUI workload.
	openwebuiName := generateName(aichat.Spec.WorkspaceName, constants.OpenwebuiName)
	result, err = r.ensureOpenWebUI(ctx, aichat, k8s.NewOpenWebUI(aichat.Spec.WorkspaceName, openwebuiName, constants.OpenwebuiContainerPort, config.OpenwebUIImageTag))
	if result != nil {
		logger.Error(
			err,
			"Creating Deployment",
			"Name", openwebuiName,
			"Workload", "openwebui",
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	// ensureService - creating the Service used to route traffic to the Open WebUI pod.
	openwebuiServiceDefaultLabels := defaultLabels(aichat.Spec.WorkspaceName, openwebuiName, constants.ServiceLabelName)
	result, err = r.ensureService(ctx, aichat, k8s.NewService(aichat.Spec.WorkspaceName, openwebuiName, constants.OpenwebuiContainerPort, openwebuiServiceDefaultLabels))
	if result != nil {
		logger.Error(
			err,
			"Creating Service",
			"Workload", "openwebuiPipelines",
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	n8nPort := int32(5378)
	result, err = r.ensureN8N(ctx, aichat, k8s.NewN8N(aichat.Spec.WorkspaceName, "n8n", n8nPort, "1.73.1"))
	if result != nil {
		logger.Error(
			err,
			"Creating n8n StatefulSet",
			"Workload", "n8n",
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	qDrantPort := int32(6333)
	result, err = r.ensureQdrant(ctx, aichat, k8s.NewQdrant(aichat.Spec.WorkspaceName, "qdrant", qDrantPort, "v1.12.5"))
	if result != nil {
		logger.Error(
			err,
			"Creating n8n StatefulSet",
			"Workload", "qdrant",
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}
	// ensureService - creating the Service used to route traffic to the Open WebUI pod.
	// openwebuiServiceDefaultLabels = defaultLabels(aichat.Spec.WorkspaceName, openwebuiName, constants.ServiceLabelName)
	// result, err = r.ensureService(ctx, aichat, k8s.NewService(aichat.Spec.WorkspaceName, openwebuiName, constants.OpenwebuiContainerPort, openwebuiServiceDefaultLabels))
	// if result != nil {
	// 	logger.Error(
	// 		err,
	// 		"Creating Service",
	// 		"Workload", "openwebui",
	// 		"Workspace", aichat.Spec.WorkspaceName,
	// 	)
	// 	return result, err
	// }

	// ensureService - creating the Service used to route traffic to the Open WebUI pod.
	// openwebuiExternalServiceDefaultLabels := defaultLabels(aichat.Spec.WorkspaceName, openwebuiName, constants.ServiceLabelName)
	// result, err = r.ensureService(ctx, aichat, k8s.NewExternalService(aichat.Spec.WorkspaceName, openwebuiExternalServiceDefaultLabels))
	// if result != nil {
	// 	logger.Error(
	// 		err,
	// 		"Creating ExternalService",
	// 		"Workload", "openwebui",
	// 		"Workspace", aichat.Spec.WorkspaceName,
	// 	)
	// 	return result, err
	// }

	// ensureIngress - creating the Ingress used for Open WebUI service
	// proxyName := fmt.Sprintf("%s", "openwebui")
	openwebBackend := getName(aichat.Spec.WorkspaceName, constants.OpenwebuiName)
	openwebuiDNSName := setIngressDNSHost(config, aichat.Spec.WorkspaceName, constants.OpenwebuiName)
	result, err = r.ensureIngress(ctx, aichat, k8s.NewIngress(aichat.Spec.WorkspaceName, constants.OpenwebuiName, openwebBackend, openwebuiDNSName, constants.OpenwebuiContainerPort))
	if result != nil {
		logger.Error(
			err,
			"Creating Ingress",
			"Workload", "openwebui",
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	// ensureIngress - creating the Ingress used for Ollama service
	ollamaBackend := getName(aichat.Spec.WorkspaceName, constants.OllamaName)
	ollamaDNSName := setIngressDNSHost(config, aichat.Spec.WorkspaceName, constants.OllamaName)
	result, err = r.ensureIngress(ctx, aichat, k8s.NewIngress(aichat.Spec.WorkspaceName, constants.OllamaName, ollamaBackend, ollamaDNSName, constants.OllamaPort))
	if result != nil {
		logger.Error(
			err,
			"Creating ingress",
			"Workload", "ollama",
			"Workspace", aichat.Spec.WorkspaceName,
		)
		return result, err
	}

	// hosts := []string{openwebuiDNSName}
	// result, err = r.ensureHTTPScaledObject(ctx, aichat, k8s.NewHttpSo(aichat.Spec.WorkspaceName, "Deployment", constants.OpenwebuiName, constants.OpenwebuiContainerPort, hosts))
	// if result != nil {
	// 	return result, err
	// }

	return result, nil
}

func getName(workspace, workload string) string {
	name := fmt.Sprintf("%s-%s", workspace, workload)
	return name
}

// setIngressDNSHost:
// - <workspaceName>.<defaultDomain.tld> for openwebui
// - <workspaceName>-api.<defaultDomain.tld> for ollama
func setIngressDNSHost(config *config.Config, workspace string, workload string) string {
	var dnsName string
	switch workload {
	case "ollama":
		dnsName = fmt.Sprintf("%s-api.%s", workspace, config.DefaultDomain)
	case "openwebui":
		dnsName = fmt.Sprintf("%s.%s", workspace, config.DefaultDomain)
	}

	return dnsName
}
