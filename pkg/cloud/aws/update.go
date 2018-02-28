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
)

// UpdateTags is responsible for updating just the tags of a stack
func (p *provider) UpdateTags(ctx context.Context, name string, tags map[string]string) error {

	metric := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		requestDuration.WithLabelValues("tags").Observe(v)
	}))
	defer metric.ObserveDuration()

	_, err := p.client.UpdateStackWithContext(ctx, &cloudformation.UpdateStackInput{
		StackName: aws.String(name),
		Tags:      makeStackTags(tags),
	})

	return err
}
