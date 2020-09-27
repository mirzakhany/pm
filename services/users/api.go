package users

import (
	"context"
	"github.com/gogo/protobuf/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"proj/pkg/grpcgw"
	users "proj/services/users/proto"
)

type API interface {
	grpcgw.Controller
	users.UserServiceServer
}

type api struct {
	service Service
}

func (a api) Login(ctx context.Context, request *users.LoginRequest) (*users.LoginResponse, error) {
	panic("implement me")
}

func (a api) Register(ctx context.Context, request *users.RegisterRequest) (*users.RegisterResponse, error) {
	panic("implement me")
}

func (a api) InitRest(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux) {
	cl := users.NewUserServiceClient(conn)
	_ = users.RegisterUserServiceHandlerClient(ctx, mux, cl)
}

func (a api) InitGrpc(ctx context.Context, server *grpc.Server) {
	users.RegisterUserServiceServer(server, a)
}

func (a api) ListUsers(ctx context.Context, request *users.ListUsersRequest) (*users.ListUsersResponse, error) {
	offset, limit := grpcgw.GetOffsetAndLimit(request.Limit, request.Offset)
	res, err := a.service.Query(ctx, offset, limit)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) GetUser(ctx context.Context, request *users.GetUserRequest) (*users.User, error) {
	res, err := a.service.Get(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) CreateUser(ctx context.Context, request *users.CreateUserRequest) (*users.User, error) {
	res,err := a.service.Create(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) UpdateUser(ctx context.Context, request *users.UpdateUserRequest) (*users.User, error) {
	res, err := a.service.Update(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) DeleteUser(ctx context.Context, request *users.DeleteUserRequest) (*types.Empty, error) {
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
