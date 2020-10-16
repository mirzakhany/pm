package templates

const ServiceTmpl = `
package {{ .Pkg.NamePlural | lower }}

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/mirzakhany/pm/internal/entity"
	{{ .Pkg.NamePlural | lower }}Proto "github.com/mirzakhany/pm/protobuf/{{ .Pkg.NamePlural | lower }}"
)

// Service encapsulates use case logic for {{ .Pkg.NamePlural | lower }}.
type Service interface {
	Get(ctx context.Context, uuid string) (*{{ .Pkg.NamePlural | lower }}Proto.{{ .Pkg.Name }}, error)
	Query(ctx context.Context, offset, limit int64) (*{{ .Pkg.NamePlural | lower }}Proto.List{{ .Pkg.NamePlural }}Response, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, input *{{ .Pkg.NamePlural | lower }}Proto.Create{{ .Pkg.Name }}Request) (*{{ .Pkg.NamePlural | lower }}Proto.{{ .Pkg.Name }}, error)
	Update(ctx context.Context, input *{{ .Pkg.NamePlural | lower }}Proto.Update{{ .Pkg.Name }}Request) (*{{ .Pkg.NamePlural | lower }}Proto.{{ .Pkg.Name }}, error)
	Delete(ctx context.Context, uuid string) (*{{ .Pkg.NamePlural | lower }}Proto.{{ .Pkg.Name }}, error)
}

// ValidateCreateRequest validates the Create{{ .Pkg.Name }}Request fields.
func ValidateCreateRequest(c *{{ .Pkg.NamePlural | lower }}Proto.Create{{ .Pkg.Name }}Request) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Title, validation.Required, validation.Length(0, 128)),
	)
}

// Validate validates the Update{{ .Pkg.Name }}Request fields.
func ValidateUpdateRequest(u *{{ .Pkg.NamePlural | lower }}Proto.Update{{ .Pkg.Name }}Request) error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Title, validation.Required, validation.Length(0, 128)),
	)
}

type service struct {
	repo Repository
}

// NewService creates a new {{ .Pkg.Name | lower }} service.
func NewService(repo Repository) Service {
	return service{repo}
}

// Get returns the {{ .Pkg.Name | lower }} with the specified the {{ .Pkg.Name | lower }} UUID.
func (s service) Get(ctx context.Context, UUID string) (*{{ .Pkg.NamePlural | lower }}Proto.{{ .Pkg.Name }}, error) {
	{{ .Pkg.Name | lower }}, err := s.repo.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	return {{ .Pkg.Name | lower }}.ToProto(), nil
}

// Create creates a new {{ .Pkg.Name | lower }}.
func (s service) Create(ctx context.Context, req *{{ .Pkg.NamePlural | lower }}Proto.Create{{ .Pkg.Name }}Request) (*{{ .Pkg.NamePlural | lower }}Proto.{{ .Pkg.Name }}, error) {
	if err := ValidateCreateRequest(req); err != nil {
		return nil, err
	}
	now := time.Now()
	id := uuid.New().String()
	err := s.repo.Create(ctx, entity.{{ .Pkg.Name }}{
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

// Update updates the {{ .Pkg.Name | lower }} with the specified UUID.
func (s service) Update(ctx context.Context, req *{{ .Pkg.NamePlural | lower }}Proto.Update{{ .Pkg.Name }}Request) (*{{ .Pkg.NamePlural | lower }}Proto.{{ .Pkg.Name }}, error) {
	if err := ValidateUpdateRequest(req); err != nil {
		return nil, err
	}

	{{ .Pkg.Name | lower }}, err := s.repo.Get(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	{{ .Pkg.Name | lower }}.Title = req.Title
	{{ .Pkg.Name | lower }}.UpdatedAt = now

	{{ .Pkg.Name | lower }}Model := entity.{{ .Pkg.Name }}{
		ID:        {{ .Pkg.Name | lower }}.ID,
		UUID:      {{ .Pkg.Name | lower }}.UUID,
		Title:     req.Title,
		CreatedAt: {{ .Pkg.Name | lower }}.CreatedAt,
		UpdatedAt: now,
	}

	if err := s.repo.Update(ctx, {{ .Pkg.Name | lower }}Model); err != nil {
		return nil, err
	}
	return {{ .Pkg.Name | lower }}.ToProto(), nil
}

// Delete deletes the {{ .Pkg.Name | lower }} with the specified UUID.
func (s service) Delete(ctx context.Context, UUID string) (*{{ .Pkg.NamePlural | lower }}Proto.{{ .Pkg.Name }}, error) {
	{{ .Pkg.Name | lower }}, err := s.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	if err = s.repo.Delete(ctx, UUID); err != nil {
		return nil, err
	}
	return {{ .Pkg.Name | lower }}, nil
}

// Count returns the number of {{ .Pkg.NamePlural | lower }}.
func (s service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// Query returns the {{ .Pkg.NamePlural | lower }} with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int64) (*{{ .Pkg.NamePlural | lower }}Proto.List{{ .Pkg.NamePlural }}Response, error) {
	items, count, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return &{{ .Pkg.NamePlural | lower }}Proto.List{{ .Pkg.NamePlural }}Response{
		{{ .Pkg.NamePlural }}:      entity.{{ .Pkg.Name }}ToProtoList(items),
		TotalCount: int64(count),
		Offset:     offset,
		Limit:      limit,
	}, nil
}

`
