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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AIChatWorkspaceSpec defines the desired state of AIChatWorkspace.
type AIChatWorkspaceSpec struct {
	// The name of the workspace.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="WorkspaceName is immutable"
	WorkspaceName string `json:"workspaceName"`

	// +kubebuilder:validation:default:=dev
	WorkspaceEnv string `json:"workspaceENV"`

	// List of default models for this workspace.
	Models []string `json:"models"`
}

// AIChatWorkspaceStatus defines the observed state of AIChatWorkspace.
type AIChatWorkspaceStatus struct {
	IsCreated bool `json:"isCreated,omitempty"`

	// Represents the observations of a AIChatWorkspace's current state.
	// AIChatWorkspace.status.conditions.type are: "Available", "Progressing", and "Degraded"
	// AIChatWorkspace.status.conditions.status are one of True, False, Unknown.
	// AIChatWorkspace.status.conditions.reason the value should be a CamelCase string and producers of specific
	// condition types may define expected values and meanings for this field, and whether the values
	// are considered a guaranteed API.
	// AIChatWorkspace.status.conditions.Message is a human readable message indicating details about the transition.
	// For further information see: https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// AIChatWorkspace is the Schema for the aichatworkspaces API.
type AIChatWorkspace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AIChatWorkspaceSpec   `json:"spec,omitempty"`
	Status AIChatWorkspaceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AIChatWorkspaceList contains a list of AIChatWorkspace.
type AIChatWorkspaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AIChatWorkspace `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AIChatWorkspace{}, &AIChatWorkspaceList{})
}
