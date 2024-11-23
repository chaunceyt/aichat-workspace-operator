package k8s

import (
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	defaultNameLabel     = "app.kubernetes.io/name"
	defaultInstanceLabel = "app.kubernetes.io/instance"
)

func NewDeployment(namespace, name string, port int32) *appsv1.Deployment {
	appLabels := map[string]string{defaultNameLabel: name}
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
							Name:  "nginx",
							Image: "nginx:1.27.2",
							Ports: []v1.ContainerPort{{ContainerPort: port}},
						},
					},
				},
			},
		},
	}
}

func NewStatefulSet(namespace, name string, port int32) *appsv1.StatefulSet {
	appLabels := map[string]string{defaultNameLabel: name}
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{MatchLabels: appLabels},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: appLabels},
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyAlways,
					Containers: []v1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.27.2",
							Ports: []v1.ContainerPort{{ContainerPort: port}},
						},
					},
				},
			},
		},
	}
}
