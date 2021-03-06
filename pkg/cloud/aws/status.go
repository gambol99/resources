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

	"github.com/gambol99/resources/pkg/models"
)

// Status is responsible for getting the status
func (p *provider) Status(cx context.Context, name string, _ *models.GetOptions) (string, error) {
	stack, _, err := p.getStack(cx, name)
	if err != nil {
		return models.StatusUnknown, err
	}

	return getStackStatus(aws.StringValue(stack.StackStatus)), nil
}
