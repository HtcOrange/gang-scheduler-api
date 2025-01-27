/*
Copyright 2024.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GangSpec defines the desired state of Gang.
type GangSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Replicas is the replica of Gang
	Replicas int32 `json:"replicas,omitempty"`

	// Gang labels
	Labels map[string]string `json:"labels,omitempty"`

	// PodGroup template
	Template PodGroupTemplate `json:"template,omitempty"`
}

type PodGroupTemplate struct {
	// PodGroup metadata
	Metadata metav1.ObjectMeta `json:"metadata,omitempty"`

	// PodGroup spec
	Spec PodGroupSpec `json:"spec"`
}

// GangStatus defines the observed state of Gang.
type GangStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// ready replicas
	ReadyReplicas int32 `json:"ready_replicas"`

	// conditions
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Gang is the Schema for the gangs API.
type Gang struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GangSpec   `json:"spec,omitempty"`
	Status GangStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GangList contains a list of Gang.
type GangList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gang `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Gang{}, &GangList{})
}
