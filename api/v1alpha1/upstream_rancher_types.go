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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// UpstreamRancherSpec defines the desired state of UpstreamRancher
type UpstreamRancherSpec struct {
	// Important: Run "make" to regenerate code after modifying this file

	// Name is the name of the rancher
	Name string `json:"name"`

	// KubeConfig is a kubeconfig of the rancher's local account
	KubeConfig KubeConfig `json:"kubeConfig"`
}

// RancherStatus defines the observed state of UpstreamRancher
type RancherStatus struct {
	// Important: Run "make" to regenerate code after modifying this file

	// State is the state of the account
	State string `json:"state"`
}

type KubeConfig struct {
	Value     string       `json:"value,omitempty"`
	ValueFrom *ValueSource `json:"valueFrom,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// UpstreamRancher is the Schema for the ranchers API
type UpstreamRancher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UpstreamRancherSpec `json:"spec,omitempty"`
	Status RancherStatus       `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// UpstreamRancherList contains a list of UpstreamRancher
type UpstreamRancherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UpstreamRancher `json:"items"`
}

func init() {
	SchemeBuilder.Register(&UpstreamRancher{}, &UpstreamRancherList{})
}
