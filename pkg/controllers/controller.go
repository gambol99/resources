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

package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/record"

	apiv1 "github.com/gambol99/resources/pkg/apis/resources/v1"
	"github.com/gambol99/resources/pkg/client/clientset/versioned"
	"github.com/gambol99/resources/pkg/cloud/aws"
	"github.com/gambol99/resources/pkg/cloud/null"
	"github.com/gambol99/resources/pkg/controllers/api"
	"github.com/gambol99/resources/pkg/controllers/cleanup"
	"github.com/gambol99/resources/pkg/controllers/resources"
	"github.com/gambol99/resources/pkg/controllers/templates"
	"github.com/gambol99/resources/pkg/models"
	"github.com/gambol99/resources/pkg/version"
)

// ResourceController maintains the state
type ResourceController struct {
	client    kubernetes.Interface
	clientset versioned.Interface
	cloud     models.CloudProvider
	config    *api.Config
	election  api.Leadership
	recorder  record.EventRecorder
	routines  []api.Controller
}

// New creates and returns a new controller
func New(config *api.Config) (*ResourceController, error) {
	log.Infof("starting the %s controller, version: %s", apiv1.GroupName, version.GetVersion())
	// @step: set the logger level
	if config.Verbose {
		log.SetLevel(log.DebugLevel)
	}

	return &ResourceController{config: config}, nil
}

// Run starts the controller runtime
func (r *ResourceController) Run(ctx context.Context) error {
	var err error

	// @step: create a cloud provider if required
	if r.cloud == nil {
		log.Infof("initializing the cloud provider: %s", r.config.CloudProvider)
		if r.cloud, err = makeCloudProvider(r.config.CloudProvider, &models.ProviderConfig{
			ClusterName: r.config.ClusterName,
			Name:        r.config.Name,
		}); err != nil {
			return fmt.Errorf("unable to initialize cloud provider: %s", err)
		}
	}
	// @step: create the kubernetes client apis
	if r.client == nil {
		log.Info("initializing the kubernetes in-cluster client api")
		if r.client, err = makeKubernetesClient(r.config); err != nil {
			return fmt.Errorf("unable to initialize a kubernetes api client: %s", err)
		}
	}
	if r.clientset == nil {
		log.Info("initializing the resource in-cluster client api")
		if r.clientset, err = makeResourceKubernetesClient(r.config); err != nil {
			return fmt.Errorf("unable to initialize a resources api client: %s", err)
		}
	}
	if r.config.EnableMetrics {
		if err := makeMetricsEndpoint(r.config); err != nil {
			return fmt.Errorf("unable to create metrics endpoint: %s", err)
		}
	}

	if r.recorder == nil {
		bc := record.NewBroadcaster()
		bc.StartRecordingToSink(&core.EventSinkImpl{Interface: r.client.CoreV1().Events("")})
		r.recorder = bc.NewRecorder(scheme.Scheme, v1.EventSource{Component: ""})
	}

	// @step: create the election leader
	endpoint := "cloud.appvia.io"
	log.Infof("initializing the controller election, namespace: %s, endpoint: %s", r.config.ElectionNamespace, endpoint)
	r.election, err = newElection(r.client, r.recorder, endpoint, r.config.ElectionNamespace)
	if err != nil {
		return fmt.Errorf("unable to create controller election: %s", err)
	}

	options := &api.Options{
		Client:         r.client,
		Cloud:          r.cloud,
		Config:         r.config,
		Election:       r.election,
		Record:         r.recorder,
		ResourceClient: r.clientset,
		Threadness:     r.config.Threadness,
	}

	// @step: create the controllers
	cleanup, err := cleanup.New(options)
	if err != nil {
		return fmt.Errorf("unable to create cleanup controller: %s", err)
	}
	resourcesCtrl, err := resources.New(options)
	if err != nil {
		return fmt.Errorf("unablr to create the cloud resources controller: %s", err)
	}
	templatesCtrl, err := templates.New(options)
	if err != nil {
		return fmt.Errorf("unable to create the cloud templates controller: %s", err)
	}
	r.routines = []api.Controller{cleanup, resourcesCtrl, templatesCtrl}

	var errorCh chan error

	// @step; start all the controllers
	for _, c := range r.routines {
		go func(ctl api.Controller) {
			log.WithFields(log.Fields{
				"controller": ctl.Name(),
			}).Info("starting the controller service processor")
			if err := ctl.Run(ctx); err != nil {
				errorCh <- fmt.Errorf("failed to start controller %s, error: %s", ctl.Name(), err)
			}
		}(c)
	}

	// @step: wait for a stop signal
	log.Info("management controller waiting for termination signal")
	select {
	case <-ctx.Done():
		log.Info("management controller recieved a termination signal")
	case err := <-errorCh:
		log.Errorf("failed to start the controller: %s", err)
		return err
	}

	return nil
}

// Wait waits on the controllers to finish
func (r *ResourceController) Wait(timeout time.Duration) error {
	doneCh := make(chan struct{}, 0)

	go func() {
		for _, x := range r.routines {
			log.WithFields(log.Fields{
				"controller": x.Name(),
			}).Info("waiting for controller tasks to gracefully finished")
			x.Wait()
		}
		doneCh <- struct{}{}
	}()

	select {
	case <-time.After(timeout):
		return fmt.Errorf("waiting for graceful shutdown of controllers has timed after: %s", timeout)
	case <-doneCh:
	}

	return nil
}

// makeMetricsnEndpoint creates the metrics endpoint
func makeMetricsEndpoint(config *api.Config) error {
	log.WithFields(log.Fields{
		"listen": config.MetricsListen,
	}).Info("exposing prometheus metrics for controllers")

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		if err := http.ListenAndServe(config.MetricsListen, nil); err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("")
		}
	}()

	return nil
}

// makeKubernetesConfig is responsible for getting either the in-cluster config of kubeconfig
func makeKubernetesConfig(config *api.Config) (*rest.Config, error) {
	if config.KubeConfig != "" {
		return clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	}

	return rest.InClusterConfig()
}

// makeResourceKubernetesClient returns a versioned kubernetes api client
func makeResourceKubernetesClient(config *api.Config) (versioned.Interface, error) {
	cfg, err := makeKubernetesConfig(config)
	if err != nil {
		return nil, err
	}

	return versioned.NewForConfig(cfg)
}

// makeKubernetesClient is responsible for initializing a kubernetes api clients
func makeKubernetesClient(config *api.Config) (kubernetes.Interface, error) {
	cfg, err := makeKubernetesConfig(config)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(cfg)
}

// makeCloudProvider is responsible for creating a cloud provider
func makeCloudProvider(name string, config *models.ProviderConfig) (models.CloudProvider, error) {
	switch name {
	case "aws":
		return aws.New(config)
	case "null":
		return null.New(config)
	default:
		return nil, errors.New("unknown of unsupport cloud provider")
	}
}
