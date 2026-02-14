/*
Copyright 2026.

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

// ResourceType represents a Kubernetes resource type to monitor
// +kubebuilder:validation:Enum=Secret;ConfigMap
type ResourceType string

const (
	// ResourceTypeSecret represents Secret resources
	ResourceTypeSecret ResourceType = "Secret"
	// ResourceTypeConfigMap represents ConfigMap resources
	ResourceTypeConfigMap ResourceType = "ConfigMap"
)

// OrphanagePolicySpec defines the desired state of OrphanagePolicy.
type OrphanagePolicySpec struct {
	// ResourceTypes specifies the Kubernetes resource types to monitor
	// Supported values: "Secret", "ConfigMap"
	ResourceTypes []ResourceType `json:"resourceTypes,omitempty"`
}

// Orphan represents an orphaned resource
type Orphan struct {
	// Kind is the Kubernetes resource kind (e.g., "Secret", "ConfigMap")
	Kind string `json:"kind"`
	// Name is the name of the orphaned resource
	Name string `json:"name"`
}

// OrphanagePolicyStatus defines the observed state of OrphanagePolicy.
type OrphanagePolicyStatus struct {
	// OrphanCount is the total number of orphaned resources
	OrphanCount int `json:"orphanCount,omitempty"`
	// LastChanged is the timestamp when the status was last updated
	LastChanged metav1.Time `json:"lastChanged,omitempty"`
	// Orphans is the list of orphaned resources
	Orphans []Orphan `json:"orphans,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// OrphanagePolicy is the Schema for the orphanagepolicies API.
type OrphanagePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OrphanagePolicySpec   `json:"spec,omitempty"`
	Status OrphanagePolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OrphanagePolicyList contains a list of OrphanagePolicy.
type OrphanagePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OrphanagePolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OrphanagePolicy{}, &OrphanagePolicyList{})
}
