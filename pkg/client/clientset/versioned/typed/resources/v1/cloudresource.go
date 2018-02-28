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

// CloudResourcesGetter has a method to return a CloudResourceInterface.
// A group's client should implement this interface.
type CloudResourcesGetter interface {
	CloudResources(namespace string) CloudResourceInterface
}

// CloudResourceInterface has methods to work with CloudResource resources.
type CloudResourceInterface interface {
	Create(*v1.CloudResource) (*v1.CloudResource, error)
	Update(*v1.CloudResource) (*v1.CloudResource, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.CloudResource, error)
	List(opts meta_v1.ListOptions) (*v1.CloudResourceList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.CloudResource, err error)
	CloudResourceExpansion
}

// cloudResources implements CloudResourceInterface
type cloudResources struct {
	client rest.Interface
	ns     string
}

// newCloudResources returns a CloudResources
func newCloudResources(c *CloudV1Client, namespace string) *cloudResources {
	return &cloudResources{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the cloudResource, and returns the corresponding cloudResource object, and an error if there is any.
func (c *cloudResources) Get(name string, options meta_v1.GetOptions) (result *v1.CloudResource, err error) {
	result = &v1.CloudResource{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("cloudresources").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CloudResources that match those selectors.
func (c *cloudResources) List(opts meta_v1.ListOptions) (result *v1.CloudResourceList, err error) {
	result = &v1.CloudResourceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("cloudresources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested cloudResources.
func (c *cloudResources) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("cloudresources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a cloudResource and creates it.  Returns the server's representation of the cloudResource, and an error, if there is any.
func (c *cloudResources) Create(cloudResource *v1.CloudResource) (result *v1.CloudResource, err error) {
	result = &v1.CloudResource{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("cloudresources").
		Body(cloudResource).
		Do().
		Into(result)
	return
}

// Update takes the representation of a cloudResource and updates it. Returns the server's representation of the cloudResource, and an error, if there is any.
func (c *cloudResources) Update(cloudResource *v1.CloudResource) (result *v1.CloudResource, err error) {
	result = &v1.CloudResource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("cloudresources").
		Name(cloudResource.Name).
		Body(cloudResource).
		Do().
		Into(result)
	return
}

// Delete takes name of the cloudResource and deletes it. Returns an error if one occurs.
func (c *cloudResources) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("cloudresources").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *cloudResources) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("cloudresources").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched cloudResource.
func (c *cloudResources) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.CloudResource, err error) {
	result = &v1.CloudResource{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("cloudresources").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
