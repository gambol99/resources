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

// Object is generic runtime object used in the templating
type Object struct {
	// ID is the id for the resource
	ID string
	// Name is the name of the resource
	Name string
	// Tags is a collection of tags for the resources
	Tags map[string]string
}

// Network is a network interface
type Network struct {
	Object
	// CIDR is the network cidr
	CIDR string
	// AvailabilityZone is the availability zone
	AvailabilityZone string
}
