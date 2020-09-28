package issues

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	issuesProto "github.com/mirzakhany/pm/services/issues/proto"
	"github.com/mirzakhany/pm/services/users"
	usersProto "github.com/mirzakhany/pm/services/users/proto"
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
func ValidateCreateRequest(c issuesProto.CreateIssueRequest) error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Title, validation.Required, validation.Length(0, 128)),
		validation.Field(&c.Description, validation.Required, validation.Length(0, 1000)),
	)
}

// Validate validates the UpdateIssueRequest fields.
func ValidateUpdateRequest(u issuesProto.UpdateIssueRequest) error {
	return validation.ValidateStruct(&u,
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
	return &issuesProto.Issue{
		Uuid:        issue.UUID,
		Title:       issue.Title,
		Description: issue.Description,
		Status:      issue.Status,
		SprintId:    issue.SprintID,
		Assignee: &usersProto.User{
			Uuid:      issue.Assignee.UUID,
			Username:  issue.Assignee.Username,
			Email:     issue.Assignee.Email,
			Enable:    issue.Assignee.Enable,
			CreatedAt: &issue.Assignee.CreatedAt,
			UpdatedAt: &issue.Assignee.UpdatedAt,
		},
		Creator: &usersProto.User{
			Uuid:      issue.Creator.UUID,
			Username:  issue.Creator.Username,
			Email:     issue.Creator.Email,
			Enable:    issue.Creator.Enable,
			CreatedAt: &issue.Creator.CreatedAt,
			UpdatedAt: &issue.Creator.UpdatedAt,
		},
		CreatedAt: &issue.CreatedAt,
		UpdatedAt: &issue.UpdatedAt,
	}, nil
}

// Create creates a new issue.
func (s service) Create(ctx context.Context, req *issuesProto.CreateIssueRequest) (*issuesProto.Issue, error) {
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
	err = s.repo.Create(ctx, IssueModel{
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

// Update updates the issue with the specified UUID.
func (s service) Update(ctx context.Context, req *issuesProto.UpdateIssueRequest) (*issuesProto.Issue, error) {
	if err := ValidateUpdateRequest(*req); err != nil {
		return nil, err
	}

	issue, err := s.Get(ctx, req.Uuid)
	if err != nil {
		return issue, err
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

	issueModel := IssueModel{
		ID:          issue.Id,
		UUID:        issue.Uuid,
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
		CreatedAt: *issue.CreatedAt,
		UpdatedAt: now,
	}

	if err := s.repo.Update(ctx, issueModel); err != nil {
		return issue, err
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
	var result []*issuesProto.Issue
	for _, item := range items {
		result = append(result, &issuesProto.Issue{
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
	return &issuesProto.ListIssuesResponse{
		Issues:     result,
		TotalCount: int64(count),
		Offset:     offset,
		Limit:      limit,
	}, nil
}
