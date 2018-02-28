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
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/gambol99/resources/pkg/models"
)

// Logs retrieves the logs from the cloudformation stack
func (p *provider) Logs(ctx context.Context, name string, options *models.GetOptions) (string, error) {
	_, found, err := p.Exists(ctx, name)
	if err != nil {
		return "", err
	}
	if !found {
		return "", models.ErrStackNotFound
	}

	log.WithFields(log.Fields{
		"stackname": name,
	}).Debug("retrieving the cloudformation logs for stack")

	tm := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		requestDuration.WithLabelValues("logs").Observe(v)
	}))
	resp, err := p.client.DescribeStackEventsWithContext(ctx, &cloudformation.DescribeStackEventsInput{
		StackName: aws.String(name),
	})
	if err != nil {
		return "", err
	}
	tm.ObserveDuration()

	var b bytes.Buffer
	b.WriteString("\n")
	for _, x := range resp.StackEvents {
		b.WriteString(fmt.Sprintf("%s %s %s %s\n",
			aws.StringValue(x.ResourceStatus),
			aws.StringValue(x.ResourceType),
			aws.StringValue(x.LogicalResourceId),
			aws.StringValue(x.ResourceStatusReason)))
	}

	return b.String(), nil
}
