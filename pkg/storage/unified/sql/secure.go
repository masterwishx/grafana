package sql

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grafana/grafana/pkg/storage/unified/resource"
)

// UpdateSecureFields implements SecureBackend.
func (b *backend) UpdateSecureFields(ctx context.Context, key *resource.ResourceKey, fields map[string]*resource.SecureValue) *resource.ErrorResult {
	if true {
		fmt.Printf("TODO... just pretend for now: %+v\n", fields)
		return nil //
	}
	return &resource.ErrorResult{
		Message: "Not yet implemented in SQL",
		Code:    http.StatusNotImplemented,
	}
}

// ReadSecureFields implements SecureBackend.
func (b *backend) ReadSecureFields(ctx context.Context, key *resource.ResourceKey, decrypt bool) (map[string]*resource.SecureValue, *resource.ErrorResult) {
	return nil, &resource.ErrorResult{
		Message: "Not yet implemented in SQL",
		Code:    http.StatusNotImplemented,
	}
}
