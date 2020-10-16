package users

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mirzakhany/pm/internal/auth/users/auth"
	"github.com/mirzakhany/pm/pkg/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/mirzakhany/pm/pkg/grpcgw"
	"github.com/mirzakhany/pm/pkg/kv"
	"github.com/mirzakhany/pm/protobuf/users"
	"google.golang.org/grpc"
)

type API interface {
	grpcgw.Controller
	users.UserServiceServer
}

type api struct {
	users.UnimplementedUserServiceServer
	service Service
}

func (a api) InitRest(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux) {
	cl := users.NewUserServiceClient(conn)
	_ = users.RegisterUserServiceHandlerClient(ctx, mux, cl)
}

func (a api) InitGrpc(ctx context.Context, server *grpc.Server) {
	users.RegisterUserServiceServer(server, a)
}

func (a api) Logout(ctx context.Context, req *users.LogoutRequest) (*users.LogoutResponse, error) {
	token, err := auth.ExtractToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	err = auth.DeleteToken(token, true /* isAccessToken */)
	if err != nil {
		log.Error("failed to remove token", log.Err(err))
		return nil, status.Errorf(codes.Internal, "logout failed")
	}
	return nil, nil
}

func (a api) VerifyToken(ctx context.Context, request *users.VerifyTokenRequest) (*users.VerifyTokenResponse, error) {
	vTokens, err := auth.VerifyToken(request.AccessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	return &users.VerifyTokenResponse{AccessToken: vTokens}, nil
}

func (a api) RefreshToken(ctx context.Context, request *users.RefreshTokenRequest) (*users.RefreshTokenResponse, error) {

	vToken, err := auth.VerifyToken(request.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token")
	}

	user, err := auth.ExtractSessionUser(request.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "session expired")
	}

	dbUser, err := a.service.GetByUsername(ctx, user.Username)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid user")
	}

	if !dbUser.Enable {
		return nil, status.Errorf(codes.Unauthenticated, "user is not active")
	}
	err = auth.DeleteToken(vToken, false /* isAccessToken */)
	if err != nil {
		log.Error("failed to remove token", log.Err(err))
		return nil, status.Errorf(codes.Internal, "logout failed")
	}

	tokens, err := auth.CreateToken(dbUser)
	if err != nil {
		log.Error("error on create user token", log.String("user", dbUser.Username), log.Err(err))
		return nil, status.Errorf(codes.Internal, "internal server error, create token")
	}

	return &users.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
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
		Enable:   true,
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

func (a api) DeleteUser(ctx context.Context, request *users.DeleteUserRequest) (*empty.Empty, error) {
	_, err := a.service.Delete(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &empty.Empty{}, err
}

func New(srv Service) API {
	s := api{service: srv}
	grpcgw.RegisterController(s)
	return s
}

func init() {
	kv.Memory().SetString("/usersV1.UserService/Login", "open")
	kv.Memory().SetString("/usersV1.UserService/Register", "open")
	kv.Memory().SetString("/usersV1.UserService/VerifyToken", "open")
	kv.Memory().SetString("/usersV1.UserService/RefreshToken", "open")
}
