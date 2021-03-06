// Copyright 2018 The Gardener Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	v1 "github.com/gardener/gardenctl/pkg/apis/garden/v1"
	scheme "github.com/gardener/gardenctl/pkg/client/garden/clientset/versioned/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ShootsGetter has a method to return a ShootInterface.
// A group's client should implement this interface.
type ShootsGetter interface {
	Shoots(namespace string) ShootInterface
}

// ShootInterface has methods to work with Shoot resources.
type ShootInterface interface {
	Create(*v1.Shoot) (*v1.Shoot, error)
	Update(*v1.Shoot) (*v1.Shoot, error)
	UpdateStatus(*v1.Shoot) (*v1.Shoot, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.Shoot, error)
	List(opts meta_v1.ListOptions) (*v1.ShootList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Shoot, err error)
	ShootExpansion
}

// shoots implements ShootInterface
type shoots struct {
	client rest.Interface
	ns     string
}

// newShoots returns a Shoots
func newShoots(c *GardenV1Client, namespace string) *shoots {
	return &shoots{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the shoot, and returns the corresponding shoot object, and an error if there is any.
func (c *shoots) Get(name string, options meta_v1.GetOptions) (result *v1.Shoot, err error) {
	result = &v1.Shoot{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("shoots").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Shoots that match those selectors.
func (c *shoots) List(opts meta_v1.ListOptions) (result *v1.ShootList, err error) {
	result = &v1.ShootList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("shoots").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested shoots.
func (c *shoots) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("shoots").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a shoot and creates it.  Returns the server's representation of the shoot, and an error, if there is any.
func (c *shoots) Create(shoot *v1.Shoot) (result *v1.Shoot, err error) {
	result = &v1.Shoot{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("shoots").
		Body(shoot).
		Do().
		Into(result)
	return
}

// Update takes the representation of a shoot and updates it. Returns the server's representation of the shoot, and an error, if there is any.
func (c *shoots) Update(shoot *v1.Shoot) (result *v1.Shoot, err error) {
	result = &v1.Shoot{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("shoots").
		Name(shoot.Name).
		Body(shoot).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *shoots) UpdateStatus(shoot *v1.Shoot) (result *v1.Shoot, err error) {
	result = &v1.Shoot{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("shoots").
		Name(shoot.Name).
		SubResource("status").
		Body(shoot).
		Do().
		Into(result)
	return
}

// Delete takes name of the shoot and deletes it. Returns an error if one occurs.
func (c *shoots) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("shoots").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *shoots) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("shoots").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched shoot.
func (c *shoots) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Shoot, err error) {
	result = &v1.Shoot{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("shoots").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
