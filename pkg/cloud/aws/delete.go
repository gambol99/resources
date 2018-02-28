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

package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/gambol99/resources/pkg/models"
)

// Delete is responsible for removing the stack
func (p *provider) Delete(ctx context.Context, name string, _ *models.DeleteOptions) error {
	// @step: we check the stack exists
	stack, _, err := p.getStack(ctx, name)
	if err != nil {
		if err == models.ErrStackNotFound {
			return nil
		}

		return err
	}

	// @check the ownership of the stack
	if !p.isOwned(stack) {
		return models.ErrUnauthorized
	}

	// @step: used to record the delete time
	metric := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		requestDuration.WithLabelValues("delete").Observe(v)
	}))
	defer metric.ObserveDuration()

	// @step: kick off the deletion of the stack
	_, err = p.client.DeleteStack(&cloudformation.DeleteStackInput{StackName: aws.String(name)})

	return err
}
