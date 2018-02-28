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

import "errors"

// IsValid checks the options are valid
func (c *CreateOptions) IsValid() error {
	if c.Resource == nil {
		return errors.New("no resource specified")
	}
	if c.Template == nil {
		return errors.New("no template specified")
	}
	if c.Context == nil {
		return errors.New("no values specified")
	}

	return nil
}

// IsOutput checks if the output exists
func (s *Stack) IsOutput(name string) bool {
	_, found := s.Spec.Outputs[name]

	return found
}

// Output returns the value of the output
func (s *Stack) Output(name string) string {
	v, _ := s.Spec.Outputs[name]

	return v
}
