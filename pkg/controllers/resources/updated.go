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
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	apiv1 "github.com/gambol99/resources/pkg/apis/resources/v1"
	"github.com/gambol99/resources/pkg/models"
	"github.com/gambol99/resources/pkg/utils"
)

// updated is responsible for updating / creating a resoruce
func (c *controller) updated(resource *apiv1.CloudResource) error {
	stackname := getStackName(resource.Name, resource.Namespace)

	// @step: check if the resource requires updating from the checksum

	// @step: attempt to retrieve the cloud template which this resource is built off
	template, err := utils.FindCloudTemplate(c.options.ResourceClient, resource.Spec.TemplateName)
	if err != nil {
		return fmt.Errorf("unable to retrieve cloud template: %s, error: %s", resource.Spec.TemplateName, err)
	}

	// @step: lets use a default 30 minutes for now
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*30)
	defer cancel()

	// @step: attempt to update the resource
	stack, result := c.updateCloudResource(ctx, stackname, resource, template)
	if result != nil {
		log.WithFields(log.Fields{
			"error":     result.Error(),
			"resource":  resource.Name,
			"namespace": resource.Namespace,
		}).Error("failed to update / create the cloud resource")
	}

	// @step: update the status of the status of the resource
	if err := c.updateCloudStatus(ctx, stack, result, resource); err != nil {
		log.WithFields(log.Fields{
			"error":     err.Error(),
			"resource":  resource.Name,
			"namespace": resource.Namespace,
		}).Error("failed to update the cloud status")

		return fmt.Errorf("failed to update the cloud status for stack: (%s/%s)", resource.Namespace, resource.Name)
	}

	// @step: if the result was an error we cannot proceed
	if result == nil {
		return result
	}

	// @step: if the stack has any credentials we need to generate them
	var credentials map[string]models.Credential

	if template.Spec.Credentials {
		credentials, err = c.updateCloudCredentials(ctx, resource, stack)
		if err != nil {
			return fmt.Errorf("unable to update / create credentials from stack: %s", err)
		}
	}

	// @step: we need to map the outputs, secrets and credentials into the user namespace
	if err := c.updateCloudSecrets(ctx, resource, stack, credentials); err != nil {
		return fmt.Errorf("unable to update the kubernetes secrets")
	}

	return nil
}

// updateCloudStatus is responsible for updating the cloud resource status
func (c *controller) updateCloudStatus(ctx context.Context, stack *models.Stack, errMsg error, resource *apiv1.CloudResource) error {
	status := &apiv1.CloudStatus{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.Name,
			Namespace: resource.Namespace,
		},
	}
	if errMsg != nil {
		status.Status = models.StatusFailed
		status.Message = "Failed to update / create the stack"
		status.Reason = errMsg.Error()

		// @check if we have a stack to update
		if stack == nil {
			return utils.UpdateCloudStatus(c.options.ResourceClient, status)
		}
		status.Status = fmt.Sprintf("%s", stack.Status)
	}
	status.Status = fmt.Sprintf("%s", stack.Status)

	// @step: grab the logs from the stack
	logs, err := c.options.Cloud.Logs(ctx, stack.Name, &models.GetOptions{})
	if err != nil {
		return utils.UpdateCloudStatus(c.options.ResourceClient, status)
	}
	status.Logs = fmt.Sprintf("|\n%s", logs)

	return utils.UpdateCloudStatus(c.options.ResourceClient, status)
}

// updateCloudResource is resposible for updating the resource
func (c *controller) updateCloudResource(ctx context.Context, stackname string, resource *apiv1.CloudResource, template *apiv1.CloudTemplate) (*models.Stack, error) {
	// @check if the stack already exists. It then checks the status of the stack
	// waiting on those which haven't finished yet
	stack, found, err := c.options.Cloud.Exists(ctx, stackname)
	if err != nil {
		return nil, fmt.Errorf("unable to check if stack exists already: %s", err)
	}
	checksum := getResourceChecksum(resource)
	log.Debugf("calculated checksum for stack as: %s", checksum)

	// @check if the resource has changed and if not we can return
	if found {
		status := stack.Status.Status
		// if the stack is found, check the status of the stack and if not finished we need to
	RETRY:
		switch status {
		case models.StatusDone:
			return stack, nil
		case models.StatusFailed:
			return stack, fmt.Errorf("stack failed on previous creation: %s", stack.Status.Reason)
		default:
			if status, err = c.options.Cloud.Wait(ctx, stackname, &models.WaitOptions{}); err != nil {
				return stack, fmt.Errorf("unable to wait on previous stack: %s", err)
			}
			log.WithFields(log.Fields{
				"namespace": resource.Namespace,
				"resource":  resource.Name,
			}).Info("rechecking the status of the stack")
			goto RETRY
		}

		// @check we have a checksum and check if its changed
		sum := stack.CheckSum()
		if sum == "" {
			return stack, fmt.Errorf("stack does not have a checksum, refusing to continue")
		}

		if sum == checksum {
			log.WithFields(log.Fields{
				"namespace": resource.Namespace,
				"resource":  resource.Name,
			}).Info("skipping updating the stack as nothing has changed")

			return stack, nil
		}
	}
	log.WithFields(log.Fields{
		"namespace": resource.Namespace,
		"resource":  resource.Name,
	}).Debug("checking the resource and template is valid")

	// @step: validate the cloud resource is ok
	if errs := resource.IsValid(); len(errs) > 0 {
		return stack, utils.GetErrors(errs)
	}

	// @check the template is valid and ok to us
	if errs := template.IsValid(); len(errs) > 0 {
		return stack, utils.GetErrors(errs)
	}

	// @step: we need build the parameters for the
	model, err := c.makeResourceModel(template, resource)
	if err != nil {
		return stack, err
	}

	log.WithFields(log.Fields{
		"model":     model,
		"namespace": resource.Namespace,
		"resource":  resource.Name,
		"stackname": stackname,
		"template":  template.Name,
	}).Info("attempting to create the stack")

	// @step: attempt to create the resource
	err = c.options.Cloud.Create(ctx, stackname, &models.CreateOptions{
		Context:  model,
		Resource: resource,
		Tags: map[string]string{
			models.CheckSumTag:     checksum,
			models.CreatedTag:      fmt.Sprintf("%d", time.Now().Unix()),
			models.NamespaceTag:    resource.Namespace,
			models.ProviderNameTag: c.config.Name,
			models.ResourceNameTag: resource.Name,
			models.RetentionTag:    fmt.Sprintf("%d", resource.Spec.Retention.Duration),
			models.TemplateNameTag: resource.Spec.TemplateName,
		},
		Template: template,
	})
	if err != nil {
		return stack, err
	}

	log.WithFields(log.Fields{
		"namespace": resource.Namespace,
		"resource":  resource.Name,
		"stackname": stackname,
	}).Info("successfully created the stack, waiting for the stack to complete")

	// @step: attempt to wait for the stack to finish
	status, err := c.options.Cloud.Wait(ctx, stackname, nil)
	if err != nil {
		return stack, err
	}

	log.WithFields(log.Fields{
		"namespace": resource.Namespace,
		"resource":  resource.Name,
		"stackname": stackname,
		"status":    status,
	}).Info("successfully wait on the stack completion")

	// @step: get the stacks
	stack, err = c.options.Cloud.Get(ctx, stackname, &models.GetOptions{})
	if err != nil {
		return stack, err
	}

	if status != models.StatusDone {
		return stack, errors.New("stack failed to complete successfully")
	}

	return stack, nil
}

// updateCloudCredentials is resposible for generating any credentials from the stack
// effectively is scas
func (c *controller) updateCloudCredentials(ctx context.Context, resource *apiv1.CloudResource, stack *models.Stack) (map[string]models.Credential, error) {
	users := make(map[string]models.Credential, 0)

	list, err := c.options.Cloud.Credentials(ctx, stack.Name)
	if err != nil {
		return users, err
	}

	for _, x := range list {
		users[x.ID] = x
	}

	return users, nil
}
