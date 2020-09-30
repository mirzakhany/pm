package users

import (
	"context"
	"github.com/gogo/protobuf/types"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/mirzakhany/pm/pkg/grpcgw"
	"github.com/mirzakhany/pm/pkg/kv"
	"github.com/mirzakhany/pm/pkg/log"
	"github.com/mirzakhany/pm/services/users/auth"
	users "github.com/mirzakhany/pm/services/users/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type API interface {
	grpcgw.Controller
	users.UserServiceServer
}

type api struct {
	service Service
}

func (a api) Login(ctx context.Context, request *users.LoginRequest) (*users.LoginResponse, error) {

	user, err := a.service.GetByUsername(ctx, request.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "username %s not found", request.Username)
	}

	if !auth.CheckPasswordHash(request.Password, user.Password) {
		return nil, status.Error(codes.Unauthenticated, "username or password is not valid")
	}

	tokens, err := auth.CreateToken(user)
	if err != nil {
		log.Error("error on create user token", log.String("user", user.Username), log.Err(err))
		return nil, status.Errorf(codes.Internal, "internal server error, create token")
	}

	err = auth.SaveTokens(user, tokens)
	if err != nil {
		log.Error("error on set user session", log.String("user", user.Username), log.Err(err))
		return nil, status.Errorf(codes.Internal, "internal server error, session")
	}
	return tokens, nil
}

func (a api) Register(ctx context.Context, request *users.RegisterRequest) (*users.RegisterResponse, error) {

	user, err := a.CreateUser(ctx, &users.CreateUserRequest{
		Username: request.Username,
		Password: request.Password,
		Email:    request.Email,
		Enable:   false,
	})

	if err != nil {
		return nil, status.Errorf(codes.OutOfRange, "registration failed with %s", err.Error())
	}

	return &users.RegisterResponse{
		Uuid:     user.Uuid,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (a api) InitRest(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux) {
	cl := users.NewUserServiceClient(conn)
	_ = users.RegisterUserServiceHandlerClient(ctx, mux, cl)
}

func (a api) InitGrpc(ctx context.Context, server *grpc.Server) {
	users.RegisterUserServiceServer(server, a)
}

func (a api) ListUsers(ctx context.Context, request *users.ListUsersRequest) (*users.ListUsersResponse, error) {
	offset, limit := grpcgw.GetOffsetAndLimit(request.Offset, request.Limit)
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
	request.Password, _ = auth.HashPassword(request.Password)
	res, err := a.service.Create(ctx, request)
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

func init() {
	kv.Memory().SetString("/users.UserService/Login", "open")
	kv.Memory().SetString("/users.UserService/Register", "open")
}
