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
package v1

import (
	v1 "github.com/gambol99/resources/pkg/apis/resources/v1"
	scheme "github.com/gambol99/resources/pkg/client/clientset/versioned/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// CloudTemplatesGetter has a method to return a CloudTemplateInterface.
// A group's client should implement this interface.
type CloudTemplatesGetter interface {
	CloudTemplates() CloudTemplateInterface
}

// CloudTemplateInterface has methods to work with CloudTemplate resources.
type CloudTemplateInterface interface {
	Create(*v1.CloudTemplate) (*v1.CloudTemplate, error)
	Update(*v1.CloudTemplate) (*v1.CloudTemplate, error)
	UpdateStatus(*v1.CloudTemplate) (*v1.CloudTemplate, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.CloudTemplate, error)
	List(opts meta_v1.ListOptions) (*v1.CloudTemplateList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.CloudTemplate, err error)
	CloudTemplateExpansion
}

// cloudTemplates implements CloudTemplateInterface
type cloudTemplates struct {
	client rest.Interface
}

// newCloudTemplates returns a CloudTemplates
func newCloudTemplates(c *CloudV1Client) *cloudTemplates {
	return &cloudTemplates{
		client: c.RESTClient(),
	}
}

// Get takes name of the cloudTemplate, and returns the corresponding cloudTemplate object, and an error if there is any.
func (c *cloudTemplates) Get(name string, options meta_v1.GetOptions) (result *v1.CloudTemplate, err error) {
	result = &v1.CloudTemplate{}
	err = c.client.Get().
		Resource("cloudtemplates").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CloudTemplates that match those selectors.
func (c *cloudTemplates) List(opts meta_v1.ListOptions) (result *v1.CloudTemplateList, err error) {
	result = &v1.CloudTemplateList{}
	err = c.client.Get().
		Resource("cloudtemplates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested cloudTemplates.
func (c *cloudTemplates) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Resource("cloudtemplates").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a cloudTemplate and creates it.  Returns the server's representation of the cloudTemplate, and an error, if there is any.
func (c *cloudTemplates) Create(cloudTemplate *v1.CloudTemplate) (result *v1.CloudTemplate, err error) {
	result = &v1.CloudTemplate{}
	err = c.client.Post().
		Resource("cloudtemplates").
		Body(cloudTemplate).
		Do().
		Into(result)
	return
}

// Update takes the representation of a cloudTemplate and updates it. Returns the server's representation of the cloudTemplate, and an error, if there is any.
func (c *cloudTemplates) Update(cloudTemplate *v1.CloudTemplate) (result *v1.CloudTemplate, err error) {
	result = &v1.CloudTemplate{}
	err = c.client.Put().
		Resource("cloudtemplates").
		Name(cloudTemplate.Name).
		Body(cloudTemplate).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *cloudTemplates) UpdateStatus(cloudTemplate *v1.CloudTemplate) (result *v1.CloudTemplate, err error) {
	result = &v1.CloudTemplate{}
	err = c.client.Put().
		Resource("cloudtemplates").
		Name(cloudTemplate.Name).
		SubResource("status").
		Body(cloudTemplate).
		Do().
		Into(result)
	return
}

// Delete takes name of the cloudTemplate and deletes it. Returns an error if one occurs.
func (c *cloudTemplates) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("cloudtemplates").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *cloudTemplates) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Resource("cloudtemplates").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched cloudTemplate.
func (c *cloudTemplates) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.CloudTemplate, err error) {
	result = &v1.CloudTemplate{}
	err = c.client.Patch(pt).
		Resource("cloudtemplates").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
