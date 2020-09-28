package tasks

import (
	"context"
	"errors"
	"github.com/go-pg/pg/v10"
	"github.com/mirzakhany/pm/services/tasks/proto"
	userSrv "github.com/mirzakhany/pm/services/users"
	usersProto "github.com/mirzakhany/pm/services/users/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

var errCRUD = errors.New("error crud")

func TestCreateTaskRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     tasks.CreateTaskRequest
		wantError bool
	}{
		{"success", tasks.CreateTaskRequest{
			Title:       "test",
			Description: "this is a test",
			Status:      0,
			SprintId:    0,
			Estimate:    0,
		}, false},
		{"required", tasks.CreateTaskRequest{
			Title:       "",
			Description: "this is a test",
			Status:      0,
			SprintId:    0,
			Estimate:    0,
		}, true},
		{"too long", tasks.CreateTaskRequest{
			Description: "this is a test",
			Status:      0,
			SprintId:    0,
			Estimate:    0,
			Title:       "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateRequest(tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateTaskRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     tasks.UpdateTaskRequest
		wantError bool
	}{
		{"success", tasks.UpdateTaskRequest{
			Title:       "test",
			Description: "this is a test",
			Status:      0,
			SprintId:    0,
			Estimate:    0,
		}, false},
		{"required", tasks.UpdateTaskRequest{
			Title:       "",
			Description: "this is a test",
			Status:      0,
			SprintId:    0,
			Estimate:    0,
		}, true},
		{"too long", tasks.UpdateTaskRequest{
			Description: "this is a test",
			Status:      0,
			SprintId:    0,
			Estimate:    0,
			Title:       "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpdateRequest(tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func Test_service_CRUD(t *testing.T) {
	userServices := userSrv.NewServiceForTest()
	s := NewService(&mockRepository{}, userServices)
	ctx := context.Background()

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, int64(0), count)

	user1, err := userServices.Create(ctx, &usersProto.CreateUserRequest{
		Username: "test", Password: "test", Email: "test@test.com", Enable: true,
	})
	assert.Nil(t, err)


	user2, err := userServices.Create(ctx, &usersProto.CreateUserRequest{
		Username: "test", Password: "test", Email: "test@test.com", Enable: true,
	})
	assert.Nil(t, err)
	// successful creation
	task, err := s.Create(ctx, &tasks.CreateTaskRequest{
		Title:       "test",
		Description: "this is a test",
		CreatorUuid: user1.Uuid,
		AssigneeUuid: user2.Uuid,
		Status:      0,
		SprintId:    0,
		Estimate:    0,
	})

	assert.Nil(t, err)
	assert.NotEmpty(t, task.Uuid)
	id := task.Uuid
	assert.Equal(t, "test", task.Title)
	assert.NotEmpty(t, task.CreatedAt)
	assert.NotEmpty(t, task.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// validation error in creation
	_, err = s.Create(ctx, &tasks.CreateTaskRequest{
		Title:       "",
		Description: "this is a test",
		Status:      0,
		SprintId:    0,
		Estimate:    0,
	})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// unexpected error in creation
	_, err = s.Create(ctx, &tasks.CreateTaskRequest{
		Title:       "error",
		Description: "this is a test",
		CreatorUuid: user1.Uuid,
		AssigneeUuid: user2.Uuid,
		Status:      0,
		SprintId:    0,
		Estimate:    0,
	})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	_, _ = s.Create(ctx, &tasks.CreateTaskRequest{
		Title:       "test2",
		Description: "this is a test",
		CreatorUuid: user1.Uuid,
		AssigneeUuid: user2.Uuid,
		Status:      0,
		SprintId:    0,
		Estimate:    0,
	})

	// update
	task, err = s.Update(ctx, &tasks.UpdateTaskRequest{
		Uuid: id,
		Title:       "test-updated",
		Description: "this is a test",
		CreatorUuid: user1.Uuid,
		AssigneeUuid: user2.Uuid,
		Status:      0,
		SprintId:    0,
		Estimate:    0,
	})
	assert.Nil(t, err)
	assert.Equal(t, "test-updated", task.Title)
	_, err = s.Update(ctx, &tasks.UpdateTaskRequest{
		Uuid: "none",
		Title:       "test-updated",
		Description: "this is a test",
		CreatorUuid: user1.Uuid,
		AssigneeUuid: user2.Uuid,
		Status:      0,
		SprintId:    0,
		Estimate:    0,
	})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, &tasks.UpdateTaskRequest{
		Uuid: id,
		Title:       "",
		Description: "this is a test",
		CreatorUuid: user1.Uuid,
		AssigneeUuid: user2.Uuid,
		Status:      0,
		SprintId:    0,
		Estimate:    0,
	})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// unexpected error in update
	_, err = s.Update(ctx, &tasks.UpdateTaskRequest{
		Uuid: id,
		Title:       "error",
		Description: "this is a test",
		CreatorUuid: user1.Uuid,
		AssigneeUuid: user2.Uuid,
		Status:      0,
		SprintId:    0,
		Estimate:    0,
	})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	task, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test-updated", task.Title)
	assert.Equal(t, id, task.Uuid)

	// query
	_tasks, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, int(_tasks.TotalCount))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	task, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, task.Uuid)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)
}

type mockRepository struct {
	items []TaskModel
}

func (m mockRepository) Get(ctx context.Context, id string) (TaskModel, error) {
	for _, item := range m.items {
		if item.UUID == id {
			return item, nil
		}
	}
	return TaskModel{}, pg.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int64) ([]TaskModel, int, error) {
	return m.items, len(m.items), nil
}

func (m *mockRepository) Create(ctx context.Context, task TaskModel) error {
	if task.Title == "error" {
		return errCRUD
	}
	m.items = append(m.items, task)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, task TaskModel) error {
	if task.Title == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.UUID == task.UUID {
			m.items[i] = task
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
