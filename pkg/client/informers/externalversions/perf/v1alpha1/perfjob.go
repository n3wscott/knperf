/*
Copyright 2019 The Knative Authors

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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	perfv1alpha1 "github.com/n3wscott/knperf/pkg/apis/perf/v1alpha1"
	versioned "github.com/n3wscott/knperf/pkg/client/clientset/versioned"
	internalinterfaces "github.com/n3wscott/knperf/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/n3wscott/knperf/pkg/client/listers/perf/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// PerfJobInformer provides access to a shared informer and lister for
// PerfJobs.
type PerfJobInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.PerfJobLister
}

type perfJobInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewPerfJobInformer constructs a new informer for PerfJob type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewPerfJobInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredPerfJobInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredPerfJobInformer constructs a new informer for PerfJob type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredPerfJobInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.PerfV1alpha1().PerfJobs(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.PerfV1alpha1().PerfJobs(namespace).Watch(options)
			},
		},
		&perfv1alpha1.PerfJob{},
		resyncPeriod,
		indexers,
	)
}

func (f *perfJobInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredPerfJobInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *perfJobInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&perfv1alpha1.PerfJob{}, f.defaultInformer)
}

func (f *perfJobInformer) Lister() v1alpha1.PerfJobLister {
	return v1alpha1.NewPerfJobLister(f.Informer().GetIndexer())
}
