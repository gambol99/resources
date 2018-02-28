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

package resources

import (
	"context"
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	apiv1 "github.com/gambol99/resources/pkg/apis/resources/v1"
	inform "github.com/gambol99/resources/pkg/client/informers/externalversions/resources/v1"
	"github.com/gambol99/resources/pkg/controllers/api"
)

// the namespace controller is used to monitor the changes in namespaces and
// update the resources appropreiately
type controller struct {
	// informer is the lister
	informer cache.SharedIndexInformer
	// the worker queue
	queue workqueue.RateLimitingInterface
	// config are the controller config
	config *api.Config
	// options are the controller options
	options *api.Options
	// waitgroup is a wait group for the workers
	waitgroup *sync.WaitGroup
}

// New returns a new namespace controller
func New(options *api.Options) (api.Controller, error) {
	return &controller{
		config:    options.Config,
		options:   options,
		waitgroup: &sync.WaitGroup{},
		queue:     workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
	}, nil
}

// Run is responsible for starting the controller up
func (c *controller) Run(ctx context.Context) error {
	log.Infof("starting the %s controller, used to handle the cloud resources", c.Name())

	// @step: we create a namespace informer
	c.informer = inform.NewCloudResourceInformer(c.options.ResourceClient, "", c.options.ResyncDuration, cache.Indexers{})
	c.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				c.queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				c.queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(newObj)
			if err == nil {
				c.queue.Add(key)
			}
		},
	})
	defer c.queue.ShutDown()

	// @step: start the shared index informer
	var stopCh chan struct{}
	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("%s controller timed out waiting for caches to sync", c.Name()))
		return fmt.Errorf("%s controller timed out waiting for cache sync", c.Name())
	}

	// @step: start the workers
	for i := 0; i < c.options.Threadness; i++ {
		go c.processItems()
	}
	// @step: wait for a signal to stop
	<-ctx.Done()
	// @step: shutdown the queue and the informer
	close(stopCh)

	return nil
}

// proccessItems is responsible for the service loop
func (c *controller) processItems() {
	for c.processNextItem() {
	}
}

// processNextItem is the responsible for pulling items off the queue
func (c *controller) processNextItem() bool {
	c.waitgroup.Add(1)
	defer c.waitgroup.Done()
	// @step: pop the task off the queue
	key, quit := c.queue.Get()
	if quit {
		return false
	}
	defer c.queue.Done(key)

	// @step: process the event in the worker queue
	err := c.processEvent(key.(string))
	if err == nil {
		c.queue.Forget(key)
		return true
	}
	metricErrorTotal.Inc()

	if c.queue.NumRequeues(key) < c.options.MaxRetries {
		c.queue.AddRateLimited(key)
	} else {
		c.queue.Forget(key)
		runtime.HandleError(err)
	}

	return true
}

// processEvent is where the action happend the method is responsible for conversion what is
// desired (defined by changes in the custom resources) with the actual backend
func (c *controller) processEvent(key string) error {
	if leader := c.options.IsLeader(); !leader {
		log.Debug("skipping the maintenance, controller not the leader")

		return nil
	}

	obj, exists, err := c.informer.GetIndexer().GetByKey(key)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"key":   key,
		}).Errorf("failed to fetch object with key from store")

		return err
	}

	// @step: grab the name and namespace of the object
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// @check if we are deleting the resource
	if !exists {
		return c.deleted(name, namespace)
	}

	resource, ok := obj.(*apiv1.CloudResource)
	if !ok {
		return fmt.Errorf("object should have been a cloudresource")
	}

	return c.updated(resource)
}

// Name returns the name of the controller
func (c *controller) Name() string {
	return "resources"
}

// Wait returns the task group stopped
func (c *controller) Wait() {
	c.waitgroup.Wait()
}
