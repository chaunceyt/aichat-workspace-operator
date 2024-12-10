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

package constants

const (
	Version                      = "0.0.1"
	ManagedBy                    = "aichat-workspace-operator"
	AIChatWorkspaceName          = "aichatworkspace"
	AIChatWorkspaceFinalizerName = "core.aichatworkspace.io/finalizer"
	AIChatWorkspaceNamespace     = "aichat-workspace-operator-system"
	AIChatWorspaceConfigMapName  = "aichat-workspace-operator-config"

	// Open WebUI
	OpenwebuiName               = "openwebui"
	OpenwebuiVolumeMountName    = "webui-volume"
	OpenwebuiVolumeMountPath    = "/app/backend/data"
	OpenwebuiContainerName      = "open-webui"
	OpenwebuiContainerImageName = "ghcr.io/open-webui/open-webui"

	OpenwebuiContainerPort     = int32(8080)
	OpenwebuiDefaultVolumeSize = "2Gi"

	// Ollama
	OllamaName               = "ollama"
	OllamaVolumeMountName    = "ollama-volume"
	OllamaContainerName      = "ollama"
	OllamaContainerImageName = "ollama/ollama"
	OllamaPort               = int32(11434)
	OllamaDefaultVolumeSize  = "20Gi"

	// KEDA scaled-to-zero
	KedaHttpInterceptorProxy = "keda-add-ons-http-interceptor-proxy.keda"

	// ResourceQuota
	ResourceQuotaName         = "rquota"
	MaxPods                   = "2"
	MaxPersistentVolumeClaims = "2"
	MaxService                = "5"

	// Label Names
	ServiceLabelName        = "svc"
	ServiceAccountLabelName = "sa"
	ResourceQuotaLabelName  = "resourceQuota"
	PVCLabelName            = "pvc"

	// Configmap Keys
	OpenwebUIImageTag = "openwebUIImageTag"
	OllamaImageTag    = "ollamaImageTag"
)
