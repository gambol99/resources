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

// FakeCloudTemplates implements CloudTemplateInterface
type FakeCloudTemplates struct {
	Fake *FakeCloudV1
}

var cloudtemplatesResource = schema.GroupVersionResource{Group: "cloud.appvia.io", Version: "v1", Resource: "cloudtemplates"}

var cloudtemplatesKind = schema.GroupVersionKind{Group: "cloud.appvia.io", Version: "v1", Kind: "CloudTemplate"}

// Get takes name of the cloudTemplate, and returns the corresponding cloudTemplate object, and an error if there is any.
func (c *FakeCloudTemplates) Get(name string, options v1.GetOptions) (result *resources_v1.CloudTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(cloudtemplatesResource, name), &resources_v1.CloudTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*resources_v1.CloudTemplate), err
}

// List takes label and field selectors, and returns the list of CloudTemplates that match those selectors.
func (c *FakeCloudTemplates) List(opts v1.ListOptions) (result *resources_v1.CloudTemplateList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(cloudtemplatesResource, cloudtemplatesKind, opts), &resources_v1.CloudTemplateList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &resources_v1.CloudTemplateList{}
	for _, item := range obj.(*resources_v1.CloudTemplateList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested cloudTemplates.
func (c *FakeCloudTemplates) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(cloudtemplatesResource, opts))
}

// Create takes the representation of a cloudTemplate and creates it.  Returns the server's representation of the cloudTemplate, and an error, if there is any.
func (c *FakeCloudTemplates) Create(cloudTemplate *resources_v1.CloudTemplate) (result *resources_v1.CloudTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(cloudtemplatesResource, cloudTemplate), &resources_v1.CloudTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*resources_v1.CloudTemplate), err
}

// Update takes the representation of a cloudTemplate and updates it. Returns the server's representation of the cloudTemplate, and an error, if there is any.
func (c *FakeCloudTemplates) Update(cloudTemplate *resources_v1.CloudTemplate) (result *resources_v1.CloudTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(cloudtemplatesResource, cloudTemplate), &resources_v1.CloudTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*resources_v1.CloudTemplate), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeCloudTemplates) UpdateStatus(cloudTemplate *resources_v1.CloudTemplate) (*resources_v1.CloudTemplate, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(cloudtemplatesResource, "status", cloudTemplate), &resources_v1.CloudTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*resources_v1.CloudTemplate), err
}

// Delete takes name of the cloudTemplate and deletes it. Returns an error if one occurs.
func (c *FakeCloudTemplates) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(cloudtemplatesResource, name), &resources_v1.CloudTemplate{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCloudTemplates) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(cloudtemplatesResource, listOptions)

	_, err := c.Fake.Invokes(action, &resources_v1.CloudTemplateList{})
	return err
}

// Patch applies the patch and returns the patched cloudTemplate.
func (c *FakeCloudTemplates) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *resources_v1.CloudTemplate, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(cloudtemplatesResource, name, data, subresources...), &resources_v1.CloudTemplate{})
	if obj == nil {
		return nil, err
	}
	return obj.(*resources_v1.CloudTemplate), err
}
