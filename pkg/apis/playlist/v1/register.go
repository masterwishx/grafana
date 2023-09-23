package v1

import (
	"github.com/grafana/grafana/pkg/apis"
	"github.com/grafana/grafana/pkg/services/playlist"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	common "k8s.io/kube-openapi/pkg/common"
	"k8s.io/kube-openapi/pkg/spec3"
)

// GroupName is the group name for this API.
const GroupName = "playlist.x.grafana.com"
const VersionID = "v0-alpha" //
const APIVersion = GroupName + "/" + VersionID

// This is used just so wire has something unique to return
type PlaylistDummyService struct{}

func RegisterAPIService(c apis.GroupBuilderCollection, p playlist.Service) *PlaylistDummyService {
	c.AddAPI(&builder{
		service: p,
	})
	return &PlaylistDummyService{}
}

type builder struct {
	service playlist.Service
}

func (b *builder) InstallSchema(scheme *runtime.Scheme) error {
	err := AddToScheme(scheme)
	if err != nil {
		return err
	}
	return scheme.SetVersionPriority(SchemeGroupVersion)
}

func (b *builder) GetAPIGroupInfo(
	scheme *runtime.Scheme,
	codecs serializer.CodecFactory, // pointer?
) *genericapiserver.APIGroupInfo {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(GroupName, scheme, metav1.ParameterCodec, codecs)
	storage := map[string]rest.Storage{}
	storage["playlists"] = &handler{
		service: b.service,
	}

	apiGroupInfo.VersionedResourcesStorageMap[VersionID] = storage
	return &apiGroupInfo
}

func (b *builder) GetOpenAPIDefinitions() common.GetOpenAPIDefinitions {
	return getOpenAPIDefinitions
}

// Register additional routes with the server
func (b *builder) GetOpenAPIPostProcessor() func(*spec3.OpenAPI) (*spec3.OpenAPI, error) {
	return nil
}

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: VersionID}

// Resource takes an unqualified resource and returns a Group qualified GroupResource
func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	// SchemeBuilder points to a list of functions added to Scheme.
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	localSchemeBuilder = &SchemeBuilder
	// AddToScheme is a common registration function for mapping packaged scoped group & version keys to a scheme.
	AddToScheme = localSchemeBuilder.AddToScheme
)

// Adds the list of known types to the given scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Playlist{},
		&PlaylistList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
