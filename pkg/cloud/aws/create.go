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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/ghodss/yaml"

	"github.com/gambol99/resources/pkg/models"
)

// Create is responsible for creating or updating a stack
func (p *provider) Create(ctx context.Context, name string, options *models.CreateOptions) error {
	template := options.Template

	// check the options are valid
	if err := options.IsValid(); err != nil {
		return err
	}
	if name == "" {
		return errors.New("no name specified for stack")
	}

	// @step: parse and generate the template
	generated, err := NewTemplater(p.compute, p.config).Render(ctx, options.Context, template.Spec.Content)
	if err != nil {
		return err
	}

	// @step: attempt to validate the stack before sending it, we don't want to waste time
	if _, err = p.client.ValidateTemplateWithContext(ctx, &cloudformation.ValidateTemplateInput{
		TemplateBody: aws.String(generated),
	}); err != nil {
		return err
	}

	// @step: check if the resource already exists and if so is in-progress
	isTrue := true
	found, err := p.hasStack(ctx, name)
	if err != nil {
		return err
	}

	// @step: is the format of the template is YAML, convert the template to JSON before sending
	if template.Spec.Format == "yaml" {
		encoded, err := yaml.YAMLToJSON([]byte(generated))
		if err != nil {
			return fmt.Errorf("unable to convert yaml to json format: %s", err)
		}
		generated = string(encoded)
	}

	if !found {
		// we are creating a new stack
		if _, err := p.client.CreateStack(&cloudformation.CreateStackInput{
			Capabilities:                aws.StringSlice([]string{"CAPABILITY_IAM"}),
			DisableRollback:             aws.Bool(isTrue),
			EnableTerminationProtection: aws.Bool(isTrue),
			StackName:                   aws.String(name),
			Tags:                        makeStackTags(options.Tags),
			TemplateBody:                aws.String(generated),
		}); err != nil {
			return err
		}
	} else {
		// @step: we are updating a cloudformation stack
		if _, err := p.client.UpdateStack(&cloudformation.UpdateStackInput{
			StackName:    aws.String(name),
			Tags:         makeStackTags(options.Tags),
			TemplateBody: aws.String(generated),
		}); err != nil {
			return err
		}
	}

	return nil
}
