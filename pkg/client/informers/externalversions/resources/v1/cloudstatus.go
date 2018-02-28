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

// This file was automatically generated by informer-gen

package v1

import (
	time "time"

	resources_v1 "github.com/gambol99/resources/pkg/apis/resources/v1"
	versioned "github.com/gambol99/resources/pkg/client/clientset/versioned"
	internalinterfaces "github.com/gambol99/resources/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/gambol99/resources/pkg/client/listers/resources/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// CloudStatusInformer provides access to a shared informer and lister for
// CloudStatuses.
type CloudStatusInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.CloudStatusLister
}

type cloudStatusInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewCloudStatusInformer constructs a new informer for CloudStatus type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewCloudStatusInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredCloudStatusInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredCloudStatusInformer constructs a new informer for CloudStatus type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredCloudStatusInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CloudV1().CloudStatuses(namespace).List(options)
			},
			WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CloudV1().CloudStatuses(namespace).Watch(options)
			},
		},
		&resources_v1.CloudStatus{},
		resyncPeriod,
		indexers,
	)
}

func (f *cloudStatusInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredCloudStatusInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *cloudStatusInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&resources_v1.CloudStatus{}, f.defaultInformer)
}

func (f *cloudStatusInformer) Lister() v1.CloudStatusLister {
	return v1.NewCloudStatusLister(f.Informer().GetIndexer())
}