package roles

import (
	"context"
	"errors"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/mirzakhany/pm/internal/entity"
	"github.com/mirzakhany/pm/protobuf/roles"
	"github.com/stretchr/testify/assert"
)

var errCRUD = errors.New("error crud")

func TestCreateRoleRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     roles.CreateRoleRequest
		wantError bool
	}{
		{"success", roles.CreateRoleRequest{Title: "test"}, false},
		{"required", roles.CreateRoleRequest{Title: ""}, true},
		{"too long", roles.CreateRoleRequest{Title: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateRoleRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     roles.UpdateRoleRequest
		wantError bool
	}{
		{"success", roles.UpdateRoleRequest{Title: "test"}, false},
		{"required", roles.UpdateRoleRequest{Title: ""}, true},
		{"too long", roles.UpdateRoleRequest{Title: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
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
	role, err := s.Create(ctx, &roles.CreateRoleRequest{Title: "test"})
	assert.Nil(t, err)
	assert.NotEmpty(t, role.Uuid)
	id := role.Uuid
	assert.Equal(t, "test", role.Title)
	assert.NotEmpty(t, role.CreatedAt)
	assert.NotEmpty(t, role.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// validation error in creation
	_, err = s.Create(ctx, &roles.CreateRoleRequest{Title: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// unexpected error in creation
	_, err = s.Create(ctx, &roles.CreateRoleRequest{Title: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	_, _ = s.Create(ctx, &roles.CreateRoleRequest{Title: "test2"})

	// update
	role, err = s.Update(ctx, &roles.UpdateRoleRequest{Title: "test updated", Uuid: id})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", role.Title)
	_, err = s.Update(ctx, &roles.UpdateRoleRequest{Title: "test updated", Uuid: "none"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, &roles.UpdateRoleRequest{Title: "", Uuid: id})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// unexpected error in update
	_, err = s.Update(ctx, &roles.UpdateRoleRequest{Title: "error", Uuid: id})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	role, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", role.Title)
	assert.Equal(t, id, role.Uuid)

	// query
	_roles, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, int(_roles.TotalCount))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	role, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, role.Uuid)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)
}

type mockRepository struct {
	items []entity.Role
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.Role, error) {
	for _, item := range m.items {
		if item.UUID == id {
			return item, nil
		}
	}
	return entity.Role{}, pg.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int64) ([]entity.Role, int, error) {
	return m.items, len(m.items), nil
}

func (m *mockRepository) Create(ctx context.Context, role entity.Role) error {
	if role.Title == "error" {
		return errCRUD
	}
	m.items = append(m.items, role)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, role entity.Role) error {
	if role.Title == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.UUID == role.UUID {
			m.items[i] = role
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	for i, item := range m.items {
		if item.UUID == id {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}
