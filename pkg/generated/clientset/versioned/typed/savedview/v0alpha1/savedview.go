// SPDX-License-Identifier: AGPL-3.0-only

// Code generated by client-gen. DO NOT EDIT.

package v0alpha1

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v0alpha1 "github.com/grafana/grafana/pkg/apis/savedview/v0alpha1"
	savedviewv0alpha1 "github.com/grafana/grafana/pkg/generated/applyconfiguration/savedview/v0alpha1"
	scheme "github.com/grafana/grafana/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// SavedViewsGetter has a method to return a SavedViewInterface.
// A group's client should implement this interface.
type SavedViewsGetter interface {
	SavedViews(namespace string) SavedViewInterface
}

// SavedViewInterface has methods to work with SavedView resources.
type SavedViewInterface interface {
	Create(ctx context.Context, savedView *v0alpha1.SavedView, opts v1.CreateOptions) (*v0alpha1.SavedView, error)
	Update(ctx context.Context, savedView *v0alpha1.SavedView, opts v1.UpdateOptions) (*v0alpha1.SavedView, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v0alpha1.SavedView, error)
	List(ctx context.Context, opts v1.ListOptions) (*v0alpha1.SavedViewList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v0alpha1.SavedView, err error)
	Apply(ctx context.Context, savedView *savedviewv0alpha1.SavedViewApplyConfiguration, opts v1.ApplyOptions) (result *v0alpha1.SavedView, err error)
	SavedViewExpansion
}

// savedViews implements SavedViewInterface
type savedViews struct {
	client rest.Interface
	ns     string
}

// newSavedViews returns a SavedViews
func newSavedViews(c *SavedviewV0alpha1Client, namespace string) *savedViews {
	return &savedViews{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the savedView, and returns the corresponding savedView object, and an error if there is any.
func (c *savedViews) Get(ctx context.Context, name string, options v1.GetOptions) (result *v0alpha1.SavedView, err error) {
	result = &v0alpha1.SavedView{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("savedviews").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of SavedViews that match those selectors.
func (c *savedViews) List(ctx context.Context, opts v1.ListOptions) (result *v0alpha1.SavedViewList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v0alpha1.SavedViewList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("savedviews").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested savedViews.
func (c *savedViews) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("savedviews").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a savedView and creates it.  Returns the server's representation of the savedView, and an error, if there is any.
func (c *savedViews) Create(ctx context.Context, savedView *v0alpha1.SavedView, opts v1.CreateOptions) (result *v0alpha1.SavedView, err error) {
	result = &v0alpha1.SavedView{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("savedviews").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(savedView).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a savedView and updates it. Returns the server's representation of the savedView, and an error, if there is any.
func (c *savedViews) Update(ctx context.Context, savedView *v0alpha1.SavedView, opts v1.UpdateOptions) (result *v0alpha1.SavedView, err error) {
	result = &v0alpha1.SavedView{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("savedviews").
		Name(savedView.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(savedView).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the savedView and deletes it. Returns an error if one occurs.
func (c *savedViews) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("savedviews").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *savedViews) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("savedviews").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched savedView.
func (c *savedViews) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v0alpha1.SavedView, err error) {
	result = &v0alpha1.SavedView{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("savedviews").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied savedView.
func (c *savedViews) Apply(ctx context.Context, savedView *savedviewv0alpha1.SavedViewApplyConfiguration, opts v1.ApplyOptions) (result *v0alpha1.SavedView, err error) {
	if savedView == nil {
		return nil, fmt.Errorf("savedView provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(savedView)
	if err != nil {
		return nil, err
	}
	name := savedView.Name
	if name == nil {
		return nil, fmt.Errorf("savedView.Name must be provided to Apply")
	}
	result = &v0alpha1.SavedView{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("savedviews").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
