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

// This file was automatically generated by lister-gen

package v1

import (
	v1 "github.com/gambol99/resources/pkg/apis/resources/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CloudResourceLister helps list CloudResources.
type CloudResourceLister interface {
	// List lists all CloudResources in the indexer.
	List(selector labels.Selector) (ret []*v1.CloudResource, err error)
	// CloudResources returns an object that can list and get CloudResources.
	CloudResources(namespace string) CloudResourceNamespaceLister
	CloudResourceListerExpansion
}

// cloudResourceLister implements the CloudResourceLister interface.
type cloudResourceLister struct {
	indexer cache.Indexer
}

// NewCloudResourceLister returns a new CloudResourceLister.
func NewCloudResourceLister(indexer cache.Indexer) CloudResourceLister {
	return &cloudResourceLister{indexer: indexer}
}

// List lists all CloudResources in the indexer.
func (s *cloudResourceLister) List(selector labels.Selector) (ret []*v1.CloudResource, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.CloudResource))
	})
	return ret, err
}

// CloudResources returns an object that can list and get CloudResources.
func (s *cloudResourceLister) CloudResources(namespace string) CloudResourceNamespaceLister {
	return cloudResourceNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CloudResourceNamespaceLister helps list and get CloudResources.
type CloudResourceNamespaceLister interface {
	// List lists all CloudResources in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.CloudResource, err error)
	// Get retrieves the CloudResource from the indexer for a given namespace and name.
	Get(name string) (*v1.CloudResource, error)
	CloudResourceNamespaceListerExpansion
}

// cloudResourceNamespaceLister implements the CloudResourceNamespaceLister
// interface.
type cloudResourceNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all CloudResources in the indexer for a given namespace.
func (s cloudResourceNamespaceLister) List(selector labels.Selector) (ret []*v1.CloudResource, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.CloudResource))
	})
	return ret, err
}

// Get retrieves the CloudResource from the indexer for a given namespace and name.
func (s cloudResourceNamespaceLister) Get(name string) (*v1.CloudResource, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("cloudresource"), name)
	}
	return obj.(*v1.CloudResource), nil
}
