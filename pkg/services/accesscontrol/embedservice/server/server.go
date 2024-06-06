package server

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/infra/tracing"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	authzv1 "github.com/grafana/grafana/pkg/services/accesscontrol/embedservice/proto/v1"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/grpcserver"
	"github.com/grafana/grafana/pkg/setting"
	// "github.com/grafana/grafana/pkg/services/grpcserver/interceptors"
)

var _ authzv1.AuthzServiceServer = (*Server)(nil)

type Server struct {
	authzv1.UnimplementedAuthzServiceServer

	acSvc  accesscontrol.Service
	logger log.Logger
	tracer tracing.Tracer
}

func ProvideAuthZServer(cfg *setting.Cfg, acSvc accesscontrol.Service, features *featuremgmt.FeatureManager,
	grpcServer grpcserver.Provider, registerer prometheus.Registerer, tracer tracing.Tracer) (*Server, error) {
	if !features.IsEnabledGlobally(featuremgmt.FlagAuthZGRPCServer) {
		return nil, nil
	}

	s := &Server{
		acSvc:  acSvc,
		logger: log.New("authz-grpc-server"),
		tracer: tracer,
	}

	grpcServer.GetServer().RegisterService(&authzv1.AuthzService_ServiceDesc, s)

	return s, nil
}

func (s *Server) Read(ctx context.Context, req *authzv1.ReadRequest) (*authzv1.ReadResponse, error) {
	ctx, span := s.tracer.Start(ctx, "authz.grpc.Read")
	defer span.End()

	action := req.GetAction()
	subject := req.GetSubject()
	stackID := req.GetStackId() // TODO can we consider the stackID as the orgID?

	ctxLogger := s.logger.FromContext(ctx)
	ctxLogger.Debug("Read", "action", action, "subject", subject, "stackID", stackID)

	permissions, err := s.acSvc.SearchUserPermissions(ctx, stackID, accesscontrol.SearchOptions{NamespacedID: subject, Action: action})
	if err != nil {
		ctxLogger.Error("failed to search user permissions", "error", err)
		return nil, tracing.Errorf(span, "failed to search user permissions: %w", err)
	}

	data := make([]*authzv1.ReadResponse_Data, 0, len(permissions))
	for _, perm := range permissions {
		data = append(data, &authzv1.ReadResponse_Data{Object: perm.Scope})
	}
	return &authzv1.ReadResponse{Data: data}, nil
}