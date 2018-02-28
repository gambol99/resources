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

// CloudStatusesGetter has a method to return a CloudStatusInterface.
// A group's client should implement this interface.
type CloudStatusesGetter interface {
	CloudStatuses(namespace string) CloudStatusInterface
}

// CloudStatusInterface has methods to work with CloudStatus resources.
type CloudStatusInterface interface {
	Create(*v1.CloudStatus) (*v1.CloudStatus, error)
	Update(*v1.CloudStatus) (*v1.CloudStatus, error)
	UpdateStatus(*v1.CloudStatus) (*v1.CloudStatus, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.CloudStatus, error)
	List(opts meta_v1.ListOptions) (*v1.CloudStatusList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.CloudStatus, err error)
	CloudStatusExpansion
}

// cloudStatuses implements CloudStatusInterface
type cloudStatuses struct {
	client rest.Interface
	ns     string
}

// newCloudStatuses returns a CloudStatuses
func newCloudStatuses(c *CloudV1Client, namespace string) *cloudStatuses {
	return &cloudStatuses{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the cloudStatus, and returns the corresponding cloudStatus object, and an error if there is any.
func (c *cloudStatuses) Get(name string, options meta_v1.GetOptions) (result *v1.CloudStatus, err error) {
	result = &v1.CloudStatus{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("cloudstatuses").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CloudStatuses that match those selectors.
func (c *cloudStatuses) List(opts meta_v1.ListOptions) (result *v1.CloudStatusList, err error) {
	result = &v1.CloudStatusList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("cloudstatuses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested cloudStatuses.
func (c *cloudStatuses) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("cloudstatuses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a cloudStatus and creates it.  Returns the server's representation of the cloudStatus, and an error, if there is any.
func (c *cloudStatuses) Create(cloudStatus *v1.CloudStatus) (result *v1.CloudStatus, err error) {
	result = &v1.CloudStatus{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("cloudstatuses").
		Body(cloudStatus).
		Do().
		Into(result)
	return
}

// Update takes the representation of a cloudStatus and updates it. Returns the server's representation of the cloudStatus, and an error, if there is any.
func (c *cloudStatuses) Update(cloudStatus *v1.CloudStatus) (result *v1.CloudStatus, err error) {
	result = &v1.CloudStatus{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("cloudstatuses").
		Name(cloudStatus.Name).
		Body(cloudStatus).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *cloudStatuses) UpdateStatus(cloudStatus *v1.CloudStatus) (result *v1.CloudStatus, err error) {
	result = &v1.CloudStatus{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("cloudstatuses").
		Name(cloudStatus.Name).
		SubResource("status").
		Body(cloudStatus).
		Do().
		Into(result)
	return
}

// Delete takes name of the cloudStatus and deletes it. Returns an error if one occurs.
func (c *cloudStatuses) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("cloudstatuses").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *cloudStatuses) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("cloudstatuses").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched cloudStatus.
func (c *cloudStatuses) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.CloudStatus, err error) {
	result = &v1.CloudStatus{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("cloudstatuses").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
