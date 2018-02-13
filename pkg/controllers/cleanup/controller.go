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

package cleanup

import (
	"context"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/gambol99/resources/pkg/controllers/api"
	"github.com/gambol99/resources/pkg/models"
	"github.com/gambol99/resources/pkg/utils"
)

// controller is responsible for cycling through the resources
// and deleting and cloudresources marked for deletion
type controller struct {
	// config are the controller config
	config *api.Config
	// options are the controller options
	options *api.Options
	// waitgroup for the tasks
	waitgroup *sync.WaitGroup
}

// New creates and returns a maintenance controller
func New(options *api.Options) (api.Controller, error) {
	return &controller{
		config:    options.Config,
		options:   options,
		waitgroup: &sync.WaitGroup{},
	}, nil
}

// Run is responsible for kicking off the controller
func (c *controller) Run(ctx context.Context) error {
	log.Info("starting the maintenance controller")
	// @step: start the service in the background
	go c.processItems(ctx)
	// @step: wait for an exit signal
	<-ctx.Done()

	return nil
}

// processEvent is responsible for performing the maintenace task
func (c *controller) processEvent(ctx context.Context) error {
	cleanupCounter.Inc()

	c.waitgroup.Add(1)
	defer c.waitgroup.Done()

	// @step: we need to check if we are the leader
	if leader := c.options.IsLeader(); !leader {
		log.WithFields(log.Fields{
			"controller": c.Name(),
		}).Debug("skipping the maintenance, controller not the leader")

		return nil
	}

	timer := prometheus.NewTimer(cleanupDuration)
	defer timer.ObserveDuration()

	log.WithFields(log.Fields{
		"controller": c.Name(),
	}).Debug("running the maintenance service")

	err := func() error {
		// @step: retrieve a list of stacks from the cloud provider
		stacks, err := c.options.Cloud.List(ctx, &models.ListOptions{})
		if err != nil {
			return err
		}

		for _, x := range stacks {
			log.WithFields(log.Fields{
				"deleting":  x.HasDeleteTag(),
				"namespace": x.Namespace,
				"resource":  x.Spec.Name,
				"stack":     x.Name,
			}).Debug("checking if stack is up for deletion")

			if !x.HasDeleteTag() {
				continue
			}

			// @step: check is the stack is not yet up for deletion
			if !x.RequiresDeletion() {
				log.WithFields(log.Fields{
					"expires":   x.ExpiresIn(),
					"namespace": x.Namespace,
					"resource":  x.Spec.Name,
					"stack":     x.Name,
				}).Info("stack scheduled for but up for deletion yet")

				continue
			}

			// @check if the stack is a in deletion failed state
			switch x.Status.Status {
			case models.StatusDone, models.StatusFailed:
			default:
				log.WithFields(log.Fields{
					"name":     x.Name,
					"reason":   x.Status.Reason,
					"resource": x.Spec.Name,
					"stack":    x.Name,
					"status":   x.Status.Status,
				}).Info("refusing to delete the stack due to current failed state")

				continue
			}

			log.WithFields(log.Fields{
				"name":      x.Name,
				"namespace": x.Namespace,
				"resource":  x.Spec.Name,
				"template":  x.Spec.Template,
			}).Info("stack is scheduled for deletion, deleting now")

			utils.Retry(3, time.Second*10, func() error {
				if err := c.options.Cloud.Delete(ctx, x.Name, nil); err != nil {
					log.WithFields(log.Fields{
						"error":    err.Error(),
						"resource": x.Spec.Name,
						"stack":    x.Name,
					}).Error("unable to delete the stack")

					return err
				}

				return utils.DeleteCloudStatus(c.options.ResourceClient, x.Spec.Name, x.Namespace)
			})
		}

		return nil
	}()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("failed to perform an cleanup run")

		return err
	}

	return nil
}

// processItems is the main entrypoint for the service loop
func (c *controller) processItems(ctx context.Context) {
	// @step: we wait for rotation or the signal to quit
	ticker := time.NewTicker(time.Second * 30)

	defer func() {
		log.Info("maintenance controller recieved signal to terminate the controller")
	}()

	// @step: we wait until a ticker of the maintenance
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.processEvent(ctx); err != nil {
				cleanupErrors.Inc()
			}
		}
	}
}

// Name returns the name of the controller
func (c *controller) Name() string {
	return "cleanup"
}

// Wait returns the task group stopped
func (c *controller) Wait() {
	c.waitgroup.Wait()
}
