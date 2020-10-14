package issues

import (
	"context"
	"time"

	"github.com/mirzakhany/pm/internal/entity"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"github.com/mirzakhany/pm/internal/cycles"
	"github.com/mirzakhany/pm/internal/users"
	issuesProto "github.com/mirzakhany/pm/protobuf/issues"
)

// Service encapsulates use case logic for issues.
type Service interface {
	Get(ctx context.Context, uuid string) (*issuesProto.Issue, error)
	Query(ctx context.Context, offset, limit int64) (*issuesProto.ListIssuesResponse, error)
	Count(ctx context.Context) (int64, error)
	Create(ctx context.Context, input *issuesProto.CreateIssueRequest) (*issuesProto.Issue, error)
	Update(ctx context.Context, input *issuesProto.UpdateIssueRequest) (*issuesProto.Issue, error)
	Delete(ctx context.Context, uuid string) (*issuesProto.Issue, error)

	GetStatus(ctx context.Context, Uuid string) (*issuesProto.IssueStatus, error)
	QueryStatus(ctx context.Context, offset, limit int64) (*issuesProto.ListIssueStatusResponse, error)
	CountStatus(ctx context.Context) (int64, error)
	CreateStatus(ctx context.Context, input *issuesProto.CreateIssueStatusRequest) (*issuesProto.IssueStatus, error)
	UpdateStatus(ctx context.Context, input *issuesProto.UpdateIssueStatusRequest) (*issuesProto.IssueStatus, error)
	DeleteStatus(ctx context.Context, Uuid string) (*issuesProto.IssueStatus, error)
}

// ValidateCreateRequest validates the CreateIssueRequest fields.
func ValidateCreateRequest(c *issuesProto.CreateIssueRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Title, validation.Required, validation.Length(0, 128)),
		validation.Field(&c.Description, validation.Required, validation.Length(0, 1000)),
	)
}

// ValidateUpdateRequest validates the UpdateIssueRequest fields.
func ValidateUpdateRequest(u *issuesProto.UpdateIssueRequest) error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Title, validation.Required, validation.Length(0, 128)),
		validation.Field(&u.Description, validation.Required, validation.Length(0, 1000)),
	)
}

// ValidateStatusCreateRequest validates the CreateIssueStatusRequest fields.
func ValidateStatusCreateRequest(c *issuesProto.CreateIssueStatusRequest) error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Title, validation.Required, validation.Length(0, 128)),
	)
}

// ValidateStatusUpdateRequest validates the UpdateIssueStatusRequest fields.
func ValidateStatusUpdateRequest(u *issuesProto.UpdateIssueStatusRequest) error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Title, validation.Required, validation.Length(0, 128)),
	)
}

type service struct {
	repo      Repository
	usersSrv  users.Service
	cyclesSrv cycles.Service
}

// NewService creates a new issue service.
func NewService(repo Repository, userSrv users.Service, cyclesSrv cycles.Service) Service {
	return service{repo, userSrv, cyclesSrv}
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

	creator, err := s.usersSrv.Get(ctx, req.CreatorUuid)
	if err != nil {
		return nil, err
	}

	assignee, err := s.usersSrv.Get(ctx, req.AssigneeUuid)
	if err != nil {
		return nil, err
	}

	cycle, err := s.cyclesSrv.Get(ctx, req.CycleUuid)
	if err != nil {
		return nil, err
	}

	status, err := s.repo.GetStatus(ctx, req.StatusUuid)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	id := uuid.New().String()

	assigneeModel := entity.UserFromProto(assignee)
	creatorModel := entity.UserFromProto(creator)
	cycleModel := entity.CycleFromProto(cycle)

	err = s.repo.Create(ctx, entity.Issue{
		UUID:        id,
		Title:       req.Title,
		Description: req.Description,
		Status:      &status,
		Cycle:       &cycleModel,
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

	creator, err := s.usersSrv.Get(ctx, req.CreatorUuid)
	if err != nil {
		return nil, err
	}

	assignee, err := s.usersSrv.Get(ctx, req.AssigneeUuid)
	if err != nil {
		return nil, err
	}

	cycle, err := s.cyclesSrv.Get(ctx, req.CycleUuid)
	if err != nil {
		return nil, err
	}

	status, err := s.repo.GetStatus(ctx, req.StatusUuid)
	if err != nil {
		return nil, err
	}

	assigneeModel := entity.UserFromProto(assignee)
	creatorModel := entity.UserFromProto(creator)
	cycleModel := entity.CycleFromProto(cycle)

	issueModel := entity.Issue{
		ID:          issue.ID,
		UUID:        issue.UUID,
		Title:       req.Title,
		Description: req.Description,
		Status:      &status,
		Cycle:       &cycleModel,
		CycleID:     cycle.Id,
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
		Issues:     entity.IssueToProtoList(items, true),
		TotalCount: int64(count),
		Offset:     offset,
		Limit:      limit,
	}, nil
}

func (s service) GetStatus(ctx context.Context, Uuid string) (*issuesProto.IssueStatus, error) {
	issueStatus, err := s.repo.GetStatus(ctx, Uuid)
	if err != nil {
		return nil, err
	}
	return issueStatus.ToProto(true), nil
}

func (s service) QueryStatus(ctx context.Context, offset, limit int64) (*issuesProto.ListIssueStatusResponse, error) {
	items, count, err := s.repo.QueryStatus(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	return &issuesProto.ListIssueStatusResponse{
		IssueStatus: entity.IssueStatusToProtoList(items, true),
		TotalCount:  int64(count),
		Offset:      offset,
		Limit:       limit,
	}, nil
}

func (s service) CountStatus(ctx context.Context) (int64, error) {
	return s.repo.CountStatus(ctx)
}

func (s service) CreateStatus(ctx context.Context, req *issuesProto.CreateIssueStatusRequest) (*issuesProto.IssueStatus, error) {
	if err := ValidateStatusCreateRequest(req); err != nil {
		return nil, err
	}

	id := uuid.New().String()
	now := time.Now()
	err := s.repo.CreateStatus(ctx, entity.IssueStatus{
		UUID:      id,
		Title:     req.Title,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, err
	}
	return s.GetStatus(ctx, id)
}

func (s service) UpdateStatus(ctx context.Context, req *issuesProto.UpdateIssueStatusRequest) (*issuesProto.IssueStatus, error) {
	if err := ValidateStatusUpdateRequest(req); err != nil {
		return nil, err
	}

	issueStatus, err := s.repo.GetStatus(ctx, req.Uuid)
	if err != nil {
		return nil, err
	}
	now := time.Now()

	issueStatusModel := entity.IssueStatus{
		ID:        issueStatus.ID,
		UUID:      issueStatus.UUID,
		Title:     req.Title,
		CreatedAt: issueStatus.CreatedAt,
		UpdatedAt: now,
	}

	if err := s.repo.UpdateStatus(ctx, issueStatusModel); err != nil {
		return nil, err
	}
	return s.GetStatus(ctx, req.Uuid)
}

func (s service) DeleteStatus(ctx context.Context, Uuid string) (*issuesProto.IssueStatus, error) {
	issueStatus, err := s.GetStatus(ctx, Uuid)
	if err != nil {
		return nil, err
	}
	if err = s.repo.DeleteStatus(ctx, Uuid); err != nil {
		return nil, err
	}
	return issueStatus, nil
}
