package auth

import (
	"context"
	"errors"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/mirzakhany/pm/pkg/kv"
	"github.com/mirzakhany/pm/pkg/session"
	users "github.com/mirzakhany/pm/services/users/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mirzakhany/pm/pkg/grpcgw"
	"google.golang.org/grpc"
)

type contextKey int

const (
	resourceKey contextKey = iota
	userKey
	tokenKey
	fullMethodKey
)

func streamExtractor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return handler(srv, ss)
}

func unaryExtractor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	res, ok := kv.Memory().Get(info.FullMethod)
	if !ok {
		ctx = context.WithValue(ctx, resourceKey, res)
	}
	ctx = context.WithValue(ctx, fullMethodKey, info.FullMethod)
	return handler(ctx, req)
}

func authHandler(ctx context.Context) (context.Context, error) {
	r := ctx.Value(resourceKey)
	if r == nil { // No user requested here
		return ctx, nil
	}
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return ctx, status.Errorf(codes.InvalidArgument, "invalid token format")
	}
	var user *users.User
	err = session.Get(token, &user)
	if err != nil {
		return ctx, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	return context.WithValue(context.WithValue(ctx, userKey, user), tokenKey, token), nil
}

// ExtractUser try to extract the current user from the context
func ExtractUser(ctx context.Context) (*users.User, error) {
	u, ok := ctx.Value(userKey).(*users.User)
	if !ok {
		return nil, errors.New("no user in context")
	}
	return u, nil
}

// ExtractToken try to extract token from context
func ExtractToken(ctx context.Context) (string, error) {
	tok, ok := ctx.Value(tokenKey).(string)
	if !ok {
		return "", errors.New("no token in context")
	}
	return tok, nil
}

func init() {
	grpcgw.RegisterInterceptors(grpcgw.Interceptor{
		Unary:  unaryExtractor,
		Stream: streamExtractor,
	})

	grpcgw.RegisterInterceptors(grpcgw.Interceptor{
		Unary:  grpc_auth.UnaryServerInterceptor(authHandler),
		Stream: grpc_auth.StreamServerInterceptor(authHandler),
	})
}
