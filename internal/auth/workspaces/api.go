package workspaces

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/mirzakhany/pm/pkg/grpcgw"
	"github.com/mirzakhany/pm/protobuf/workspaces"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type API interface {
	grpcgw.Controller
	workspaces.WorkspaceServiceServer
}

type api struct {
	service Service
}

func (a api) InitRest(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux) {
	cl := workspaces.NewWorkspaceServiceClient(conn)
	_ = workspaces.RegisterWorkspaceServiceHandlerClient(ctx, mux, cl)
}

func (a api) InitGrpc(ctx context.Context, server *grpc.Server) {
	workspaces.RegisterWorkspaceServiceServer(server, a)
}

func (a api) ListWorkspaces(ctx context.Context, request *workspaces.ListWorkspacesRequest) (*workspaces.ListWorkspacesResponse, error) {
	offset, limit := grpcgw.GetOffsetAndLimit(request.Offset, request.Limit)
	res, err := a.service.Query(ctx, offset, limit)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) GetWorkspace(ctx context.Context, request *workspaces.GetWorkspaceRequest) (*workspaces.Workspace, error) {
	res, err := a.service.Get(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) CreateWorkspace(ctx context.Context, request *workspaces.CreateWorkspaceRequest) (*workspaces.Workspace, error) {
	res, err := a.service.Create(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) UpdateWorkspace(ctx context.Context, request *workspaces.UpdateWorkspaceRequest) (*workspaces.Workspace, error) {
	res, err := a.service.Update(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) DeleteWorkspace(ctx context.Context, request *workspaces.DeleteWorkspaceRequest) (*empty.Empty, error) {
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
