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

	"github.com/gambol99/resources/pkg/models"
)

// Exists checks if a stack exists
func (p *provider) Exists(ctx context.Context, name string) (*models.Stack, bool, error) {
	stack, err := p.Get(ctx, name, &models.GetOptions{})
	if err != nil {
		if err == models.ErrStackNotFound {
			return nil, false, nil
		}

		return nil, false, err
	}

	return stack, true, nil
}
