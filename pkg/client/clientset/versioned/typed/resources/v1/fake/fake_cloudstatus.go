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
package fake

import (
	resources_v1 "github.com/gambol99/resources/pkg/apis/resources/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeCloudStatuses implements CloudStatusInterface
type FakeCloudStatuses struct {
	Fake *FakeCloudV1
	ns   string
}

var cloudstatusesResource = schema.GroupVersionResource{Group: "cloud.appvia.io", Version: "v1", Resource: "cloudstatuses"}

var cloudstatusesKind = schema.GroupVersionKind{Group: "cloud.appvia.io", Version: "v1", Kind: "CloudStatus"}

// Get takes name of the cloudStatus, and returns the corresponding cloudStatus object, and an error if there is any.
func (c *FakeCloudStatuses) Get(name string, options v1.GetOptions) (result *resources_v1.CloudStatus, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(cloudstatusesResource, c.ns, name), &resources_v1.CloudStatus{})

	if obj == nil {
		return nil, err
	}
	return obj.(*resources_v1.CloudStatus), err
}

// List takes label and field selectors, and returns the list of CloudStatuses that match those selectors.
func (c *FakeCloudStatuses) List(opts v1.ListOptions) (result *resources_v1.CloudStatusList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(cloudstatusesResource, cloudstatusesKind, c.ns, opts), &resources_v1.CloudStatusList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &resources_v1.CloudStatusList{}
	for _, item := range obj.(*resources_v1.CloudStatusList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested cloudStatuses.
func (c *FakeCloudStatuses) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(cloudstatusesResource, c.ns, opts))

}

// Create takes the representation of a cloudStatus and creates it.  Returns the server's representation of the cloudStatus, and an error, if there is any.
func (c *FakeCloudStatuses) Create(cloudStatus *resources_v1.CloudStatus) (result *resources_v1.CloudStatus, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(cloudstatusesResource, c.ns, cloudStatus), &resources_v1.CloudStatus{})

	if obj == nil {
		return nil, err
	}
	return obj.(*resources_v1.CloudStatus), err
}

// Update takes the representation of a cloudStatus and updates it. Returns the server's representation of the cloudStatus, and an error, if there is any.
func (c *FakeCloudStatuses) Update(cloudStatus *resources_v1.CloudStatus) (result *resources_v1.CloudStatus, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(cloudstatusesResource, c.ns, cloudStatus), &resources_v1.CloudStatus{})

	if obj == nil {
		return nil, err
	}
	return obj.(*resources_v1.CloudStatus), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeCloudStatuses) UpdateStatus(cloudStatus *resources_v1.CloudStatus) (*resources_v1.CloudStatus, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(cloudstatusesResource, "status", c.ns, cloudStatus), &resources_v1.CloudStatus{})

	if obj == nil {
		return nil, err
	}
	return obj.(*resources_v1.CloudStatus), err
}

// Delete takes name of the cloudStatus and deletes it. Returns an error if one occurs.
func (c *FakeCloudStatuses) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(cloudstatusesResource, c.ns, name), &resources_v1.CloudStatus{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCloudStatuses) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(cloudstatusesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &resources_v1.CloudStatusList{})
	return err
}

// Patch applies the patch and returns the patched cloudStatus.
func (c *FakeCloudStatuses) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *resources_v1.CloudStatus, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(cloudstatusesResource, c.ns, name, data, subresources...), &resources_v1.CloudStatus{})

	if obj == nil {
		return nil, err
	}
	return obj.(*resources_v1.CloudStatus), err
}
