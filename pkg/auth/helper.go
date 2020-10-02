package auth

import (
	"context"
	"errors"

	users "github.com/mirzakhany/pm/services/users/proto"
)

type contextKey int

const (
	userKey contextKey = iota + 1
)

// ExtractUser try to extract the current user from the context
func ExtractUser(ctx context.Context) (*users.User, error) {
	u, ok := ctx.Value(userKey).(*users.User)
	if !ok {
		return nil, errors.New("no user in context")
	}
	return u, nil
}

// ContextWithUser return context with user
func ContextWithUser(ctx context.Context, user *users.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}
