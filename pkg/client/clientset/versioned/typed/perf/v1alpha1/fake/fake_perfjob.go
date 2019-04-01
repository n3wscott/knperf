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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/n3wscott/knperf/pkg/apis/perf/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakePerfJobs implements PerfJobInterface
type FakePerfJobs struct {
	Fake *FakePerfV1alpha1
	ns   string
}

var perfjobsResource = schema.GroupVersionResource{Group: "perf.knative.dev", Version: "v1alpha1", Resource: "perfjobs"}

var perfjobsKind = schema.GroupVersionKind{Group: "perf.knative.dev", Version: "v1alpha1", Kind: "PerfJob"}

// Get takes name of the perfJob, and returns the corresponding perfJob object, and an error if there is any.
func (c *FakePerfJobs) Get(name string, options v1.GetOptions) (result *v1alpha1.PerfJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(perfjobsResource, c.ns, name), &v1alpha1.PerfJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PerfJob), err
}

// List takes label and field selectors, and returns the list of PerfJobs that match those selectors.
func (c *FakePerfJobs) List(opts v1.ListOptions) (result *v1alpha1.PerfJobList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(perfjobsResource, perfjobsKind, c.ns, opts), &v1alpha1.PerfJobList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.PerfJobList{ListMeta: obj.(*v1alpha1.PerfJobList).ListMeta}
	for _, item := range obj.(*v1alpha1.PerfJobList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested perfJobs.
func (c *FakePerfJobs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(perfjobsResource, c.ns, opts))

}

// Create takes the representation of a perfJob and creates it.  Returns the server's representation of the perfJob, and an error, if there is any.
func (c *FakePerfJobs) Create(perfJob *v1alpha1.PerfJob) (result *v1alpha1.PerfJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(perfjobsResource, c.ns, perfJob), &v1alpha1.PerfJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PerfJob), err
}

// Update takes the representation of a perfJob and updates it. Returns the server's representation of the perfJob, and an error, if there is any.
func (c *FakePerfJobs) Update(perfJob *v1alpha1.PerfJob) (result *v1alpha1.PerfJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(perfjobsResource, c.ns, perfJob), &v1alpha1.PerfJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PerfJob), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakePerfJobs) UpdateStatus(perfJob *v1alpha1.PerfJob) (*v1alpha1.PerfJob, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(perfjobsResource, "status", c.ns, perfJob), &v1alpha1.PerfJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PerfJob), err
}

// Delete takes name of the perfJob and deletes it. Returns an error if one occurs.
func (c *FakePerfJobs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(perfjobsResource, c.ns, name), &v1alpha1.PerfJob{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePerfJobs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(perfjobsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.PerfJobList{})
	return err
}

// Patch applies the patch and returns the patched perfJob.
func (c *FakePerfJobs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PerfJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(perfjobsResource, c.ns, name, data, subresources...), &v1alpha1.PerfJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PerfJob), err
}
