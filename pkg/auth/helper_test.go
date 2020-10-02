package auth

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
	users "github.com/mirzakhany/pm/services/users/proto"
)

func TestAuthHelper(t *testing.T) {
	now := time.Now()
	Uuid := uuid.New().String()
	user := &users.User{
		Id:        1,
		Uuid:      Uuid,
		Username:  "test",
		Password:  "test",
		Email:     "user@web.com",
		Enable:    true,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	ctx := context.Background()
	ctx = ContextWithUser(ctx, user)
	assert.NotNil(t, ctx)

	user1, err := ExtractUser(ctx)
	assert.Nil(t, err)
	assert.Equal(t, user, user1)

	// test not found user
	ctx1 := context.Background()
	_, err1 := ExtractUser(ctx1)
	assert.NotNil(t, err1)
}
