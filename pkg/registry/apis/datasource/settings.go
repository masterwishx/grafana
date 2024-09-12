package datasource

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"

	"github.com/grafana/grafana/pkg/apimachinery/utils"
	grafanaregistry "github.com/grafana/grafana/pkg/apiserver/registry/generic"
)

func newSettingsStorage(scheme *runtime.Scheme, resourceInfo utils.ResourceInfo, optsGetter generic.RESTOptionsGetter) (*genericregistry.Store, error) {
	strategy := grafanaregistry.NewStrategy(scheme, resourceInfo.GroupVersion())
	store := &genericregistry.Store{
		NewFunc:                   resourceInfo.NewFunc,
		NewListFunc:               resourceInfo.NewListFunc,
		KeyRootFunc:               grafanaregistry.KeyRootFunc(resourceInfo.GroupResource()),
		KeyFunc:                   grafanaregistry.NamespaceKeyFunc(resourceInfo.GroupResource()),
		PredicateFunc:             grafanaregistry.Matcher,
		DefaultQualifiedResource:  resourceInfo.GroupResource(),
		SingularQualifiedResource: resourceInfo.SingularGroupResource(),
		CreateStrategy:            strategy,
		UpdateStrategy:            strategy,
		DeleteStrategy:            strategy,
		TableConvertor:            resourceInfo.TableConverter(),
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: grafanaregistry.GetAttrs}
	err := store.CompleteWithOptions(options)
	return store, err
}
