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
	"fmt"
	"strings"

	apiv1 "github.com/gambol99/resources/pkg/apis/resources/v1"
	"github.com/gambol99/resources/pkg/models"
	"github.com/gambol99/resources/pkg/utils"
)

// updateCloudSecrets is responsible for injecting the secrets into the namespace
func (c *controller) updateCloudSecrets(ctx context.Context, resource *apiv1.CloudResource, stack *models.Stack, creds map[string]models.Credential) error {
	for _, x := range resource.Spec.Secrets {
		values := make(map[string]string, 0)
		for _, k := range x.Values {
			switch k.Type {
			case apiv1.SecretTypeOutput:
				values[k.Key] = stack.Output(k.Value)
			case apiv1.SecretTypeCredential:
				items := strings.Split(k.Value, ".")
				if len(items) != 2 {
					return fmt.Errorf("invalid credential value: %s, should username.attribute for secret: %s", k.Value, x.Name)
				}
				user, found := creds[items[0]]
				if !found {
					return fmt.Errorf("credentials not found for secret: %s, reference: %s", x.Name, k.Value)
				}
				switch items[1] {
				case "username":
					values[k.Key] = user.User
				case "secret":
					values[k.Key] = user.Secret
				}
			}
		}
		// @step: inject the secret into the user namespace
		if err := utils.UpdateKubernetesSecret(c.options.Client, x.Name, resource.Namespace, values); err != nil {
			return err
		}
	}

	return nil
}
