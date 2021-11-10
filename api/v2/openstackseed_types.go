/*
Copyright 2021 SAP.

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

package v2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// OpenstackSeedSpec defines the desired state of OpenstackSeed
type OpenstackSeedSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of OpenstackSeed. Edit openstackseed_types.go to remove/update
	// list of required specs that need to be resolved before the current one
	Dependencies []string `json:"requires,omitempty" yaml:"requires,omitempty"`
	// list of keystone roles
	Roles []RoleSpec `json:"roles,omitempty" yaml:"roles,omitempty"`
	// list of implied roles
	RoleInferences []RoleInferenceSpec `json:"role_inferences,omitempty" yaml:"role_inferences,omitempty"`
	// list keystone regions
	Regions []RegionSpec `json:"regions,omitempty" yaml:"regions,omitempty"`
	// list keystone services and their endpoints
	Services []ServiceSpec `json:"services,omitempty" yaml:"services,omitempty"`
	// list of nova flavors
	Flavors []FlavorSpec `json:"flavors,omitempty" yaml:"flavors,omitempty"`
	// list of Manila share types
	ShareTypes []ShareTypeSpec `json:"share_types,omitempty" yaml:"share_types,omitempty"`
	// list of resource classes for the placement service (currently still part of nova)
	ResourceClasses []string `json:"resource_classes,omitempty" yaml:"resource_classes,omitempty"`
	// list keystone domains with their configuration, users, groups, projects, etc
	Domains []DomainSpec `json:"domains,omitempty" yaml:"domains,omitempty"`
	// list of neutron rbac polices (currently only network rbacs are supported)
	RBACPolicies []RBACPolicySpec `json:"rbac_policies,omitempty" yaml:"rbac_policies,omitempty"`
	// list of cinder volume types
	VolumeTypes []VolumeTypeSpec `json:"volume_types,omitempty" yaml:"volume_types,omitempty"`
}

// OpenstackSeedStatus defines the observed state of OpenstackSeed
type OpenstackSeedStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	UnfinishedSeeds           map[string]string `json:"unfinished_seeds,omitempty" yaml:"unfinished_seeds,omitempty"`
	ReconciledResourceVersion string            `json:"reconciled_resource_version,omitempty" yaml:"reconciled_resource_version,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// OpenstackSeed is the Schema for the openstackseeds API
type OpenstackSeed struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenstackSeedSpec   `json:"spec,omitempty"`
	Status OpenstackSeedStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OpenstackSeedList contains a list of OpenstackSeed
type OpenstackSeedList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenstackSeed `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenstackSeed{}, &OpenstackSeedList{})
}
