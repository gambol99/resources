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

package templates

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	apiv1 "github.com/gambol99/resources/pkg/apis/resources/v1"
	inform "github.com/gambol99/resources/pkg/client/informers/externalversions/resources/v1"
	"github.com/gambol99/resources/pkg/controllers/api"
	"github.com/gambol99/resources/pkg/models"
	"github.com/gambol99/resources/pkg/utils"
)

// the controller is used to monitor the changes in cloud templates
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

// updated is called when a template has been updated or created
func (c *controller) updated(template *apiv1.CloudTemplate) error {
	log.WithFields(log.Fields{
		"name": template.Name,
	}).Info("checking the cloud template is valid")

	// @check the template is valid and if not we need to update the status
	switch errs := template.IsValid(); len(errs) > 0 {
	case true:
		template.Status = apiv1.TemplateSpecStatus{
			Message: "The cloud template specification is invalid",
			Reason:  utils.GetErrors(errs).Error(),
			Status:  models.StatusTemplateInvalid,
		}
	default:
		template.Status = apiv1.TemplateSpecStatus{Status: models.StatusTemplateOK}
	}

	log.WithFields(log.Fields{
		"name":   template.Name,
		"status": template.Status.Status,
		"reason": template.Status.Reason,
	}).Debug("updating the cloud template status")

	// @step: attempt to update the template statue
	return utils.Retry(5, time.Duration(time.Second*3), func() error {
		_, err := c.options.ResourceClient.CloudV1().CloudTemplates().Update(template)
		if err != nil {
			log.WithFields(log.Fields{
				"error":    err.Error(),
				"template": template.Name,
			}).Error("failed to update the template status")
		}

		return err
	})
}

// Run is responsible for starting the controller up
func (c *controller) Run(ctx context.Context) error {
	// @step: we create a namespace informer
	c.informer = inform.NewCloudTemplateInformer(c.options.ResourceClient, c.options.ResyncDuration, cache.Indexers{})
	c.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				c.queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err == nil {
				c.queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {},
	})
	defer c.queue.ShutDown()

	// @step: start the shared index informer
	stopCh := make(chan struct{}, 0)
	go c.informer.Run(stopCh)

	log.WithFields(log.Fields{"controller": c.Name()}).Info("waiting for controller caches to synchronize")
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("%s controller timed out waiting for caches to sync", c.Name()))
		return fmt.Errorf("%s controller timed out waiting for cache sync", c.Name())
	}

	// @step: start the workers
	for i := 0; i < c.options.Threadness; i++ {
		go c.processItems()
	}
	// @step: wait for a signal to stop
	select {
	case <-ctx.Done():
		close(stopCh)
	}
	// @step: shutdown the queue and the informer
	log.WithFields(log.Fields{"controller": c.Name()}).Info("shutting down the controller")

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
	} else if c.queue.NumRequeues(key) < c.options.MaxRetries {
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
		log.WithFields(log.Fields{
			"controller": c.Name(),
		}).Debug("skipping the maintenance, controller not the leader")

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
	_, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// @check if we are deleting the resource
	if !exists {
		log.WithFields(log.Fields{
			"name": name,
		}).Info("cloud resource template has been remove")

		return nil
	}

	resource, ok := obj.(*apiv1.CloudTemplate)
	if !ok {
		return fmt.Errorf("object should have been a cloudresource")
	}

	return c.updated(resource)
}

// Name returns the name of the controller
func (c *controller) Name() string {
	return "templates"
}

// Wait returns the task group stopped
func (c *controller) Wait() {
	c.waitgroup.Wait()
}
