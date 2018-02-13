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

package v1

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// HasParameter checks the parameter has been set
func (c *CloudResource) HasParameter(name string) bool {
	for _, x := range c.Spec.Parameters {
		if x.Name == name {
			return true
		}
	}

	return false
}

// AddSecret adds the secret to the resource if it doesn't exists already
func (c *CloudResource) AddSecret(secret Secret) {
	if c.HasSecret(secret.Name) {
		return
	}

	c.Spec.Secrets = append(c.Spec.Secrets, secret)
}

// HasSecret checks if a secret exists
func (c *CloudResource) HasSecret(name string) bool {
	for _, x := range c.Spec.Secrets {
		if name == x.Name {
			return true
		}
	}

	return false
}

// IsValid checks the secret is valid
func (s *Secret) IsValid(path *field.Path) field.ErrorList {
	var errs field.ErrorList

	if s.Name == "" {
		errs = append(errs, field.Invalid(path.Key("name"), "", "no name defined"))
	}
	if len(s.Values) <= 0 {
		errs = append(errs, field.Invalid(path.Key("values"), "", "no values defined"))
	}

	for i, x := range s.Values {
		errs = append(errs, x.IsValid(path.Index(i))...)
	}

	return errs
}

// IsValid checks the secret value is valid
func (s *SecretValue) IsValid(path *field.Path) field.ErrorList {
	var errs field.ErrorList
	if s.Type == "" {
		errs = append(errs, field.Invalid(path.Key("type"), s.Type, "no type defined"))
	}
	if s.Type != SecretTypeOutput && s.Type != SecretTypeCredential {
		errs = append(errs, field.Invalid(path.Key("type"), s.Type, "supported secret type"))
	}
	if s.Key == "" {
		errs = append(errs, field.Invalid(path.Key("key"), s.Key, "no key defined"))
	}
	if s.Value == "" {
		errs = append(errs, field.Invalid(path.Key("value"), s.Key, "no value defined"))
	}

	return errs
}

// IsValid checks the parameter is valid
func (p *Parameter) IsValid(path *field.Path, allowEmpty bool) field.ErrorList {
	var errs field.ErrorList

	if p.Name == "" {
		errs = append(errs, field.Invalid(path.Key("name"), p.Name, "no name given"))
	}
	if p.Value == nil && p.SecretName == nil && !allowEmpty {
		errs = append(errs, field.Invalid(path, "", "neither parameter value or secret reference set"))
	}

	return errs
}

// IsValid checks the cloud resource is valid
func (c *CloudResource) IsValid() field.ErrorList {
	var errs field.ErrorList

	if c.Spec.TemplateName == "" {
		errs = append(errs, field.Invalid(field.NewPath("spec").Key("templateName"), c.Spec.TemplateName, "no template name defined"))
	}
	for i, x := range c.Spec.Parameters {
		errs = append(errs, x.IsValid(field.NewPath("spec").Key("parameters").Index(i), false)...)
	}
	for i, x := range c.Spec.Secrets {
		errs = append(errs, x.IsValid(field.NewPath("spec").Key("secrets").Index(i))...)
	}

	return errs
}

// IsValid checks the template is valid
func (c *CloudTemplate) IsValid() field.ErrorList {
	var errs field.ErrorList

	spec := field.NewPath("spec")

	if c.Spec.Content == "" {
		errs = append(errs, field.Invalid(spec.Key("content"), c.Spec.Content, "no stack template specified"))
	}
	if c.Spec.Retention == nil {
		errs = append(errs, field.Invalid(spec.Key("retention"), c.Spec.Content, "no retention policy defined"))
	}
	if c.Spec.Format == "" {
		errs = append(errs, field.Invalid(spec.Key("format"), c.Spec.Format, "no format defined"))
	}
	if c.Spec.Format != FormatJSON && c.Spec.Format != FormatYAML {
		errs = append(errs, field.Invalid(spec.Key("format"), c.Spec.Format, "unsupported format"))
	}
	for i, x := range c.Spec.Parameters {
		errs = append(errs, x.IsValid(spec.Key("parameters").Index(i), true)...)
	}
	for i, x := range c.Spec.Secrets {
		errs = append(errs, x.IsValid(spec.Key("secrets").Index(i))...)
	}

	return errs
}
