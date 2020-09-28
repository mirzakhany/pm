package users

import (
	"context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	usersProto "github.com/mirzakhany/pm/services/users/proto"
	"time"
)

// Service encapsulates use case logic for users.
type Service interface {
	Get(ctx context.Context, uuid string) (*usersProto.User, error)
	Query(ctx context.Context, offset, limit int64) (*usersProto.ListUsersResponse, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, input *usersProto.CreateUserRequest) (*usersProto.User, error)
	Update(ctx context.Context, input *usersProto.UpdateUserRequest) (*usersProto.User, error)
	Delete(ctx context.Context, uuid string) (*usersProto.User, error)
}

// ValidateCreateRequest validates the CreateUserRequest fields.
func ValidateCreateRequest(c usersProto.CreateUserRequest) error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Username, validation.Required, validation.Length(0, 128)),
		validation.Field(&c.Email, validation.Required, validation.Length(0, 128)),
		validation.Field(&c.Password, validation.Required, validation.Length(0, 128)),
	)
}

// Validate validates the UpdateUserRequest fields.
func ValidateUpdateRequest(u usersProto.UpdateUserRequest) error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Username, validation.Required, validation.Length(0, 128)),
		validation.Field(&u.Email, validation.Required, validation.Length(0, 128)),
		validation.Field(&u.Enable, validation.Required),
	)
}

type service struct {
	repo Repository
}

// NewService creates a new user service.
func NewService(repo Repository) Service {
	return service{repo}
}

// NewServiceForTest creates a new user service for test.
func NewServiceForTest() Service {
	return NewService(&mockRepository{})
}

// Get returns the user with the specified the user UUID.
func (s service) Get(ctx context.Context, UUID string) (*usersProto.User, error) {
	user, err := s.repo.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	return &usersProto.User{
		Uuid:      user.UUID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: &user.CreatedAt,
		UpdatedAt: &user.UpdatedAt,
	}, nil
}

// Create creates a new user.
func (s service) Create(ctx context.Context, req *usersProto.CreateUserRequest) (*usersProto.User, error) {
	if err := ValidateCreateRequest(*req); err != nil {
		return nil, err
	}
	now := time.Now()
	id := uuid.New().String()
	err := s.repo.Create(ctx, UserModel{
		UUID:      id,
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		Enable:    req.Enable,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// Update updates the user with the specified UUID.
func (s service) Update(ctx context.Context, req *usersProto.UpdateUserRequest) (*usersProto.User, error) {
	if err := ValidateUpdateRequest(*req); err != nil {
		return nil, err
	}

	user, err := s.Get(ctx, req.Uuid)
	if err != nil {
		return user, err
	}
	now := time.Now()

	user.Id = 0
	user.Username = req.Username
	user.Email = req.Email
	user.Enable = req.Enable
	user.UpdatedAt = &now

	userModel := UserModel{
		ID:        user.Id,
		UUID:      user.Uuid,
		Username:  req.Username,
		Password:  user.Password,
		Email:     req.Email,
		Enable:    req.Enable,
		CreatedAt: *user.CreatedAt,
		UpdatedAt: now,
	}

	if err := s.repo.Update(ctx, userModel); err != nil {
		return user, err
	}
	return user, nil
}

// Delete deletes the user with the specified UUID.
func (s service) Delete(ctx context.Context, UUID string) (*usersProto.User, error) {
	user, err := s.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	if err = s.repo.Delete(ctx, UUID); err != nil {
		return nil, err
	}
	return user, nil
}

// Count returns the number of users.
func (s service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// Query returns the users with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int64) (*usersProto.ListUsersResponse, error) {
	items, count, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	var result []*usersProto.User
	for _, item := range items {
		result = append(result, &usersProto.User{
			Uuid:      item.UUID,
			Username:  item.Username,
			Email:     item.Email,
			Enable:    item.Enable,
			CreatedAt: &item.CreatedAt,
			UpdatedAt: &item.CreatedAt,
		})
	}
	return &usersProto.ListUsersResponse{
		Users:      result,
		TotalCount: int64(count),
		Offset:     offset,
		Limit:      limit,
	}, nil
}
