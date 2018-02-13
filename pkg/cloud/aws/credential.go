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

// Credentials is responsible for creating a credential in the cloud
func (p *provider) Credentials(ctx context.Context, name string) ([]models.Credential, error) {
	var list []models.Credential

	// @step: we get a list of users from the stack
	users, err := p.findIAMUsers(ctx, name)
	if err != nil {
		return list, err
	}
	if len(users) < 0 {
		return list, nil
	}

	// @step: we attempt to create a credential for the user
	for _, x := range users {
		access, key, err := p.getAccessToken(ctx, x)
		if err != nil {
			return list, err
		}
		list = append(list, models.Credential{ID: x, User: access, Secret: key})
	}

	return list, nil
}
