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
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gambol99/resources/pkg/models"
	"github.com/gambol99/resources/pkg/utils"
)

// deleted is responsible for handling the removal of a cloudresource
func (c *controller) deleted(name, namespace string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*30)
	defer cancel()

	// @step: pull the stack from the cloud provider
	stackname := getStackName(name, namespace)
	stack, err := c.options.Cloud.Get(ctx, stackname, &models.GetOptions{})
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err.Error(),
			"namespace": namespace,
			"name":      name,
		}).Error("failed to retrieve the stack from cloud provider")

		return err
	}

	log.WithFields(log.Fields{
		"created":   stack.Created.Format(time.RFC822Z),
		"name":      name,
		"namespace": namespace,
		"retention": stack.Spec.Retention,
		"template":  stack.Spec.Template,
	}).Info("cloud resource stack deletion event")

	// @check if not retention, in which case we can delete straight away
	if stack.Spec.Retention <= 0 {
		log.WithFields(log.Fields{
			"name":      name,
			"namespace": namespace,
			"template":  stack.Spec.Template,
		}).Info("deleting stack as it has not retention period")

		return utils.Retry(3, time.Second*10, func() error {
			err := c.options.Cloud.Delete(ctx, stack.Name, &models.DeleteOptions{})
			if err != nil {
				log.WithFields(log.Fields{
					"error":     err.Error(),
					"name":      name,
					"namespace": namespace,
				}).Error("failed to delete the stack")

				return err
			}

			return utils.DeleteCloudStatus(c.options.ResourceClient, name, namespace)
		})
	}

	// @logic to need to update tags for set a maintenance deletion
	expiration := time.Now().Add(stack.Spec.Retention)
	log.WithFields(log.Fields{
		"expiry":    time.Now().Add(stack.Spec.Retention).String(),
		"name":      name,
		"namespace": namespace,
		"template":  stack.Spec.Template,
	}).Info("attempting to mark stack deletion later")

	stack.Spec.Tags[models.DeletionTimeTag] = fmt.Sprintf("%d", expiration.Unix())

	if err := c.options.Cloud.UpdateTags(ctx, stackname, stack.Spec.Tags); err != nil {
		log.WithFields(log.Fields{
			"error":     err.Error(),
			"name":      name,
			"namespace": namespace,
		}).Error("unable to update the stack tags")

		return err
	}

	return nil
}
