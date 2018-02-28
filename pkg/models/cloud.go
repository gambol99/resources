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

package models

import (
	"context"
	"errors"
	"time"

	apiv1 "github.com/gambol99/resources/pkg/apis/resources/v1"
)

var (
	// ErrOperationAborted indicates the operation was aborted
	ErrOperationAborted = errors.New("operation aborted")
	// ErrStackNotFound indicates the resource stack was not found
	ErrStackNotFound = errors.New("stack not found")
	// ErrUnauthorized indicates you trying to delete a stacks was you do not own
	ErrUnauthorized = errors.New("unauthorized to operate on this stack")
)

// ProviderConfig are configuration options for the providers
type ProviderConfig struct {
	// ClusterName is the name of the cluster
	// +required
	ClusterName string
	// Regon is region the cluster resides
	// +optional
	Region string
	// Name is the name of the provider
	// +required
	Name string
}

// CreateOptions is a set of providers for the provider
type CreateOptions struct {
	// Context is a set of contextual values
	// +required
	Context map[string]string
	// Resource is the resource we are creating
	// +required
	Resource *apiv1.CloudResource
	// Tags is a collection of tags to use when creating the resource
	// +optional
	Tags map[string]string
	// Template is the templete it's based off
	// +required
	Template *apiv1.CloudTemplate
	// WaitOn indicates we should wait on the creation
	// +optional
	WaitOn bool
}

// DeleteOptions is the delete options
type DeleteOptions struct {
	// WaitOn indicates we should wait on the creation
	// +optional
	WaitOn bool
}

// GetOptions are options for get options
type GetOptions struct{}

// WaitOptions are options for wait operations
type WaitOptions struct {
	// CheckInterval is the duration to wait before checking
	// +optional
	CheckInterval time.Duration
}

// ListOptions are used for list options
type ListOptions struct {
	// Before are any stacks created before this time
	// +optional
	Before *time.Time
	// After is any stacks create after this time
	// +optional
	After *time.Time
	//
}

// CredentialsOptions are the options for creating creational
type CredentialsOptions struct {
}

// CloudProvider defined the cloud provider contract
type CloudProvider interface {
	// Credentials generates the credentials from a stack
	Credentials(context.Context, string) ([]Credential, error)
	// Create is responsible for creating or updating a stack
	Create(context.Context, string, *CreateOptions) error
	// Delete is responsible for removing the stack
	Delete(context.Context, string, *DeleteOptions) error
	// Exists is responsible for checking is stack already exists
	Exists(context.Context, string) (*Stack, bool, error)
	// Get is responisble for retrieving a stack
	Get(context.Context, string, *GetOptions) (*Stack, error)
	// List is responsible for getting a list of stacks
	List(context.Context, *ListOptions) ([]*Stack, error)
	// Logs gets the logs on the stack
	Logs(context.Context, string, *GetOptions) (string, error)
	// Status is responsible for getting the status
	Status(context.Context, string, *GetOptions) (string, error)
	// UpdateTags is responsible for updating just the tags of a stack
	UpdateTags(context.Context, string, map[string]string) error
	// Wait is responsible for waiting for a stack to complete or fail
	Wait(context.Context, string, *WaitOptions) (string, error)
}

const (
	// StatusDone indicates the operation is complete
	StatusDone = "OK"
	// StatusFailed indicates the stack has failed
	StatusFailed = "Failed"
	// StatusDeleting indicates the stack is being deleted
	StatusDeleting = "Deleting"
	// StatusInProgress indicates the stack is in progress
	StatusInProgress = "InProgress"
	// StatusInRollback indicates the stack is in rollback
	StatusInRollback = "Rollback"
	// StatusUnknown indicates an unknown status
	StatusUnknown = ""
	// StatusTemplateOK indicates the template has passed validation
	StatusTemplateOK = "OK"
	// StatusTemplateInvalid indicates the template is invalid
	StatusTemplateInvalid = "Invalid"
)

const (
	// CreatedTag is when the resource was created
	CreatedTag = ProviderTag + "/created"
	// CheckSumTag is the checksum tag
	CheckSumTag = ProviderTag + "/checksum"
	// DeletionTimeTag is the time the resource is up for deletion
	DeletionTimeTag = ProviderTag + "/removal"
	// NamespaceTag is the namespace tag
	NamespaceTag = ProviderTag + "/namespace"
	// ProviderNameTag is the owner
	ProviderNameTag = ProviderTag + "/provider"
	// ProviderTag is the name of the provider
	ProviderTag = "resources.appvia.io"
	// ResourceNameTag is the resource tag
	ResourceNameTag = ProviderTag + "/resource"
	// RetentionTag is the tag used for the retention
	RetentionTag = ProviderTag + "/retention"
	// TemplateNameTag is the name template used to create it
	TemplateNameTag = ProviderTag + "/template"
)

// Stack is an instance of a resource in the cloud
type Stack struct {
	// Created is the time the resource was created
	Created time.Time `json:"created" yaml:"created"`
	// Name is the name of the stack
	Name string `json:"name" yaml:"name"`
	// Namespace is the namespace it was created
	Namespace string `json:"namespace" yaml:"namespace"`
	// Spec is the specification
	Spec StackSpec `json:"spec" yaml:"spec"`
	// Status the status of the stack
	Status StackStatus `json:"status" yaml:"status"`
}

// StackSpec is the specification for the for a stack
type StackSpec struct {
	// DeleteOn is the deletion time if it has one
	DeleteOn time.Time `json:"deleteOn" yaml:"deleteOn"`
	// Name is the name of the actual cloud resource
	Name string `json:"stackName" yaml:"stackName"`
	// Outputs the outputs from a stack
	Outputs map[string]string `json:"outputs" yaml:"outputs"`
	// Retention is the duration a stack shoult be kept
	Retention time.Duration `json:"retention" yaml:"retention"`
	// Tags is a series of tags for the stack
	Tags map[string]string `json:"tags" yaml:"tags"`
	// Template is the name of the template used to build off
	Template string `json:"template" yaml:"template"`
}

// StackStatus the status of a stack
type StackStatus struct {
	// Status the status of the stack
	Status string `json:"status" yaml:"status"`
	// Reasion is a reasion for the status
	Reason string `json:"reason" yaml:"reason"`
}

// Credential is response from a credential creation
type Credential struct {
	// ID is the name of this credential
	ID string `json:"id"`
	// User is the username / id of the credential
	User string `json:"userID"`
	// Secret is the credential password
	Secret string `json:"secret"`
}
