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
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/gambol99/resources/pkg/models"
	"github.com/gambol99/resources/pkg/utils"
)

// getStack is responsible for retrieving the stack
func (p *provider) getStack(ctx context.Context, name string) (*cloudformation.Stack, string, error) {
	tm := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		requestDuration.WithLabelValues("get").Observe(v)
	}))

	log.WithFields(log.Fields{
		"stackname": name,
	}).Debug("retrieving cloudformation stack")

	resp, err := p.client.DescribeStacksWithContext(ctx, &cloudformation.DescribeStacksInput{
		StackName: aws.String(name),
	})
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil, "", models.ErrStackNotFound
		}

		return nil, "", err
	}

	if len(resp.Stacks) <= 0 {
		return nil, "", models.ErrStackNotFound
	}

	content, err := p.getStackTemplate(ctx, name)
	if err != nil {
		return nil, "", err
	}

	defer tm.ObserveDuration()

	return resp.Stacks[0], content, nil
}

// getStackTemplate is responsible for retrieving the underlining template body of a stack
func (p *provider) getStackTemplate(ctx context.Context, name string) (string, error) {
	resp, err := p.client.GetTemplateWithContext(ctx, &cloudformation.GetTemplateInput{
		StackName: aws.String(name),
	})
	if err != nil {
		return "", err
	}

	if resp.TemplateBody == nil {
		return "", models.ErrStackNotFound
	}

	return aws.StringValue(resp.TemplateBody), nil
}

// hasStack checks to see if the cloudformation already exists
func (p *provider) hasStack(ctx context.Context, name string) (bool, error) {
	if _, _, err := p.getStack(ctx, name); err != nil {
		if err == models.ErrStackNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// isOwned is responsible for checking the stack is owned by us
func (p *provider) isOwned(stack *cloudformation.Stack) bool {
	if stack == nil {
		return false
	}
	if len(stack.Tags) <= 0 {
		return false
	}
	// @check for provider tags
	if !checkForTags(models.ProviderNameTag, p.config.Name, stack.Tags) {
		return false
	}

	return true
}

// getAccessToken is responsible for generating an access token for the user
func (p *provider) getAccessToken(ctx context.Context, username string) (string, string, error) {
	// @step: get the user has not gone over the limit
	resp, err := p.accounts.ListAccessKeysWithContext(ctx, &iam.ListAccessKeysInput{
		UserName: aws.String(username),
	})
	if err != nil {
		return "", "", err
	}
	// @check if we have reached the max number of keys
	if len(resp.AccessKeyMetadata) >= maxAccessKeys {
		return "", "", fmt.Errorf("user: %s has reached max number of access keys", username)
	}

	res, err := p.accounts.CreateAccessKeyWithContext(ctx, &iam.CreateAccessKeyInput{
		UserName: aws.String(username),
	})
	if err != nil {
		return "", "", err
	}
	if res.AccessKey == nil {
		return "", "", fmt.Errorf("no access returns in response for user: %s", username)
	}

	return aws.StringValue(res.AccessKey.AccessKeyId), aws.StringValue(res.AccessKey.SecretAccessKey), nil
}

// getPolicyArn is responsible for resoling a policy name to aws ARN
func (p *provider) getPolicyArn(ctx context.Context, name string) (string, error) {
	if strings.HasPrefix(name, "arn:") {
		return name, nil
	}
	// @step: try and resolve the arn
	resp, err := p.accounts.ListPoliciesWithContext(ctx, &iam.ListPoliciesInput{
		Scope: aws.String("LOCAL"),
	})
	if err != nil {
		return "", err
	}
	for _, x := range resp.Policies {
		if name == aws.StringValue(x.PolicyName) {
			return aws.StringValue(x.Arn), nil
		}
	}

	return "", fmt.Errorf("policy: %s not found", name)
}

// findIAMUsers is responsible for finding any IAM users in the stack
// we effectively get the cloudformation stack, we iterate the resources and then we
// extract any IAM::User from the template
func (p *provider) findIAMUsers(ctx context.Context, name string) ([]string, error) {
	var list []string

	resp, err := p.client.DescribeStackResourcesWithContext(ctx, &cloudformation.DescribeStackResourcesInput{
		StackName: aws.String(name),
	})
	if err != nil {
		return list, err
	}

	for _, x := range resp.StackResources {
		if aws.StringValue(x.ResourceType) == "AWS::IAM::User" {
			list = append(list, aws.StringValue(x.PhysicalResourceId))
		}
	}

	return list, nil
}

// getClusterFilters returns a series of filters for aws resources
func getClusterFilters(cluster string) []*ec2.Filter {
	return []*ec2.Filter{
		{
			Name:   aws.String("tag:KubernetesCluster"),
			Values: []*string{aws.String(cluster)},
		},
		{
			Name:   aws.String("tag:kubernetes.io/cluster/" + cluster),
			Values: []*string{aws.String("owned"), aws.String("shared")},
		},
	}
}

// makeStackTags is responsible for converting a collection resource tags to cloudformation tags
func makeStackTags(values map[string]string) []*cloudformation.Tag {
	tags := make([]*cloudformation.Tag, 0)
	for k, v := range values {
		tags = append(tags, &cloudformation.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	return tags
}

// makeStack converts an aws stack a api.Stack
func makeStack(stack *cloudformation.Stack, body string) (*models.Stack, error) {
	s := &models.Stack{
		Name:    aws.StringValue(stack.StackName),
		Created: aws.TimeValue(stack.CreationTime),
		Spec: models.StackSpec{
			Outputs:  make(map[string]string, 0),
			Tags:     make(map[string]string, 0),
			Template: body,
		},
		Status: models.StackStatus{
			Status: getStackStatus(aws.StringValue(stack.StackStatus)),
		},
	}

	// @step: copy the outputs from the stack
	for _, x := range stack.Outputs {
		s.Spec.Outputs[aws.StringValue(x.OutputKey)] = aws.StringValue(x.OutputValue)
	}

	for _, x := range stack.Tags {
		// @step: add all the tags into the resource tags
		s.Spec.Tags[aws.StringValue(x.Key)] = aws.StringValue(x.Value)

		// @step: filter out the elements
		switch aws.StringValue(x.Key) {
		case models.DeletionTimeTag:
			tm, err := strconv.ParseInt(aws.StringValue(x.Value), 10, 64)
			if err != nil {
				return nil, errors.New("cannot convert to deletion time")
			}
			s.Spec.DeleteOn = time.Unix(tm, 0)
		case models.NamespaceTag:
			s.Namespace = aws.StringValue(x.Value)
		case models.ResourceNameTag:
			s.Spec.Name = aws.StringValue(x.Value)
		case models.RetentionTag:
			tm, err := strconv.ParseInt(aws.StringValue(x.Value), 10, 64)
			if err != nil {
				return nil, errors.New("cannot convert to duration")
			}
			s.Spec.Retention = time.Duration(tm)
		case models.TemplateNameTag:
			s.Spec.Template = aws.StringValue(x.Value)
		}
	}

	if errs := s.IsValid(); len(errs) > 0 {
		return nil, utils.GetErrors(errs)
	}

	return s, nil
}

// checkForTags checks if the aws tag is included
func checkForTags(key, value string, tags []*cloudformation.Tag) bool {
	for _, x := range tags {
		if aws.StringValue(x.Key) == key && aws.StringValue(x.Value) == value {
			return true
		}
	}

	return false
}

// getStackStatus coverts the cloudformation status a stack status
func getStackStatus(status string) string {
	switch status {
	case "CREATE_COMPLETE":
		return models.StatusDone
	case "CREATE_IN_PROGRESS":
		return models.StatusInProgress
	case "CREATE_FAILED":
		return models.StatusFailed
	case "DELETE_COMPLETE":
		return models.StatusDone
	case "DELETE_FAILED":
		return models.StatusFailed
	case "DELETE_IN_PROGRESS":
		return models.StatusInProgress
	case "REVIEW_IN_PROGRESS":
		return models.StatusInProgress
	case "ROLLBACK_COMPLETE":
		return models.StatusDone
	case "ROLLBACK_FAILED":
		return models.StatusFailed
	case "ROLLBACK_IN_PROGRESS":
		return models.StatusInProgress
	case "UPDATE_COMPLETE":
		return models.StatusDone
	case "UPDATE_COMPLETE_CLEANUP_IN_PROGRESS":
		return models.StatusInProgress
	case "UPDATE_IN_PROGRESS":
		return models.StatusInProgress
	case "UPDATE_ROLLBACK_COMPLETE":
		return models.StatusDone
	case "UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS":
		return models.StatusInProgress
	case "UPDATE_ROLLBACK_FAILED":
		return models.StatusDone
	case "UPDATE_ROLLBACK_IN_PROGRESS":
		return models.StatusInProgress
	}

	return ""
}
