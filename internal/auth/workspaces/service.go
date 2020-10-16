package workspaces

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/mirzakhany/pm/internal/entity"
	workspacesProto "github.com/mirzakhany/pm/protobuf/workspaces"
)

// Service encapsulates use case logic for workspaces.
type Service interface {
	Get(ctx context.Context, uuid string) (*workspacesProto.Workspace, error)
	Query(ctx context.Context, offset, limit int64) (*workspacesProto.ListWorkspacesResponse, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, input *workspacesProto.CreateWorkspaceRequest) (*workspacesProto.Workspace, error)
	Update(ctx context.Context, input *workspacesProto.UpdateWorkspaceRequest) (*workspacesProto.Workspace, error)
	Delete(ctx context.Context, uuid string) (*workspacesProto.Workspace, error)
}

// ValidateCreateRequest validates the CreateWorkspaceRequest fields.
func ValidateCreateRequest(c *workspacesProto.CreateWorkspaceRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Title, validation.Required, validation.Length(0, 128)),
		validation.Field(&c.Domain, validation.Required, validation.Length(0, 128)),
	)
}

// Validate validates the UpdateWorkspaceRequest fields.
func ValidateUpdateRequest(u *workspacesProto.UpdateWorkspaceRequest) error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Title, validation.Required, validation.Length(0, 128)),
		validation.Field(&u.Domain, validation.Required, validation.Length(0, 128)),
	)
}

type service struct {
	repo Repository
}

// NewService creates a new workspace service.
func NewService(repo Repository) Service {
	return service{repo}
}

// Get returns the workspace with the specified the workspace UUID.
func (s service) Get(ctx context.Context, UUID string) (*workspacesProto.Workspace, error) {
	workspace, err := s.repo.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	return workspace.ToProto(), nil
}

// Create creates a new workspace.
func (s service) Create(ctx context.Context, req *workspacesProto.CreateWorkspaceRequest) (*workspacesProto.Workspace, error) {
	if err := ValidateCreateRequest(req); err != nil {
		return nil, err
	}
	now := time.Now()
	id := uuid.New().String()
	err := s.repo.Create(ctx, entity.Workspace{
		UUID:      id,
		Title:     req.Title,
		Domain:    req.Domain,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// Update updates the workspace with the specified UUID.
func (s service) Update(ctx context.Context, req *workspacesProto.UpdateWorkspaceRequest) (*workspacesProto.Workspace, error) {
	if err := ValidateUpdateRequest(req); err != nil {
		return nil, err
	}

	workspace, err := s.repo.Get(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	workspace.Title = req.Title
	workspace.Domain = req.Domain
	workspace.UpdatedAt = now

	workspaceModel := entity.Workspace{
		ID:        workspace.ID,
		UUID:      workspace.UUID,
		Title:     req.Title,
		Domain:    req.Domain,
		CreatedAt: workspace.CreatedAt,
		UpdatedAt: now,
	}

	if err := s.repo.Update(ctx, workspaceModel); err != nil {
		return nil, err
	}
	return workspace.ToProto(), nil
}

// Delete deletes the workspace with the specified UUID.
func (s service) Delete(ctx context.Context, UUID string) (*workspacesProto.Workspace, error) {
	workspace, err := s.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	if err = s.repo.Delete(ctx, UUID); err != nil {
		return nil, err
	}
	return workspace, nil
}

// Count returns the number of workspaces.
func (s service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// Query returns the workspaces with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int64) (*workspacesProto.ListWorkspacesResponse, error) {
	items, count, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return &workspacesProto.ListWorkspacesResponse{
		Workspaces: entity.WorkspaceToProtoList(items),
		TotalCount: int64(count),
		Offset:     offset,
		Limit:      limit,
	}, nil
}
