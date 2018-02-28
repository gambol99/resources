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

// FakeCloudResources implements CloudResourceInterface
type FakeCloudResources struct {
	Fake *FakeCloudV1
	ns   string
}

var cloudresourcesResource = schema.GroupVersionResource{Group: "cloud.appvia.io", Version: "v1", Resource: "cloudresources"}

var cloudresourcesKind = schema.GroupVersionKind{Group: "cloud.appvia.io", Version: "v1", Kind: "CloudResource"}

// Get takes name of the cloudResource, and returns the corresponding cloudResource object, and an error if there is any.
func (c *FakeCloudResources) Get(name string, options v1.GetOptions) (result *resources_v1.CloudResource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(cloudresourcesResource, c.ns, name), &resources_v1.CloudResource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*resources_v1.CloudResource), err
}

// List takes label and field selectors, and returns the list of CloudResources that match those selectors.
func (c *FakeCloudResources) List(opts v1.ListOptions) (result *resources_v1.CloudResourceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(cloudresourcesResource, cloudresourcesKind, c.ns, opts), &resources_v1.CloudResourceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &resources_v1.CloudResourceList{}
	for _, item := range obj.(*resources_v1.CloudResourceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested cloudResources.
func (c *FakeCloudResources) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(cloudresourcesResource, c.ns, opts))

}

// Create takes the representation of a cloudResource and creates it.  Returns the server's representation of the cloudResource, and an error, if there is any.
func (c *FakeCloudResources) Create(cloudResource *resources_v1.CloudResource) (result *resources_v1.CloudResource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(cloudresourcesResource, c.ns, cloudResource), &resources_v1.CloudResource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*resources_v1.CloudResource), err
}

// Update takes the representation of a cloudResource and updates it. Returns the server's representation of the cloudResource, and an error, if there is any.
func (c *FakeCloudResources) Update(cloudResource *resources_v1.CloudResource) (result *resources_v1.CloudResource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(cloudresourcesResource, c.ns, cloudResource), &resources_v1.CloudResource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*resources_v1.CloudResource), err
}

// Delete takes name of the cloudResource and deletes it. Returns an error if one occurs.
func (c *FakeCloudResources) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(cloudresourcesResource, c.ns, name), &resources_v1.CloudResource{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCloudResources) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(cloudresourcesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &resources_v1.CloudResourceList{})
	return err
}

// Patch applies the patch and returns the patched cloudResource.
func (c *FakeCloudResources) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *resources_v1.CloudResource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(cloudresourcesResource, c.ns, name, data, subresources...), &resources_v1.CloudResource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*resources_v1.CloudResource), err
}
