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

package k8s

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/chaunceyt/aichat-workspace-operator/internal/constants"
)

const (
	defaultNameLabel            = "app.kubernetes.io/name"
	defaultInstanceLabel        = "app.kubernetes.io/instance"
	openwebuiVolumeMountName    = "webui-volume"
	openwebuiContainerName      = "open-webui"
	openwebuiContainerImageName = "ghcr.io/open-webui/open-webui"
	ollamaVolumeMountName       = "ollama-volume"
	ollamaContainerName         = "ollama"
	ollamaContainerImageName    = "ollama/ollama"
)

// NewDeployment is responsible for creating the Open WebUI workload.
func NewDeployment(namespace, name string, port int32, openwebuiContainerImageTag string) *appsv1.Deployment {
	appLabels := map[string]string{defaultNameLabel: name}

	// config for ollama service
	containerImage := fmt.Sprintf("%s:%s", constants.OpenwebuiContainerImageName, openwebuiContainerImageTag)
	serviceName := fmt.Sprintf("%s-ollama", namespace)
	ollamaServerURI := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", serviceName, namespace, constants.OllamaPort)
	openAIURI := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d/v1", serviceName, namespace, constants.OllamaPort)
	workspaceName := fmt.Sprintf("AIChat Workspace: %s", namespace)
	saName := fmt.Sprintf("%s-openwebui", namespace)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{MatchLabels: appLabels},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: appLabels},
				Spec: v1.PodSpec{
					RestartPolicy:                v1.RestartPolicyAlways,
					ServiceAccountName:           saName,
					AutomountServiceAccountToken: ptr.To[bool](false),
					// at the moment having issues getting open webui to run as non-root
					// SecurityContext:              defaultPodSecurityContext(),
					Containers: []v1.Container{
						{
							Name:  constants.OpenwebuiContainerName,
							Image: containerImage,
							// SecurityContext: defaultSecurityContext(),
							Env: []v1.EnvVar{
								{
									Name:  "OLLAMA_BASE_URL",
									Value: ollamaServerURI,
								},
								{
									Name:  "OPENAI_API_BASE_URL",
									Value: openAIURI,
								},
								{
									Name:  "ENV",
									Value: "dev",
								},
								{
									Name:  "WEBUI_NAME",
									Value: workspaceName,
								},
								{
									Name:  "KEY_FILE",
									Value: "/tmp/.webui_secret_key",
								},
							},
							Ports: []v1.ContainerPort{{ContainerPort: port}},
							TTY:   true,
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      constants.OpenwebuiVolumeMountName,
									MountPath: constants.OpenwebuiVolumeMountPath,
								},
							},
						},
					},
					Volumes: []v1.Volume{
						{
							Name: openwebuiVolumeMountName,
							VolumeSource: v1.VolumeSource{
								PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
									ClaimName: name,
									ReadOnly:  false,
								},
							},
						},
					},
				},
			},
		},
	}
}

// NewStatefulSet is responsible for creating the Ollama workload.
func NewStatefulSet(namespace, name string, port int32, volumeSize string, ollamaContainerImageTag string) *appsv1.StatefulSet {
	appLabels := map[string]string{defaultNameLabel: name}

	// config for Open WebUI
	containerImage := fmt.Sprintf("%s:%s", constants.OllamaContainerImageName, ollamaContainerImageTag)
	saName := fmt.Sprintf("%s-ollama", namespace)
	serviceName := fmt.Sprintf("%s-%s", namespace, constants.OllamaName)

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector:    &metav1.LabelSelector{MatchLabels: appLabels},
			ServiceName: serviceName,
			VolumeClaimTemplates: []v1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      ollamaVolumeMountName,
						Namespace: namespace,
					},
					Spec: v1.PersistentVolumeClaimSpec{
						AccessModes: []v1.PersistentVolumeAccessMode{
							v1.ReadWriteOnce,
						},
						Resources: v1.VolumeResourceRequirements{
							Requests: v1.ResourceList{"storage": resource.MustParse(volumeSize)},
						},
					},
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: appLabels},
				Spec: v1.PodSpec{
					RestartPolicy:                v1.RestartPolicyAlways,
					ServiceAccountName:           saName,
					AutomountServiceAccountToken: ptr.To[bool](false),
					SecurityContext:              defaultPodSecurityContext(),
					Containers: []v1.Container{
						{
							Name:  constants.OllamaContainerName,
							Image: containerImage,
							Env: []v1.EnvVar{
								{
									Name:  "OLLAMA_DEBUG",
									Value: "1",
								},
							},
							SecurityContext: defaultSecurityContext(),
							Ports:           []v1.ContainerPort{{ContainerPort: port}},
							TTY:             true,
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      ollamaVolumeMountName,
									MountPath: "/.ollama",
								},
							},
						},
					},
				},
			},
		},
	}
}

// defaultSecurityContext - sets the security context for each container
func defaultSecurityContext() *v1.SecurityContext {
	return &v1.SecurityContext{
		AllowPrivilegeEscalation: ptr.To[bool](false),
		Capabilities: &v1.Capabilities{
			Drop: []v1.Capability{
				"ALL",
			},
		},
		Privileged:             ptr.To[bool](false),
		ReadOnlyRootFilesystem: ptr.To[bool](true),
		RunAsNonRoot:           ptr.To[bool](true),
		SeccompProfile: &v1.SeccompProfile{
			Type: v1.SeccompProfileType("RuntimeDefault"),
		},
	}
}

// defaultPodSecurityContext - sets the pod's security context
func defaultPodSecurityContext() *v1.PodSecurityContext {
	return &v1.PodSecurityContext{
		FSGroup:    ptr.To[int64](10001),
		RunAsUser:  ptr.To[int64](10001),
		RunAsGroup: ptr.To[int64](10001),
	}
}
