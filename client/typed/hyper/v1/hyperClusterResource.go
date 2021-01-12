/*
Copyright The Kubernetes Authors.
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

package v1

import (
	"context"
	"time"

	v1 "github.com/tmax-cloud/hypercloud-server/external/hyper/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	scheme "k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

// HyperClusterResourceGetter has a method to return a HyperClusterResourceInterface.
// A group's client should implement this interface.
type HyperClusterResourceGetter interface {
	HyperClusterResources(namespace string) HyperClusterResourceInterface
}

// HyperClusterResourceInterface has methods to work with hyperClusterResource resources.
type HyperClusterResourceInterface interface {
	Create(ctx context.Context, hyperClusterResource *v1.HyperClusterResource, opts metav1.CreateOptions) (*v1.HyperClusterResource, error)
	Update(ctx context.Context, hyperClusterResource *v1.HyperClusterResource, opts metav1.UpdateOptions) (*v1.HyperClusterResource, error)
	UpdateStatus(ctx context.Context, hyperClusterResource *v1.HyperClusterResource, opts metav1.UpdateOptions) (*v1.HyperClusterResource, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.HyperClusterResource, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.HyperClusterResourceList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.HyperClusterResource, err error)
	hyperClusterResourceExpansion
}

// hyperClusterResources implements HyperClusterResourceInterface
type hyperClusterResources struct {
	client rest.Interface
	ns     string
}

// newfluentBitConfigurations returns a hyperClusterResources
func newHyperClusterResources(c *HyperV1Client, namespace string) *hyperClusterResources {
	return &hyperClusterResources{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the hyperClusterResource, and returns the corresponding hyperClusterResource object, and an error if there is any.
func (c *hyperClusterResources) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.HyperClusterResource, err error) {
	result = &v1.HyperClusterResource{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("hyperClusterResources").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of hyperClusterResources that match those selectors.
func (c *hyperClusterResources) List(ctx context.Context, opts metav1.ListOptions) (result *v1.HyperClusterResourceList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.HyperClusterResourceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("hyperClusterResources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested hyperClusterResources.
func (c *hyperClusterResources) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("hyperClusterResources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a hyperClusterResource and creates it.  Returns the server's representation of the hyperClusterResource, and an error, if there is any.
func (c *hyperClusterResources) Create(ctx context.Context, hyperClusterResource *v1.HyperClusterResource, opts metav1.CreateOptions) (result *v1.HyperClusterResource, err error) {
	result = &v1.HyperClusterResource{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("hyperClusterResources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(hyperClusterResource).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a hyperClusterResource and updates it. Returns the server's representation of the hyperClusterResource, and an error, if there is any.
func (c *hyperClusterResources) Update(ctx context.Context, hyperClusterResource *v1.HyperClusterResource, opts metav1.UpdateOptions) (result *v1.HyperClusterResource, err error) {
	result = &v1.HyperClusterResource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("hyperClusterResources").
		Name(hyperClusterResource.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(hyperClusterResource).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *hyperClusterResources) UpdateStatus(ctx context.Context, hyperClusterResource *v1.HyperClusterResource, opts metav1.UpdateOptions) (result *v1.HyperClusterResource, err error) {
	result = &v1.HyperClusterResource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("hyperClusterResources").
		Name(hyperClusterResource.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(hyperClusterResource).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the hyperClusterResource and deletes it. Returns an error if one occurs.
func (c *hyperClusterResources) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("hyperClusterResources").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *hyperClusterResources) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("hyperClusterResources").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched hyperClusterResource.
func (c *hyperClusterResources) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.HyperClusterResource, err error) {
	result = &v1.HyperClusterResource{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("hyperClusterResources").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}