package k8s

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/utils"
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

var (
	openwebuiContainerImageTag = "main"
	ollamaContainerImageTag    = "0.4.1"
)

// NewDeployment is responsible for creating the Open WebUI workload.
func NewDeployment(namespace, name string, port int32) *appsv1.Deployment {
	appLabels := map[string]string{defaultNameLabel: name}

	// config for ollama service
	ollamaPort := int32(11434)
	serviceName := fmt.Sprintf("%s-ollama", namespace)
	ollamaServerURI := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", serviceName, namespace, ollamaPort)
	openAIURI := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d/v1", serviceName, namespace, ollamaPort)
	workspaceName := fmt.Sprintf("AIChat Workspace: %s", namespace)

	saName := fmt.Sprintf("%s-openwebui", namespace)

	// containerImage
	containerImage := fmt.Sprintf("%s:%s", openwebuiContainerImageName, openwebuiContainerImageTag)

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
					AutomountServiceAccountToken: utils.PtrBool(false),
					// at the moment having issues getting open webui to run as non-root
					// SecurityContext:              defaultPodSecurityContext(),
					Containers: []v1.Container{
						{
							Name:  openwebuiContainerName,
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
									Name:      openwebuiVolumeMountName,
									MountPath: "/app/backend/data",
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

func NewStatefulSet(namespace, name string, port int32, volumeSize int32) *appsv1.StatefulSet {
	appLabels := map[string]string{defaultNameLabel: name}
	// containerImage
	containerImage := fmt.Sprintf("%s:%s", ollamaContainerImageName, ollamaContainerImageTag)
	saName := fmt.Sprintf("%s-ollama", namespace)

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: appLabels},
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
							Requests: v1.ResourceList{"storage": resource.MustParse("10Gi")},
						},
					},
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: appLabels},
				Spec: v1.PodSpec{
					RestartPolicy:                v1.RestartPolicyAlways,
					ServiceAccountName:           saName,
					AutomountServiceAccountToken: utils.PtrBool(false),
					SecurityContext:              defaultPodSecurityContext(),
					Containers: []v1.Container{
						{
							Name:  ollamaContainerName,
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
		AllowPrivilegeEscalation: utils.PtrBool(false),
		Capabilities: &v1.Capabilities{
			Drop: []v1.Capability{
				"ALL",
			},
		},
		Privileged:             utils.PtrBool(false),
		ReadOnlyRootFilesystem: utils.PtrBool(true),
		RunAsNonRoot:           utils.PtrBool(true),
		SeccompProfile: &v1.SeccompProfile{
			Type: v1.SeccompProfileType("RuntimeDefault"),
		},
	}
}

// defaultPodSecurityContext - sets the pod's security context
func defaultPodSecurityContext() *v1.PodSecurityContext {
	return &v1.PodSecurityContext{
		FSGroup:    utils.PtrInt64(10001),
		RunAsUser:  utils.PtrInt64(10001),
		RunAsGroup: utils.PtrInt64(10001),
	}
}
