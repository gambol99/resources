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
	"time"
)

// HasDeleteTag check of the stack has a deletion tag
func (s *Stack) HasDeleteTag() bool {
	if _, found := s.Spec.Tags[DeletionTimeTag]; found {
		return true
	}

	return false
}

// CheckSum returns the checksum of the stack
func (s *Stack) CheckSum() string {
	return s.Spec.Tags[CheckSumTag]
}

// RequiresDeletion checks if the stack should be deleted
func (s *Stack) RequiresDeletion() bool {
	if !s.HasDeleteTag() {
		return false
	}

	return s.Spec.DeleteOn.Before(time.Now())
}

// ExpiresIn returns the stack expiration
func (s *Stack) ExpiresIn() string {
	return time.Now().Sub(s.Spec.DeleteOn).String()
}
