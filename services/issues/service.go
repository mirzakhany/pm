package issues

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	issuesProto "github.com/mirzakhany/pm/services/issues/proto"
	"github.com/mirzakhany/pm/services/users"
)

// Service encapsulates use case logic for issues.
type Service interface {
	Get(ctx context.Context, uuid string) (*issuesProto.Issue, error)
	Query(ctx context.Context, offset, limit int64) (*issuesProto.ListIssuesResponse, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, input *issuesProto.CreateIssueRequest) (*issuesProto.Issue, error)
	Update(ctx context.Context, input *issuesProto.UpdateIssueRequest) (*issuesProto.Issue, error)
	Delete(ctx context.Context, uuid string) (*issuesProto.Issue, error)
}

// ValidateCreateRequest validates the CreateIssueRequest fields.
func ValidateCreateRequest(c *issuesProto.CreateIssueRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Title, validation.Required, validation.Length(0, 128)),
		validation.Field(&c.Description, validation.Required, validation.Length(0, 1000)),
	)
}

// Validate validates the UpdateIssueRequest fields.
func ValidateUpdateRequest(u *issuesProto.UpdateIssueRequest) error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Title, validation.Required, validation.Length(0, 128)),
		validation.Field(&u.Description, validation.Required, validation.Length(0, 1000)),
	)
}

type service struct {
	repo    Repository
	userSrv users.Service
}

// NewService creates a new issue service.
func NewService(repo Repository, userSrv users.Service) Service {
	return service{repo, userSrv}
}

// Get returns the issue with the specified the issue UUID.
func (s service) Get(ctx context.Context, UUID string) (*issuesProto.Issue, error) {
	issue, err := s.repo.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	return issue.ToProto(true), nil
}

// Create creates a new issue.
func (s service) Create(ctx context.Context, req *issuesProto.CreateIssueRequest) (*issuesProto.Issue, error) {
	if err := ValidateCreateRequest(req); err != nil {
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

	assigneeModel := users.FromProto(assignee)
	creatorModel := users.FromProto(creator)

	err = s.repo.Create(ctx, IssueModel{
		UUID:        id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		SprintID:    req.SprintId,
		Estimate:    req.Estimate,
		AssigneeID:  assignee.Id,
		CreatorID:   creator.Id,
		Assignee:    &assigneeModel,
		Creator:     &creatorModel,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return nil, err
	}
	return s.Get(ctx, id)
}

// Update updates the issue with the specified UUID.
func (s service) Update(ctx context.Context, req *issuesProto.UpdateIssueRequest) (*issuesProto.Issue, error) {
	if err := ValidateUpdateRequest(req); err != nil {
		return nil, err
	}

	issue, err := s.repo.Get(ctx, req.Uuid)
	if err != nil {
		return nil, err
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

	assigneeModel := users.FromProto(assignee)
	creatorModel := users.FromProto(creator)

	issueModel := IssueModel{
		ID:          issue.ID,
		UUID:        issue.UUID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		SprintID:    req.SprintId,
		Estimate:    req.Estimate,
		AssigneeID:  assignee.Id,
		CreatorID:   creator.Id,
		Assignee:    &assigneeModel,
		Creator:     &creatorModel,
		CreatedAt:   issue.CreatedAt,
		UpdatedAt:   now,
	}

	if err := s.repo.Update(ctx, issueModel); err != nil {
		return nil, err
	}
	return s.Get(ctx, req.Uuid)
}

// Delete deletes the issue with the specified UUID.
func (s service) Delete(ctx context.Context, UUID string) (*issuesProto.Issue, error) {
	issue, err := s.Get(ctx, UUID)
	if err != nil {
		return nil, err
	}
	if err = s.repo.Delete(ctx, UUID); err != nil {
		return nil, err
	}
	return issue, nil
}

// Count returns the number of issues.
func (s service) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}

// Query returns the issues with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int64) (*issuesProto.ListIssuesResponse, error) {
	items, count, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return &issuesProto.ListIssuesResponse{
		Issues:     ToProtoList(items, true),
		TotalCount: int64(count),
		Offset:     offset,
		Limit:      limit,
	}, nil
}
