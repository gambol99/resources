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
	"crypto/md5"
	"fmt"
	"io"
	"reflect"

	apiv1 "github.com/gambol99/resources/pkg/apis/resources/v1"
	"github.com/gambol99/resources/pkg/utils"
)

// makeResourceModel is resposible for consolidating the parameteres, secrets and attributes
func (c *controller) makeResourceModel(template *apiv1.CloudTemplate, resource *apiv1.CloudResource) (map[string]string, error) {
	values := make(map[string]string, 0)

	{
		// @step: in the deletion policy and retention
		if resource.Spec.Retention == nil {
			resource.Spec.Retention = template.Spec.Retention
		}
	}

	{
		// @step: copy the default parameters from the template into the model
		for _, x := range template.Spec.Parameters {
			if x.Value != nil {
				values[x.Name] = *x.Value
			}
			// @check if no default is set and the parameter is required that parameter is set
			if !resource.HasParameter(x.Name) {
				return values, fmt.Errorf("resource parameter: '%s' is required", x.Name)
			}
		}
		// @step: we need to iterate the resource parameters and pull in the values or the kubernetes secrets
		for _, x := range resource.Spec.Parameters {
			// @check if this a static parameter
			if x.Value != nil {
				values[x.Name] = *x.Value
				continue
			}
			if x.SecretName != nil {
				// @step: pull the kubernetes secret from the resource's namespace
				secret, err := utils.FindKubernetesSecret(c.options.Client, *x.SecretName, resource.Namespace)
				if err != nil {
					return values, fmt.Errorf("paramater: '%s' unable to pull from kubernetes secret: %s", x.Name, err)
				}
				switch len(secret) {
				case 0:
					return values, fmt.Errorf("parameter: '%s' kubernetes secret has no value", x.Name)
				case 1:
					keys := reflect.ValueOf(secret).MapKeys()
					values[x.Name] = fmt.Sprintf("%s", keys[0])
				default:
					return values, fmt.Errorf("parameter: '%s' kubernetes secret has multiple keys", x.Name)
				}
				continue
			}
			// @step: thrown an error and nothing has been set
			return values, fmt.Errorf("resource parameter: '%s' has no value or kubernetes secret set", x.Name)
		}
	}

	{
		// @step: inject the secrets from the template
		for _, x := range template.Spec.Secrets {
			resource.AddSecret(x)
		}
	}

	return values, nil
}

// getResourceChecksum is responsible for checking if the resource parameters have changed
func getResourceChecksum(resource *apiv1.CloudResource) string {
	h := md5.New()
	for _, x := range resource.Spec.Parameters {
		io.WriteString(h, x.Name)
		if x.SecretName != nil {
			io.WriteString(h, *x.SecretName)
		}
		if x.Value != nil {
			io.WriteString(h, *x.Value)
		}
	}

	return string(h.Sum(nil))
}

// getStackName is the default naming convertion for all formation stacks
func getStackName(name, namespace string) string {
	return fmt.Sprintf("stacks_%s_%s", namespace, name)
}
