package apistore

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/grafana/grafana/pkg/apimachinery/apis/common/v0alpha1"
	"github.com/grafana/grafana/pkg/apimachinery/identity"
	"github.com/grafana/grafana/pkg/apimachinery/utils"
	"github.com/grafana/grafana/pkg/storage/unified/resource"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/storage"
)

// Called on create
func (s *Storage) prepareObjectForStorage(ctx context.Context, newObject runtime.Object) ([]byte, error) {
	user, err := identity.GetRequester(ctx)
	if err != nil {
		return nil, err
	}

	obj, err := utils.MetaAccessor(newObject)
	if err != nil {
		return nil, err
	}
	if obj.GetName() == "" {
		return nil, fmt.Errorf("new object must have a name")
	}
	if obj.GetResourceVersion() != "" {
		return nil, storage.ErrResourceVersionSetOnCreate
	}
	obj.SetGenerateName("") // Clear the random name field
	obj.SetResourceVersion("")
	obj.SetSelfLink("")

	// Read+write will verify that origin format is accurate
	origin, err := obj.GetOriginInfo()
	if err != nil {
		return nil, err
	}
	obj.SetOriginInfo(origin)
	obj.SetUpdatedBy("")
	obj.SetUpdatedTimestamp(nil)
	obj.SetCreatedBy(user.GetUID())

	// Secure fields exist
	err = s.updateSecureFields(ctx, obj)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = s.codec.Encode(newObject, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Called on update
func (s *Storage) prepareObjectForUpdate(ctx context.Context, updateObject runtime.Object, previousObject runtime.Object) ([]byte, error) {
	user, err := identity.GetRequester(ctx)
	if err != nil {
		return nil, err
	}

	obj, err := utils.MetaAccessor(updateObject)
	if err != nil {
		return nil, err
	}
	if obj.GetName() == "" {
		return nil, fmt.Errorf("updated object must have a name")
	}

	previous, err := utils.MetaAccessor(previousObject)
	if err != nil {
		return nil, err
	}
	obj.SetUID(previous.GetUID())
	obj.SetCreatedBy(previous.GetCreatedBy())
	obj.SetCreationTimestamp(previous.GetCreationTimestamp())

	// Secure fields exist
	err = s.updateSecureFields(ctx, obj)
	if err != nil {
		return nil, err
	}

	// Read+write will verify that origin format is accurate
	origin, err := obj.GetOriginInfo()
	if err != nil {
		return nil, err
	}
	obj.SetOriginInfo(origin)
	obj.SetUpdatedBy(user.GetUID())
	obj.SetUpdatedTimestampMillis(time.Now().UnixMilli())

	var buf bytes.Buffer
	err = s.codec.Encode(updateObject, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Called on update
func (s *Storage) updateSecureFields(ctx context.Context, obj utils.GrafanaMetaAccessor) error {
	// Secure fields exist
	secure, ok := obj.GetSecureValues()
	if !ok || len(secure) < 1 {
		return nil
	}

	// Find the fields we need to replace
	req := &resource.WriteSecureFieldsRequest{
		Fields: make(map[string]*resource.SecureValue),
	}
	for k, v := range secure {
		if !v.IsValid() {
			return fmt.Errorf("invalid secure value: " + k)
		}
		if v.GUID != "" {
			continue // no need to update anything
		}
		if v.Ref != "" {
			return fmt.Errorf("invalid secure value: " + k + " // reference fields not yet supported")
		}
		req.Fields[k] = &resource.SecureValue{
			Guid:  v.GUID,
			Value: v.Value,
		}
	}

	if len(req.Fields) > 0 {
		req.Key = &resource.ResourceKey{
			Namespace: obj.GetNamespace(),
			Group:     s.gr.Group,
			Resource:  s.gr.Resource,
			Name:      obj.GetName(),
		}
		rsp, err := s.store.WriteSecureFields(ctx, req)
		if err != nil {
			return err
		}
		if rsp.Error != nil {
			return resource.GetError(rsp.Error)
		}
		for k, v := range rsp.Fields {
			obj.SetSecureValue(k, v0alpha1.SecureValue{
				GUID:  v.Guid,
				Value: v.Value,
				Ref:   v.Refid,
			})
		}
	}
	return nil
}
