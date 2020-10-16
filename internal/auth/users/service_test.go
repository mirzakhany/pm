package users

import (
	"context"
	"testing"

	users "github.com/mirzakhany/pm/protobuf/users"
	"github.com/stretchr/testify/assert"
)

func TestCreateUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     users.CreateUserRequest
		wantError bool
	}{
		{"success", users.CreateUserRequest{Username: "test", Password: "test", Email: "test@test.com", Enable: true}, false},
		{"required", users.CreateUserRequest{Username: "", Password: "test", Email: "test@test.com", Enable: true}, true},
		{"too long", users.CreateUserRequest{Password: "test", Email: "test@test.com", Enable: true, Username: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     users.UpdateUserRequest
		wantError bool
	}{
		{"success", users.UpdateUserRequest{Username: "test", Email: "test@test.com", Enable: true}, false},
		{"required", users.UpdateUserRequest{Username: "", Email: "test@test.com", Enable: true}, true},
		{"too long", users.UpdateUserRequest{Email: "test@test.com", Enable: true, Username: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpdateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func Test_service_CRUD(t *testing.T) {
	s := NewService(&mockRepository{})
	ctx := context.Background()

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, int64(0), count)

	// successful creation
	user, err := s.Create(ctx, &users.CreateUserRequest{Username: "test", Password: "test", Email: "test@test.com", Enable: true})
	assert.Nil(t, err)
	assert.NotEmpty(t, user.Uuid)
	id := user.Uuid
	assert.Equal(t, "test", user.Username)
	assert.NotEmpty(t, user.CreatedAt)
	assert.NotEmpty(t, user.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// validation error in creation
	_, err = s.Create(ctx, &users.CreateUserRequest{Username: "", Password: "test", Email: "test@test.com", Enable: true})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// unexpected error in creation
	_, err = s.Create(ctx, &users.CreateUserRequest{Username: "error", Password: "test", Email: "test@test.com", Enable: true})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	_, _ = s.Create(ctx, &users.CreateUserRequest{Username: "test2", Password: "test", Email: "test@test.com", Enable: true})

	// update
	user, err = s.Update(ctx, &users.UpdateUserRequest{Username: "test-updated", Email: "test@test.com", Enable: true, Uuid: id})
	assert.Nil(t, err)
	assert.Equal(t, "test-updated", user.Username)
	_, err = s.Update(ctx, &users.UpdateUserRequest{Username: "test-updated", Email: "test@test.com", Enable: true, Uuid: "none"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, &users.UpdateUserRequest{Username: "", Email: "test@test.com", Enable: true, Uuid: id})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// unexpected error in update
	_, err = s.Update(ctx, &users.UpdateUserRequest{Username: "error", Email: "test@test.com", Enable: true, Uuid: id})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	user, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test-updated", user.Username)
	assert.Equal(t, id, user.Uuid)

	// query
	_users, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, int(_users.TotalCount))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	user, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, user.Uuid)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)
}
