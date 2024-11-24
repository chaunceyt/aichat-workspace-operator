package k8s

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultNameLabel         = "app.kubernetes.io/name"
	defaultInstanceLabel     = "app.kubernetes.io/instance"
	openwebuiVolumeMountName = "webui-volume"
	ollamaVolumeMountName    = "ollama-volume"
)

func NewDeployment(namespace, name string, port int32) *appsv1.Deployment {
	appLabels := map[string]string{defaultNameLabel: name}
	ollamaPort := int32(11434)
	serviceName := fmt.Sprintf("%s-ollam", namespace)
	ollamaServerURI := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", serviceName, namespace, ollamaPort)

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
					RestartPolicy: v1.RestartPolicyAlways,
					Containers: []v1.Container{
						{
							Name:  "open-webui",
							Image: "ghcr.io/open-webui/open-webui:main",
							Env: []v1.EnvVar{
								{
									Name:  "OLLAMA_BASE_URL",
									Value: ollamaServerURI,
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
					RestartPolicy: v1.RestartPolicyAlways,
					Containers: []v1.Container{
						{
							Name:  "ollama",
							Image: "ollama/ollama:latest",
							Ports: []v1.ContainerPort{{ContainerPort: port}},
							TTY:   true,
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      ollamaVolumeMountName,
									MountPath: "/root/.ollama",
								},
							},
						},
					},
				},
			},
		},
	}
}
