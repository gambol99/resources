/*
Copyright 2018 All rights reserved - Appvia

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

const (
	// GroupName is the API group name
	GroupName = "cloud.appvia.io"
	// GroupVersion is the version of the API
	GroupVersion = "v1"
)

const (
	// DeleteOnRetention indicates a holding period for deletion
	DeleteOnRetention = "retain"
	// DeleteNever indicates we do not delete at all
	DeleteNever = "never"
)

const (
	// SecretTypeOutput indicates an output
	SecretTypeOutput = "output"
	// SecretTypeCredential indicates a credential mapping
	SecretTypeCredential = "credential"
)

// Secret defines a mapping for a output to a kubernetes secret
type Secret struct {
	// Name is the name of the secret
	// +required
	Name string `json:"name" protobuf:"bytes,1,rep,name=name"`
	// Description is a short description of the parameter
	// +optional
	Description string `json:"description,omitempty" protobuf:"bytes,2,opt,name=description"`
	// Values provides the mapping to the secret
	// +required
	Values []SecretValue `json:"values" protobuf:"bytes,3,opt,name=values"`
}

// SecretValue defines the specification for a secret value
type SecretValue struct {
	// Type is type secret (output, credential)
	// +required
	Type string `json:"type" protobuf:"bytes,1,opt,name=type"`
	// Key is keyname for the kubernetes secret
	// +required
	Key string `json:"key" protobuf:"bytes,2,opt,name=key"`
	// Value is the mapping name to the type
	// +required
	Value string `json:"name" protobuf:"bytes,3,opt,name=name"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudResource is a generic type for a in-cloud templated resource
type CloudResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Spec is the specification of the resource
	Spec CloudResourceSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudStatus is a status object for the
type CloudStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// A human readable message indicating condition.
	// +optional
	Status string `json:"status,omitempty" protobuf:"bytes,2,opt,name=status"`
	// A human readable message indicating details about why the pod is in this condition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,3,opt,name=message"`
	// A brief CamelCase message indicating details about why the pod is in this state.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
	// Logs are the logs from from the stack
	// +optional
	Logs string `json:"logs,omitempty" protobuf:"bytes,5,opt,name=reason"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudStatusList is a list of CloudStatus's
type CloudStatusList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of CloudStatus
	Items []CloudStatus `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// Parameter defined a parameter for the resource
type Parameter struct {
	// Name is the key name of the parameter
	// +required
	Name string `json:"name,omitempty" protobuf:"bytes,1,req,name=name"`
	// Description is a short description of the parameter
	// +optional
	Description string `json:"description,omitempty" protobuf:"bytes,2,opt,name=description"`
	// SecretName is optional name of a secret holding the value
	// +optional
	SecretName *string `json:"secretName,omitempty" protobuf:"bytes,2,opt,name=secretName"`
	// Value is an optional of the value of the parameter
	// +optional
	Value *string `json:"value,omitempty" protobuf:"bytes,3,opt,name=value"`
}

// CloudResourceSpec is the definition for a requested cloud resource
type CloudResourceSpec struct {
	// Credentials indicates the template has credentials embedded
	// +optional
	Credentials bool `json:"credentials" protobuf:"bytes,2,rep,name=credentials"`
	// DeletedOn indicates the deletion policy
	// +optional
	DeleteOn *string `json:"deleteOn" protobuf:"bytes,1,opt,name=deletedOn"`
	// TemplateName is the name of the template to use
	// +required
	TemplateName string `json:"templateName,omitempty" protobuf:"bytes,2,rep,name=templateName"`
	// Retention is used with the deletion policy
	// +optional
	Retention *metav1.Duration `json:"retention,omitempty" protobuf:"bytes,3,opt,name=retention"`
	// Parameters is collection of parameters for this resource
	// +optional
	Parameters []Parameter `json:"parameters,omitempty" protobuf:"bytes,4,opt,name=parameters"`
	// Secrets is a mapping for outputs to kube secrets
	// +optional
	Secrets []Secret `json:"secrets,omitempty" protobuf:"bytes,5,ops,name=secrets,casttype=Secret"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudResourceList is a list of Resource items
type CloudResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of Resources
	Items []CloudResource `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudTemplate is a cloud resource template
type CloudTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Spec is the specification of the resource
	Spec TemplateSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	// Status the current state of the resource
	Status TemplateSpecStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CloudTemplateList is a list of Templte items
type CloudTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items is a list of Templates
	Items []CloudTemplate `json:"items" protobuf:"bytes,2,rep,name=items"`
}

const (
	// FormatYAML is the yaml template format
	FormatYAML = "yaml"
	// FormatJSON is the json format
	FormatJSON = "json"
)

// TemplateSpec defines the specification for a template
type TemplateSpec struct {
	// Content is the tempate content
	// +required
	Content string `json:"content" protobuf:"bytes,1,rep,name=content"`
	// Credentials indicates the template has credentials embedded
	// +optional
	Credentials bool `json:"credentials" protobuf:"bytes,2,rep,name=credentials"`
	// Format is the format of the template i.e. yaml or json
	// +required
	Format string `json:"format" protobuf:"bytes,4,req,name=format"`
	// Parameters is collection of parameters for this resource
	// +optional
	Parameters []Parameter `json:"parameters,omitempty" protobuf:"bytes,5,opt,name=parameters"`
	// Retention is used with the deletion policy
	// +optional
	Retention *metav1.Duration `json:"retention,omitempty" protobuf:"bytes,6,opt,name=retention"`
	// Secrets is a mapping for outputs to kube secrets
	// +optional
	Secrets []Secret `json:"secrets,omitempty" protobuf:"bytes,7,ops,name=secrets,casttype=Secret"`
}

// TemplateSpecStatus is the status information related to a template
type TemplateSpecStatus struct {
	// A human readable message indicating condition.
	// +optional
	Status string `json:"status,omitempty" protobuf:"bytes,1,opt,name=status"`
	// A human readable message indicating details about why the template is in this condition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,2,opt,name=message"`
	// A brief CamelCase message indicating details about why the template is in this state.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`
}
