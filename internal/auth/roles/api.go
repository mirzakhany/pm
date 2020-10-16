package roles

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/mirzakhany/pm/pkg/grpcgw"
	"github.com/mirzakhany/pm/protobuf/roles"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type API interface {
	grpcgw.Controller
	roles.RoleServiceServer
}

type api struct {
	service Service
}

func (a api) InitRest(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux) {
	cl := roles.NewRoleServiceClient(conn)
	_ = roles.RegisterRoleServiceHandlerClient(ctx, mux, cl)
}

func (a api) InitGrpc(ctx context.Context, server *grpc.Server) {
	roles.RegisterRoleServiceServer(server, a)
}

func (a api) ListRoles(ctx context.Context, request *roles.ListRolesRequest) (*roles.ListRolesResponse, error) {
	offset, limit := grpcgw.GetOffsetAndLimit(request.Offset, request.Limit)
	res, err := a.service.Query(ctx, offset, limit)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) GetRole(ctx context.Context, request *roles.GetRoleRequest) (*roles.Role, error) {
	res, err := a.service.Get(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) CreateRole(ctx context.Context, request *roles.CreateRoleRequest) (*roles.Role, error) {
	res, err := a.service.Create(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) UpdateRole(ctx context.Context, request *roles.UpdateRoleRequest) (*roles.Role, error) {
	res, err := a.service.Update(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) DeleteRole(ctx context.Context, request *roles.DeleteRoleRequest) (*empty.Empty, error) {
	_, err := a.service.Delete(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return nil, err
}

func New(srv Service) API {
	s := api{service: srv}
	grpcgw.RegisterController(s)
	return s
}
