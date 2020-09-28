package tasks

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	tasksProto "github.com/mirzakhany/pm/services/tasks/proto"
	"github.com/mirzakhany/pm/services/users"
	usersProto "github.com/mirzakhany/pm/services/users/proto"
)

// Service encapsulates use case logic for tasks.
type Service interface {
	Get(ctx context.Context, uuid string) (*tasksProto.Task, error)
	Query(ctx context.Context, offset, limit int64) (*tasksProto.ListTasksResponse, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, input *tasksProto.CreateTaskRequest) (*tasksProto.Task, error)
	Update(ctx context.Context, input *tasksProto.UpdateTaskRequest) (*tasksProto.Task, error)
	Delete(ctx context.Context, uuid string) (*tasksProto.Task, error)
}

// ValidateCreateRequest validates the CreateTaskRequest fields.
func ValidateCreateRequest(c tasksProto.CreateTaskRequest) error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Title, validation.Required, validation.Length(0, 128)),
		validation.Field(&c.Description, validation.Required, validation.Length(0, 1000)),
	)
}

// Validate validates the UpdateTaskRequest fields.
func ValidateUpdateRequest(u tasksProto.UpdateTaskRequest) error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Title, validation.Required, validation.Length(0, 128)),
		validation.Field(&u.Description, validation.Required, validation.Length(0, 1000)),
	)
}

type service struct {
	repo    Repository
	userSrv users.Service
}

// NewService creates a new task service.
func NewService(repo Repository, userSrv users.Service) Service {
	return service{repo, userSrv}
}

// Get returns the task with the specified the task UUID.
func (s service) Get(ctx context.Context, UUID string) (*tasksProto.Task, error) {
	task, err := s.repo.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	return &tasksProto.Task{
		Uuid:        task.UUID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		SprintId:    task.SprintID,
		Assignee: &usersProto.User{
			Uuid:      task.Assignee.UUID,
			Username:  task.Assignee.Username,
			Email:     task.Assignee.Email,
			Enable:    task.Assignee.Enable,
			CreatedAt: &task.Assignee.CreatedAt,
			UpdatedAt: &task.Assignee.UpdatedAt,
		},
		Creator: &usersProto.User{
			Uuid:      task.Creator.UUID,
			Username:  task.Creator.Username,
			Email:     task.Creator.Email,
			Enable:    task.Creator.Enable,
			CreatedAt: &task.Creator.CreatedAt,
			UpdatedAt: &task.Creator.UpdatedAt,
		},
		CreatedAt: &task.CreatedAt,
		UpdatedAt: &task.UpdatedAt,
	}, nil
}

// Create creates a new task.
func (s service) Create(ctx context.Context, req *tasksProto.CreateTaskRequest) (*tasksProto.Task, error) {
	if err := ValidateCreateRequest(*req); err != nil {
		return nil, err
	}

	creator, err := s.userSrv.Get(ctx, req.CreatorUuid)
	if err != nil {
		return nil, err
	}

	assignee, err := s.userSrv.Get(ctx, req.AssigneeUuid)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	id := uuid.New().String()
	err = s.repo.Create(ctx, TaskModel{
		UUID:        id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		SprintID:    req.SprintId,
		Estimate:    req.Estimate,
		AssigneeID:  assignee.Id,
		CreatorID:   creator.Id,
		Assignee: &users.UserModel{
			ID:        assignee.Id,
			UUID:      assignee.Uuid,
			Username:  assignee.Username,
			Password:  assignee.Password,
			Email:     assignee.Email,
			Enable:    assignee.Enable,
			CreatedAt: *assignee.CreatedAt,
			UpdatedAt: *assignee.UpdatedAt,
		},
		Creator: &users.UserModel{
			ID:        creator.Id,
			UUID:      creator.Uuid,
			Username:  creator.Username,
			Password:  creator.Password,
			Email:     creator.Email,
			Enable:    creator.Enable,
			CreatedAt: *creator.CreatedAt,
			UpdatedAt: *creator.UpdatedAt,
		},
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// Update updates the task with the specified UUID.
func (s service) Update(ctx context.Context, req *tasksProto.UpdateTaskRequest) (*tasksProto.Task, error) {
	if err := ValidateUpdateRequest(*req); err != nil {
		return nil, err
	}

	task, err := s.Get(ctx, req.Uuid)
	if err != nil {
		return task, err
	}
	now := time.Now()

	creator, err := s.userSrv.Get(ctx, req.CreatorUuid)
	if err != nil {
		return nil, err
	}

	assignee, err := s.userSrv.Get(ctx, req.AssigneeUuid)
	if err != nil {
		return nil, err
	}

	taskModel := TaskModel{
		ID:          task.Id,
		UUID:        task.Uuid,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		SprintID:    req.SprintId,
		Estimate:    req.Estimate,
		AssigneeID:  assignee.Id,
		CreatorID:   creator.Id,
		Assignee: &users.UserModel{
			ID:        assignee.Id,
			UUID:      assignee.Uuid,
			Username:  assignee.Username,
			Password:  assignee.Password,
			Email:     assignee.Email,
			Enable:    assignee.Enable,
			CreatedAt: *assignee.CreatedAt,
			UpdatedAt: *assignee.UpdatedAt,
		},
		Creator: &users.UserModel{
			ID:        creator.Id,
			UUID:      creator.Uuid,
			Username:  creator.Username,
			Password:  creator.Password,
			Email:     creator.Email,
			Enable:    creator.Enable,
			CreatedAt: *creator.CreatedAt,
			UpdatedAt: *creator.UpdatedAt,
		},
		CreatedAt: *task.CreatedAt,
		UpdatedAt: now,
	}

	if err := s.repo.Update(ctx, taskModel); err != nil {
		return task, err
	}
	return s.Get(ctx, req.Uuid)
}

// Delete deletes the task with the specified UUID.
func (s service) Delete(ctx context.Context, UUID string) (*tasksProto.Task, error) {
	task, err := s.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	if err = s.repo.Delete(ctx, UUID); err != nil {
		return nil, err
	}
	return task, nil
}

// Count returns the number of tasks.
func (s service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// Query returns the tasks with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int64) (*tasksProto.ListTasksResponse, error) {
	items, count, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	var result []*tasksProto.Task
	for _, item := range items {
		result = append(result, &tasksProto.Task{
			Uuid:        item.UUID,
			Title:       item.Title,
			Description: item.Description,
			Status:      item.Status,
			SprintId:    item.SprintID,
			Estimate:    item.Estimate,
			Assignee: &usersProto.User{
				Uuid:      item.Assignee.UUID,
				Username:  item.Assignee.Username,
				Email:     item.Assignee.Email,
				Enable:    item.Assignee.Enable,
				CreatedAt: &item.Assignee.CreatedAt,
				UpdatedAt: &item.Assignee.UpdatedAt,
			},
			Creator: &usersProto.User{
				Uuid:      item.Creator.UUID,
				Username:  item.Creator.Username,
				Email:     item.Creator.Email,
				Enable:    item.Creator.Enable,
				CreatedAt: &item.Creator.CreatedAt,
				UpdatedAt: &item.Creator.UpdatedAt,
			},
			CreatedAt: &item.CreatedAt,
			UpdatedAt: &item.UpdatedAt,
		})
	}
	return &tasksProto.ListTasksResponse{
		Tasks:      result,
		TotalCount: int64(count),
		Offset:     offset,
		Limit:      limit,
	}, nil
}
