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

package models

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// IsValid checks the Stack is valid
func (s *Stack) IsValid() field.ErrorList {
	var errs field.ErrorList

	if s.Name == "" {
		errs = append(errs, field.Invalid(field.NewPath("name"), s.Name, "no name specified"))
	}
	if s.Namespace == "" {
		errs = append(errs, field.Invalid(field.NewPath("namespace"), s.Namespace, "no namespace specifieid"))
	}
	if s.Created.Unix() <= 0 {
		errs = append(errs, field.Invalid(field.NewPath("created"), "", "no creation specified"))
	}

	spec := field.NewPath("spec")
	if s.Spec.Name == "" {
		errs = append(errs, field.Invalid(spec.Key("name"), "", "no resource name specified"))
	}

	return errs
}
