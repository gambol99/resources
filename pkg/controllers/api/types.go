/*
Copyright 2018 All rights reserved - Appvia.io

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

package api

import (
	"context"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"

	"github.com/gambol99/resources/pkg/client/clientset/versioned"
	"github.com/gambol99/resources/pkg/models"
)

// Config defines the configuraton for the controller
type Config struct {
	// CloudProvider is the actual provider i.e. aws
	CloudProvider string
	// ClusterName is the name of the cluster
	ClusterName string
	// EnableMetrics enables the metrics endpoint
	EnableMetrics bool
	// ElectionNamespace is the namespace for the endpoint election
	ElectionNamespace string
	// KubeConfig is an optional path to a kubeconfig file
	KubeConfig string
	// MetricsListen is the interface we should expose the metrics on
	MetricsListen string
	// Name is the name of the controller
	Name string
	// ResyncDuration is the default resync time duration for the controller
	ResyncDuration time.Duration
	// StackTimeout is the timeout for a stack to complete
	StackTimeout time.Duration
	// Threadness is the number of controller threads to run
	Threadness int
	// Verbose indicates verbose logging
	Verbose bool
}

// Leadership returns a true / false indicating leadership
type Leadership interface {
	// IsLeader checks if we are the leader
	IsLeader() bool
}

// Options are the options providers to a controller
type Options struct {
	Leadership

	// Client is the kubernetes client
	Client kubernetes.Interface
	// Cloud is the cloud provider
	Cloud models.CloudProvider
	// Config is the configuraton for the controller
	Config *Config
	// Election checks for leadership
	Election Leadership
	// Record is a event recorder
	Record record.EventRecorder
	// ResourceClient is the client for resources
	ResourceClient versioned.Interface
	// RsyncDuration is the duration for resyncing
	ResyncDuration time.Duration
	// Threadness is the number of workers for the controller
	Threadness int
	// MaxRetries is the max attempts to retry a event
	MaxRetries int
}

// IsLeader is just a friendly wrapper
func (o *Options) IsLeader() bool {
	return o.Election.IsLeader()
}

// Controller defined the controller contract
type Controller interface {
	// Name is the name of the controller
	Name() string
	// Run starts the controller
	Run(context.Context) error
	// Wait waits for tasks to finish
	Wait()
}
