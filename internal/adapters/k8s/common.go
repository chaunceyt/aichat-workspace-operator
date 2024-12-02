package k8s

import (
	"fmt"

	kedahttpv1alpha1 "github.com/kedacore/http-add-on/operator/apis/http/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/chaunceyt/aichat-workspace-operator/internal/adapters/utils"
)

// NewNamespace returns new K8S namespace
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

// NewServiceAccount returns new K8S service account
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

func NewService(namespace string, name string, port int32, appLabels map[string]string) *corev1.Service {
	return &corev1.Service{
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

// need to make scale-to-zero work.
// the ingress for open-webui and ollama will point to these
func NewExternalService(namespace string, name, workload string, port int32, appLabels map[string]string) *corev1.Service {
	// TODO: move to const package.
	kedaProxy := "keda-add-ons-http-interceptor-proxy.keda"
	proxyName := fmt.Sprintf("%s-%s", namespace, workload)
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getName(proxyName, "-proxy"),
			Namespace: namespace,
			Labels:    appLabels,
		},
		Spec: corev1.ServiceSpec{
			Type:         corev1.ServiceTypeExternalName,
			ExternalName: kedaProxy,
		},
	}
}

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
				"pods":                   resource.MustParse("2"),
				"persistentvolumeclaims": resource.MustParse("2"),
				"services":               resource.MustParse("5"),
			},
		},
	}
}

func NewIngress(workspacename, ns, backendName, workload string, backendPort int) *networkingv1.Ingress {
	return &networkingv1.Ingress{
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
					Host: setIngressDNSHost(workspacename, workload),
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path: "/",
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: getName(workspacename, workload),
											Port: networkingv1.ServiceBackendPort{
												Number: int32(80),
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

func NewHttpSo(name, namespace, kind, workload string, port int32, hosts []string) *kedahttpv1alpha1.HTTPScaledObject {
	return &kedahttpv1alpha1.HTTPScaledObject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getName(name, workload),
			Namespace: namespace,
		},
		Spec: kedahttpv1alpha1.HTTPScaledObjectSpec{
			Hosts:        hosts,
			PathPrefixes: []string{"/"},
			ScaleTargetRef: kedahttpv1alpha1.ScaleTargetRef{
				Name:       getName(name, workload),
				Kind:       kind,
				APIVersion: "apps/v1",
				Service:    getName(name, workload),
				Port:       int32(port),
			},
			Replicas: &kedahttpv1alpha1.ReplicaStruct{
				Min: utils.PtrInt32(0),
				Max: utils.PtrInt32(1),
			},
			ScalingMetric: &kedahttpv1alpha1.ScalingMetricSpec{
				Rate: &kedahttpv1alpha1.RateMetricSpec{
					TargetValue: 20,
				},
			},
		},
	}
}

// setIngressDNSHost: team-a-aichat.openwebui.localtest.me
func setIngressDNSHost(workspace string, workload string) string {
	defaultDomain := "localtest.me"
	dnsName := fmt.Sprintf("%s.%s.%s", workspace, workload, defaultDomain)
	return dnsName
}

func getName(workspace, workload string) string {
	name := fmt.Sprintf("%s-%s", workspace, workload)
	return name
}
