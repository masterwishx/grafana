// SPDX-License-Identifier: AGPL-3.0-only

// Code generated by informer-gen. DO NOT EDIT.

package v0alpha1

import (
	internalinterfaces "github.com/grafana/grafana/pkg/generated/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// SavedViews returns a SavedViewInformer.
	SavedViews() SavedViewInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// SavedViews returns a SavedViewInformer.
func (v *version) SavedViews() SavedViewInformer {
	return &savedViewInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
