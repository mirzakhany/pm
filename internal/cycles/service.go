package cycles

import (
	"context"
	"time"

	"github.com/mirzakhany/pm/internal/entity"

	"github.com/golang/protobuf/ptypes"

	"github.com/mirzakhany/pm/internal/auth/users"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	cyclesProto "github.com/mirzakhany/pm/protobuf/cycles"
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
func ValidateCreateRequest(c *cyclesProto.CreateCycleRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Title, validation.Required, validation.Length(0, 128)),
	)
}

// Validate validates the UpdateCycleRequest fields.
func ValidateUpdateRequest(u *cyclesProto.UpdateCycleRequest) error {
	return validation.ValidateStruct(u,
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
	return cycle.ToProto(true), nil
}

// Create creates a new cycle.
func (s service) Create(ctx context.Context, req *cyclesProto.CreateCycleRequest) (*cyclesProto.Cycle, error) {
	if err := ValidateCreateRequest(req); err != nil {
		return nil, err
	}

	//user, err := auth.ExtractUser(ctx)
	//if err != nil {
	//	return nil, err
	//}

	startAt, err := ptypes.Timestamp(req.StartAt)
	if err != nil {
		return nil, err
	}

	endAt, err := ptypes.Timestamp(req.EndAt)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	id := uuid.New().String()

	err = s.repo.Create(ctx, entity.Cycle{
		UUID:        id,
		Title:       req.Title,
		Description: req.Description,
		Active:      req.Active,
		StartAt:     startAt,
		EndAt:       endAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// Update updates the cycle with the specified UUID.
func (s service) Update(ctx context.Context, req *cyclesProto.UpdateCycleRequest) (*cyclesProto.Cycle, error) {
	if err := ValidateUpdateRequest(req); err != nil {
		return nil, err
	}

	//user, err := auth.ExtractUser(ctx)
	//if err != nil {
	//	return nil, err
	//}

	cycle, err := s.repo.Get(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	startAt, err := ptypes.Timestamp(req.StartAt)
	if err != nil {
		return nil, err
	}

	endAt, err := ptypes.Timestamp(req.EndAt)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	cycle.Title = req.Title
	cycle.UpdatedAt = now

	cycleModel := entity.Cycle{
		ID:          cycle.ID,
		UUID:        cycle.UUID,
		Title:       req.Title,
		Description: req.Description,
		Active:      req.Active,
		StartAt:     startAt,
		EndAt:       endAt,
		CreatedAt:   cycle.CreatedAt,
		UpdatedAt:   now,
	}

	if err := s.repo.Update(ctx, cycleModel); err != nil {
		return nil, err
	}
	return cycle.ToProto(true), nil
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
	return &cyclesProto.ListCyclesResponse{
		Cycles:     entity.CycleToProtoList(items, true),
		TotalCount: int64(count),
		Offset:     offset,
		Limit:      limit,
	}, nil
}
