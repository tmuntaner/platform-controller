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

// DownstreamRancherSpec defines the desired state of DownstreamRancher
type DownstreamRancherSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// UpstreamRancher is the name of the upstream rancher
	UpstreamRancher string `json:"upstreamRancher"`

	// Name is the name of the downstream rancher
	Name string `json:"name"`
}

// DownstreamRancherStatus defines the observed state of DownstreamRancher
type DownstreamRancherStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// State is the state of the account
	State string `json:"state"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DownstreamRancher is the Schema for the downstreamranchers API
type DownstreamRancher struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DownstreamRancherSpec   `json:"spec,omitempty"`
	Status DownstreamRancherStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DownstreamRancherList contains a list of DownstreamRancher
type DownstreamRancherList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DownstreamRancher `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DownstreamRancher{}, &DownstreamRancherList{})
}
