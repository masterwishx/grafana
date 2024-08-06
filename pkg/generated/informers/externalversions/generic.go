// SPDX-License-Identifier: AGPL-3.0-only

// Code generated by informer-gen. DO NOT EDIT.

package externalversions

import (
	"fmt"

	v0alpha1 "github.com/grafana/grafana/pkg/apis/alerting_notifications/v0alpha1"
	savedviewv0alpha1 "github.com/grafana/grafana/pkg/apis/savedview/v0alpha1"
	servicev0alpha1 "github.com/grafana/grafana/pkg/apis/service/v0alpha1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=notifications.alerting.grafana.app, Version=v0alpha1
	case v0alpha1.SchemeGroupVersion.WithResource("receivers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Notifications().V0alpha1().Receivers().Informer()}, nil
	case v0alpha1.SchemeGroupVersion.WithResource("timeintervals"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Notifications().V0alpha1().TimeIntervals().Informer()}, nil

		// Group=savedview.grafana.app, Version=v0alpha1
	case savedviewv0alpha1.SchemeGroupVersion.WithResource("savedviews"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Savedview().V0alpha1().SavedViews().Informer()}, nil

		// Group=service.grafana.app, Version=v0alpha1
	case servicev0alpha1.SchemeGroupVersion.WithResource("externalnames"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Service().V0alpha1().ExternalNames().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
