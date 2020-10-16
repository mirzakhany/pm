package workspaces

import (
	"context"
	"errors"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/mirzakhany/pm/internal/entity"
	"github.com/mirzakhany/pm/protobuf/workspaces"
	"github.com/stretchr/testify/assert"
)

var errCRUD = errors.New("error crud")

func TestCreateWorkspaceRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     workspaces.CreateWorkspaceRequest
		wantError bool
	}{
		{"success", workspaces.CreateWorkspaceRequest{Title: "test", Domain: "example"}, false},
		{"required", workspaces.CreateWorkspaceRequest{Title: "", Domain: "example"}, true},
		{"too long", workspaces.CreateWorkspaceRequest{Domain: "example", Title: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateWorkspaceRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     workspaces.UpdateWorkspaceRequest
		wantError bool
	}{
		{"success", workspaces.UpdateWorkspaceRequest{Title: "test", Domain: "example"}, false},
		{"required", workspaces.UpdateWorkspaceRequest{Title: "", Domain: "example"}, true},
		{"too long", workspaces.UpdateWorkspaceRequest{Domain: "example", Title: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
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
	workspace, err := s.Create(ctx, &workspaces.CreateWorkspaceRequest{Title: "test", Domain: "example"})
	assert.Nil(t, err)
	assert.NotEmpty(t, workspace.Uuid)
	id := workspace.Uuid
	assert.Equal(t, "test", workspace.Title)
	assert.NotEmpty(t, workspace.CreatedAt)
	assert.NotEmpty(t, workspace.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// validation error in creation
	_, err = s.Create(ctx, &workspaces.CreateWorkspaceRequest{Title: "", Domain: "example"})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// unexpected error in creation
	_, err = s.Create(ctx, &workspaces.CreateWorkspaceRequest{Title: "error", Domain: "example"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	_, _ = s.Create(ctx, &workspaces.CreateWorkspaceRequest{Title: "test2", Domain: "example"})

	// update
	workspace, err = s.Update(ctx, &workspaces.UpdateWorkspaceRequest{Title: "test updated", Domain: "example", Uuid: id})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", workspace.Title)
	_, err = s.Update(ctx, &workspaces.UpdateWorkspaceRequest{Title: "test updated", Uuid: "none", Domain: "example"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, &workspaces.UpdateWorkspaceRequest{Title: "", Uuid: id, Domain: "example"})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// unexpected error in update
	_, err = s.Update(ctx, &workspaces.UpdateWorkspaceRequest{Title: "error", Uuid: id, Domain: "example"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	workspace, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", workspace.Title)
	assert.Equal(t, id, workspace.Uuid)

	// query
	_workspaces, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, int(_workspaces.TotalCount))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	workspace, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, workspace.Uuid)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)
}

type mockRepository struct {
	items []entity.Workspace
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.Workspace, error) {
	for _, item := range m.items {
		if item.UUID == id {
			return item, nil
		}
	}
	return entity.Workspace{}, pg.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int64) ([]entity.Workspace, int, error) {
	return m.items, len(m.items), nil
}

func (m *mockRepository) Create(ctx context.Context, workspace entity.Workspace) error {
	if workspace.Title == "error" {
		return errCRUD
	}
	m.items = append(m.items, workspace)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, workspace entity.Workspace) error {
	if workspace.Title == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.UUID == workspace.UUID {
			m.items[i] = workspace
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
