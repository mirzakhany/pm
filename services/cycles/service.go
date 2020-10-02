package cycles

import (
	"context"
	"time"

	"github.com/mirzakhany/pm/pkg/auth"

	"github.com/mirzakhany/pm/services/users"

	usersProto "github.com/mirzakhany/pm/services/users/proto"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	cyclesProto "github.com/mirzakhany/pm/services/cycles/proto"
)

// Service encapsulates use case logic for cycles.
type Service interface {
	Get(ctx context.Context, uuid string) (*cyclesProto.Cycle, error)
	Query(ctx context.Context, offset, limit int64) (*cyclesProto.ListCyclesResponse, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, input *cyclesProto.CreateCycleRequest) (*cyclesProto.Cycle, error)
	Update(ctx context.Context, input *cyclesProto.UpdateCycleRequest) (*cyclesProto.Cycle, error)
	Delete(ctx context.Context, uuid string) (*cyclesProto.Cycle, error)
}

// ValidateCreateRequest validates the CreateCycleRequest fields.
func ValidateCreateRequest(c cyclesProto.CreateCycleRequest) error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Title, validation.Required, validation.Length(0, 128)),
	)
}

// Validate validates the UpdateCycleRequest fields.
func ValidateUpdateRequest(u cyclesProto.UpdateCycleRequest) error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Title, validation.Required, validation.Length(0, 128)),
	)
}

type service struct {
	repo    Repository
	userSrv users.Service
}

// NewService creates a new cycle service.
func NewService(repo Repository, userSrv users.Service) Service {
	return service{repo, userSrv}
}

// Get returns the cycle with the specified the cycle UUID.
func (s service) Get(ctx context.Context, UUID string) (*cyclesProto.Cycle, error) {
	cycle, err := s.repo.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	return &cyclesProto.Cycle{
		Id:          cycle.ID,
		Uuid:        cycle.UUID,
		Title:       cycle.Title,
		Description: cycle.Description,
		Active:      cycle.Active,
		StartAt:     &cycle.StartAt,
		EndAt:       &cycle.StartAt,
		Creator: &usersProto.User{
			Uuid:      cycle.Creator.UUID,
			Username:  cycle.Creator.Username,
			Email:     cycle.Creator.Email,
			Enable:    cycle.Creator.Enable,
			CreatedAt: &cycle.Creator.CreatedAt,
			UpdatedAt: &cycle.Creator.UpdatedAt,
		},
		CreatedAt: &cycle.CreatedAt,
		UpdatedAt: &cycle.UpdatedAt,
	}, nil
}

// Create creates a new cycle.
func (s service) Create(ctx context.Context, req *cyclesProto.CreateCycleRequest) (*cyclesProto.Cycle, error) {
	if err := ValidateCreateRequest(*req); err != nil {
		return nil, err
	}

	user, err := auth.ExtractUser(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	id := uuid.New().String()
	err = s.repo.Create(ctx, CycleModel{
		UUID:        id,
		Title:       req.Title,
		Description: req.Description,
		Active:      req.Active,
		StartAt:     *req.StartAt,
		EndAt:       *req.EndAt,
		Creator: &users.UserModel{
			ID:        user.Id,
			UUID:      user.Uuid,
			Username:  user.Username,
			Password:  user.Password,
			Email:     user.Email,
			Enable:    user.Enable,
			CreatedAt: *user.CreatedAt,
			UpdatedAt: *user.UpdatedAt,
		},
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// Update updates the cycle with the specified UUID.
func (s service) Update(ctx context.Context, req *cyclesProto.UpdateCycleRequest) (*cyclesProto.Cycle, error) {
	if err := ValidateUpdateRequest(*req); err != nil {
		return nil, err
	}

	cycle, err := s.Get(ctx, req.Uuid)
	if err != nil {
		return cycle, err
	}
	now := time.Now()

	cycle.Title = req.Title
	cycle.UpdatedAt = &now

	cycleModel := CycleModel{
		ID:          cycle.Id,
		UUID:        cycle.Uuid,
		Title:       req.Title,
		Description: req.Description,
		Active:      req.Active,
		StartAt:     *req.StartAt,
		EndAt:       *req.EndAt,
		Creator: &users.UserModel{
			ID:        cycle.Creator.Id,
			UUID:      cycle.Creator.Uuid,
			Username:  cycle.Creator.Username,
			Password:  cycle.Creator.Password,
			Email:     cycle.Creator.Email,
			Enable:    cycle.Creator.Enable,
			CreatedAt: *cycle.Creator.CreatedAt,
			UpdatedAt: *cycle.Creator.UpdatedAt,
		},
		CreatedAt: *cycle.CreatedAt,
		UpdatedAt: now,
	}

	if err := s.repo.Update(ctx, cycleModel); err != nil {
		return cycle, err
	}
	return cycle, nil
}

// Delete deletes the cycle with the specified UUID.
func (s service) Delete(ctx context.Context, UUID string) (*cyclesProto.Cycle, error) {
	cycle, err := s.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	if err = s.repo.Delete(ctx, UUID); err != nil {
		return nil, err
	}
	return cycle, nil
}

// Count returns the number of cycles.
func (s service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// Query returns the cycles with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int64) (*cyclesProto.ListCyclesResponse, error) {
	items, count, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	var result []*cyclesProto.Cycle
	for _, item := range items {
		result = append(result, &cyclesProto.Cycle{
			Id:          item.ID,
			Uuid:        item.UUID,
			Title:       item.Title,
			Description: item.Description,
			Active:      item.Active,
			StartAt:     &item.StartAt,
			EndAt:       &item.StartAt,
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
	return &cyclesProto.ListCyclesResponse{
		Cycles:     result,
		TotalCount: int64(count),
		Offset:     offset,
		Limit:      limit,
	}, nil
}
