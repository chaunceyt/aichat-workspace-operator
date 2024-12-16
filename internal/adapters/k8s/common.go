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

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"

	kedahttpv1alpha1 "github.com/kedacore/http-add-on/operator/apis/http/v1alpha1"

	"github.com/chaunceyt/aichat-workspace-operator/internal/constants"
)

/**
 * Creates a new Kubernetes Namespace object.
 *
 * @param name The name of the namespace to create.
 * @param appLabels A map of labels to apply to the namespace.
 * @return A pointer to a new corev1.Namespace object.
 */
func NewNamespace(name string, appLabels map[string]string) *corev1.Namespace {
	return &corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: appLabels,
		},
	}
}

/**
 * Creates a new Kubernetes ServiceAccount object.
 *
 * @param name The name of the service account to create.
 * @param namespace The namespace where the service account will be created.
 * @param appLabels A map of labels to apply to the service account.
 * @return A pointer to a new corev1.ServiceAccount object.
 */
func NewServiceAccount(name string, namespace string, appLabels map[string]string) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ServiceAccount",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    appLabels,
		},
	}
}

/**
 * Creates a new Kubernetes PersistentVolumeClaim object.
 *
 * @param name The name of the persistent volume claim to create.
 * @param namespace The namespace where the persistent volume claim will be created.
 * @param storageSize The size of storage requested for the persistent volume claim (e.g. "1Gi").
 * @param appLabels A map of labels to apply to the persistent volume claim.
 * @return A pointer to a new corev1.PersistentVolumeClaim object.
 */
func NewPersistentVolumeClaim(name string, namespace string, storageSize string, appLabels map[string]string) *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    appLabels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{

			AccessModes: []corev1.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},

			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(storageSize),
				},
			},
		},
	}
}

/**
 * Creates a new Kubernetes Service object.
 *
 * @param namespace The namespace where the service will be created.
 * @param name The name of the service to create.
 * @param port The port number that the service will listen on.
 * @param appLabels A map of labels to apply to the service.
 * @return A pointer to a new corev1.Service object.
 */
func NewService(namespace string, name string, port int32, appLabels map[string]string) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    appLabels,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app.kubernetes.io/name": name,
			},
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       port,
				TargetPort: intstr.FromInt32(port),
			}},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
}

/**
 * Creates a new Kubernetes ResourceQuota object.
 *
 * @param namespace The namespace where the resource quota will be created.
 * @param name The name of the resource quota to create.
 * @param appLabels A map of labels to apply to the resource quota.
 * @return A pointer to a new corev1.ResourceQuota object.
 */
func NewResourceQuota(namespace string, name string, appLabels map[string]string) *corev1.ResourceQuota {
	return &corev1.ResourceQuota{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ResourceQuota",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    appLabels,
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: corev1.ResourceList{
				"pods":                   resource.MustParse(constants.MaxPods),
				"persistentvolumeclaims": resource.MustParse(constants.MaxPersistentVolumeClaims),
				"services":               resource.MustParse(constants.MaxService),
			},
		},
	}
}

/**
 * Creates a new Kubernetes Ingress object.
 *
 * @param workspacename The name of the workspace where the ingress will be created.
 * @param workload The name of the workload to create an ingress for.
 * @param backendName The name of the service that the ingress will route traffic to.
 * @param hostname The hostname that the ingress will listen on (e.g. example.com).
 * @param backendPort The port number that the service is listening on.
 * @return A pointer to a new networkingv1.Ingress object.
 */
func NewIngress(workspacename, workload, backendName, hostname string, backendPort int32) *networkingv1.Ingress {
	pathType := networkingv1.PathTypePrefix
	return &networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      getName(workspacename, workload),
			Namespace: workspacename,
			Labels: map[string]string{
				"app.kubernetes.io/instance":  workspacename,
				"app.kubernetes.io/component": workspacename,
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{
				{
					Host: hostname,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: backendName,
											Port: networkingv1.ServiceBackendPort{
												Number: backendPort,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

/**
 * Creates a new Kubernetes Service object of type ExternalName.
 *
 * @param namespace The namespace where the service will be created.
 * @param appLabels A map of labels to apply to the service.
 * @return A pointer to a new corev1.Service object of type ExternalName.
 *
 * needed to make scale-to-zero work.
 * the ingress for open-webui and ollama will point to these
 */
func NewExternalService(namespace string, appLabels map[string]string) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      constants.ExternalServiceName,
			Namespace: namespace,
			Labels:    appLabels,
		},
		Spec: corev1.ServiceSpec{
			Type:         corev1.ServiceTypeExternalName,
			ExternalName: constants.KedaHttpInterceptorProxy,
		},
	}
}

/**
 * Creates a new Kubernetes HTTP Scaled Object.
 *
 * @param workspacename The name of the workspace where the scaled object will be created.
 * @param kind The kind of the workload to scale (e.g. "Deployment", "StatefulSet").
 * @param workload The name of the workload to create an HTTP scaled object for.
 * @param port The port number that the service is listening on.
 * @param hosts A list of hostnames that the scaled object will listen on.
 * @return A pointer to a new kedahttpv1alpha1.HTTPScaledObject object.
 */
func NewHttpSo(workspacename, kind, workload string, port int32, hosts []string) *kedahttpv1alpha1.HTTPScaledObject {
	return &kedahttpv1alpha1.HTTPScaledObject{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "http.keda.sh/v1alpha1",
			Kind:       "HTTPScaledObject",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      getName(workspacename, workload),
			Namespace: workspacename,
		},
		Spec: kedahttpv1alpha1.HTTPScaledObjectSpec{
			Hosts:        hosts,
			PathPrefixes: []string{"/"},
			ScaleTargetRef: kedahttpv1alpha1.ScaleTargetRef{
				Name:       getName(workspacename, workload),
				Kind:       kind,
				APIVersion: "apps/v1",
				Service:    getName(workspacename, workload),
				Port:       int32(port),
			},
			Replicas: &kedahttpv1alpha1.ReplicaStruct{
				Min: ptr.To[int32](0),
				Max: ptr.To[int32](1),
			},
			ScalingMetric: &kedahttpv1alpha1.ScalingMetricSpec{
				Rate: &kedahttpv1alpha1.RateMetricSpec{
					TargetValue: 20,
				},
			},
		},
	}
}

func getName(workspace, workload string) string {
	name := fmt.Sprintf("%s-%s", workspace, workload)
	return name
}
