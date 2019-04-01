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

package v1alpha1

import (
	v1alpha1 "github.com/n3wscott/knperf/pkg/apis/perf/v1alpha1"
	scheme "github.com/n3wscott/knperf/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PerfJobsGetter has a method to return a PerfJobInterface.
// A group's client should implement this interface.
type PerfJobsGetter interface {
	PerfJobs(namespace string) PerfJobInterface
}

// PerfJobInterface has methods to work with PerfJob resources.
type PerfJobInterface interface {
	Create(*v1alpha1.PerfJob) (*v1alpha1.PerfJob, error)
	Update(*v1alpha1.PerfJob) (*v1alpha1.PerfJob, error)
	UpdateStatus(*v1alpha1.PerfJob) (*v1alpha1.PerfJob, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.PerfJob, error)
	List(opts v1.ListOptions) (*v1alpha1.PerfJobList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PerfJob, err error)
	PerfJobExpansion
}

// perfJobs implements PerfJobInterface
type perfJobs struct {
	client rest.Interface
	ns     string
}

// newPerfJobs returns a PerfJobs
func newPerfJobs(c *PerfV1alpha1Client, namespace string) *perfJobs {
	return &perfJobs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the perfJob, and returns the corresponding perfJob object, and an error if there is any.
func (c *perfJobs) Get(name string, options v1.GetOptions) (result *v1alpha1.PerfJob, err error) {
	result = &v1alpha1.PerfJob{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("perfjobs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of PerfJobs that match those selectors.
func (c *perfJobs) List(opts v1.ListOptions) (result *v1alpha1.PerfJobList, err error) {
	result = &v1alpha1.PerfJobList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("perfjobs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested perfJobs.
func (c *perfJobs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("perfjobs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a perfJob and creates it.  Returns the server's representation of the perfJob, and an error, if there is any.
func (c *perfJobs) Create(perfJob *v1alpha1.PerfJob) (result *v1alpha1.PerfJob, err error) {
	result = &v1alpha1.PerfJob{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("perfjobs").
		Body(perfJob).
		Do().
		Into(result)
	return
}

// Update takes the representation of a perfJob and updates it. Returns the server's representation of the perfJob, and an error, if there is any.
func (c *perfJobs) Update(perfJob *v1alpha1.PerfJob) (result *v1alpha1.PerfJob, err error) {
	result = &v1alpha1.PerfJob{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("perfjobs").
		Name(perfJob.Name).
		Body(perfJob).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *perfJobs) UpdateStatus(perfJob *v1alpha1.PerfJob) (result *v1alpha1.PerfJob, err error) {
	result = &v1alpha1.PerfJob{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("perfjobs").
		Name(perfJob.Name).
		SubResource("status").
		Body(perfJob).
		Do().
		Into(result)
	return
}

// Delete takes name of the perfJob and deletes it. Returns an error if one occurs.
func (c *perfJobs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("perfjobs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *perfJobs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("perfjobs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched perfJob.
func (c *perfJobs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PerfJob, err error) {
	result = &v1alpha1.PerfJob{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("perfjobs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
