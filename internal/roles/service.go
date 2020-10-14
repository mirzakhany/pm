package roles

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/mirzakhany/pm/internal/entity"
	rolesProto "github.com/mirzakhany/pm/protobuf/roles"
)

// Service encapsulates use case logic for roles.
type Service interface {
	Get(ctx context.Context, uuid string) (*rolesProto.Role, error)
	Query(ctx context.Context, offset, limit int64) (*rolesProto.ListRolesResponse, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, input *rolesProto.CreateRoleRequest) (*rolesProto.Role, error)
	Update(ctx context.Context, input *rolesProto.UpdateRoleRequest) (*rolesProto.Role, error)
	Delete(ctx context.Context, uuid string) (*rolesProto.Role, error)
}

// ValidateCreateRequest validates the CreateRoleRequest fields.
func ValidateCreateRequest(c *rolesProto.CreateRoleRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Title, validation.Required, validation.Length(0, 128)),
	)
}

// Validate validates the UpdateRoleRequest fields.
func ValidateUpdateRequest(u *rolesProto.UpdateRoleRequest) error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Title, validation.Required, validation.Length(0, 128)),
	)
}

type service struct {
	repo Repository
}

// NewService creates a new role service.
func NewService(repo Repository) Service {
	return service{repo}
}

// Get returns the role with the specified the role UUID.
func (s service) Get(ctx context.Context, UUID string) (*rolesProto.Role, error) {
	role, err := s.repo.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	return role.ToProto(), nil
}

// Create creates a new role.
func (s service) Create(ctx context.Context, req *rolesProto.CreateRoleRequest) (*rolesProto.Role, error) {
	if err := ValidateCreateRequest(req); err != nil {
		return nil, err
	}
	now := time.Now()
	id := uuid.New().String()
	err := s.repo.Create(ctx, entity.Role{
		UUID:      id,
		Title:     req.Title,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// Update updates the role with the specified UUID.
func (s service) Update(ctx context.Context, req *rolesProto.UpdateRoleRequest) (*rolesProto.Role, error) {
	if err := ValidateUpdateRequest(req); err != nil {
		return nil, err
	}

	role, err := s.repo.Get(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	role.Title = req.Title
	role.UpdatedAt = now

	roleModel := entity.Role{
		ID:        role.ID,
		UUID:      role.UUID,
		Title:     req.Title,
		CreatedAt: role.CreatedAt,
		UpdatedAt: now,
	}

	if err := s.repo.Update(ctx, roleModel); err != nil {
		return nil, err
	}
	return role.ToProto(), nil
}

// Delete deletes the role with the specified UUID.
func (s service) Delete(ctx context.Context, UUID string) (*rolesProto.Role, error) {
	role, err := s.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	if err = s.repo.Delete(ctx, UUID); err != nil {
		return nil, err
	}
	return role, nil
}

// Count returns the number of roles.
func (s service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// Query returns the roles with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int64) (*rolesProto.ListRolesResponse, error) {
	items, count, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return &rolesProto.ListRolesResponse{
		Roles:      entity.RoleToProtoList(items),
		TotalCount: int64(count),
		Offset:     offset,
		Limit:      limit,
	}, nil
}
