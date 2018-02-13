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

package utils

import (
	"errors"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

// GetErrors returns the list of errors in a single error
func GetErrors(list field.ErrorList) error {
	if len(list) <= 0 {
		return nil
	}

	var reasons []string
	for _, x := range list {
		reasons = append(reasons, fmt.Sprintf("%s=%v : %s", x.Field, x.BadValue, x.Detail))
	}

	return errors.New(strings.Join(reasons, ","))
}
