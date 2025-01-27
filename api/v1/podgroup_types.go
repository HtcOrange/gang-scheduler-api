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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PodGroupSpec defines the desired state of PodGroup.
type PodGroupSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// GroupName is the name of PodGroup
	GroupName string `json:"group_name,omitempty"`
	// MinNum is the min number of the pods created that need to be scheduled
	MinNum int `json:"min_num,omitempty"`
	// Template is the pod tempalte
	Template corev1.PodTemplateSpec `json:"template,omitempty"`
}

// PodGroupStatus defines the observed state of PodGroup.
type PodGroupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Phase indicates current state of PodGroup
	Phase string `json:"phase,omitempty"`

	// PodMetaList contains pods information in the same PodGroup
	PodMetaList []PodMeta `json:"pod_meta_list,omitempty"`

	// Message containers additional information of PodGroup
	Message string `json:"message,omitempty"`
}

// PodMeta is the meta informatino of one pod
type PodMeta struct {
	// Name is the name of pod
	Name string `json:"name,omitempty"`

	// Ip is podIP in a k8s cluster, a pod can find other pods ips in the same PodGroup by this
	IP string `json:"ip,omitempty"`

	// HostIP is the node ip that pod started at
	HostIP string `json:"host_ip,omitempty"`

	// Phase of pod
	Phase string `json:"phase,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// PodGroup is the Schema for the podgroups API.
type PodGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodGroupSpec   `json:"spec,omitempty"`
	Status PodGroupStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PodGroupList contains a list of PodGroup.
type PodGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodGroup{}, &PodGroupList{})
}
