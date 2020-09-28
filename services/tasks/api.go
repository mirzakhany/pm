package tasks

import (
	"context"

	"github.com/gogo/protobuf/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/mirzakhany/pm/pkg/grpcgw"
	tasks "github.com/mirzakhany/pm/services/tasks/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type API interface {
	grpcgw.Controller
	tasks.TaskServiceServer
}

type api struct {
	service Service
}

func (a api) InitRest(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux) {
	cl := tasks.NewTaskServiceClient(conn)
	_ = tasks.RegisterTaskServiceHandlerClient(ctx, mux, cl)
}

func (a api) InitGrpc(ctx context.Context, server *grpc.Server) {
	tasks.RegisterTaskServiceServer(server, a)
}

func (a api) ListTasks(ctx context.Context, request *tasks.ListTasksRequest) (*tasks.ListTasksResponse, error) {
	offset, limit := grpcgw.GetOffsetAndLimit(request.Offset, request.Limit)
	res, err := a.service.Query(ctx, offset, limit)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) GetTask(ctx context.Context, request *tasks.GetTaskRequest) (*tasks.Task, error) {
	res, err := a.service.Get(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) CreateTask(ctx context.Context, request *tasks.CreateTaskRequest) (*tasks.Task, error) {
	res, err := a.service.Create(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) UpdateTask(ctx context.Context, request *tasks.UpdateTaskRequest) (*tasks.Task, error) {
	res, err := a.service.Update(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) DeleteTask(ctx context.Context, request *tasks.DeleteTaskRequest) (*types.Empty, error) {
	_, err := a.service.Delete(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &types.Empty{}, err
}

func New(srv Service) API {
	s := api{service: srv}
	grpcgw.RegisterController(s)
	return s
}
