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

// List is responsible for getting a list of stacks
func (p *provider) List(ctx context.Context, options *models.ListOptions) ([]*models.Stack, error) {
	var list []*models.Stack

	metric := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		requestDuration.WithLabelValues("list").Observe(v)
	}))

	resp, err := p.client.ListStacksWithContext(ctx, &cloudformation.ListStacksInput{})
	if err != nil {
		return list, err
	}
	metric.ObserveDuration()

	for _, x := range resp.StackSummaries {
		stack, content, err := p.getStack(ctx, aws.StringValue(x.StackName))
		if err != nil {
			return list, err
		}
		// @step: filter out the stack if its not owned by use
		if !p.isOwned(stack) {
			continue
		}
		s, err := makeStack(stack, content)
		if err != nil {
			return list, err
		}

		list = append(list, s)
	}

	return list, nil
}
