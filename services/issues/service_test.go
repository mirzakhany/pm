package issues

import (
	"context"
	"errors"
	"testing"

	"github.com/go-pg/pg/v10"
	issues "github.com/mirzakhany/pm/services/issues/proto"
	userSrv "github.com/mirzakhany/pm/services/users"
	usersProto "github.com/mirzakhany/pm/services/users/proto"
	"github.com/stretchr/testify/assert"
)

var errCRUD = errors.New("error crud")

func TestCreateIssueRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     issues.CreateIssueRequest
		wantError bool
	}{
		{"success", issues.CreateIssueRequest{
			Title:       "test",
			Description: "this is a test",
			Status:      0,
			SprintId:    0,
			Estimate:    0,
		}, false},
		{"required", issues.CreateIssueRequest{
			Title:       "",
			Description: "this is a test",
			Status:      0,
			SprintId:    0,
			Estimate:    0,
		}, true},
		{"too long", issues.CreateIssueRequest{
			Description: "this is a test",
			Status:      0,
			SprintId:    0,
			Estimate:    0,
			Title:       "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateIssueRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     issues.UpdateIssueRequest
		wantError bool
	}{
		{"success", issues.UpdateIssueRequest{
			Title:       "test",
			Description: "this is a test",
			Status:      0,
			SprintId:    0,
			Estimate:    0,
		}, false},
		{"required", issues.UpdateIssueRequest{
			Title:       "",
			Description: "this is a test",
			Status:      0,
			SprintId:    0,
			Estimate:    0,
		}, true},
		{"too long", issues.UpdateIssueRequest{
			Description: "this is a test",
			Status:      0,
			SprintId:    0,
			Estimate:    0,
			Title:       "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpdateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func Test_service_CRUD(t *testing.T) {
	userServices := userSrv.NewServiceForTest()
	s := NewService(&mockRepository{}, userServices)
	ctx := context.Background()

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, int64(0), count)

	user1, err := userServices.Create(ctx, &usersProto.CreateUserRequest{
		Username: "test", Password: "test", Email: "test@test.com", Enable: true,
	})
	assert.Nil(t, err)

	user2, err := userServices.Create(ctx, &usersProto.CreateUserRequest{
		Username: "test", Password: "test", Email: "test@test.com", Enable: true,
	})
	assert.Nil(t, err)
	// successful creation
	issue, err := s.Create(ctx, &issues.CreateIssueRequest{
		Title:        "test",
		Description:  "this is a test",
		CreatorUuid:  user1.Uuid,
		AssigneeUuid: user2.Uuid,
		Status:       0,
		SprintId:     0,
		Estimate:     0,
	})

	assert.Nil(t, err)
	assert.NotEmpty(t, issue.Uuid)
	id := issue.Uuid
	assert.Equal(t, "test", issue.Title)
	assert.NotEmpty(t, issue.CreatedAt)
	assert.NotEmpty(t, issue.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// validation error in creation
	_, err = s.Create(ctx, &issues.CreateIssueRequest{
		Title:       "",
		Description: "this is a test",
		Status:      0,
		SprintId:    0,
		Estimate:    0,
	})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// unexpected error in creation
	_, err = s.Create(ctx, &issues.CreateIssueRequest{
		Title:        "error",
		Description:  "this is a test",
		CreatorUuid:  user1.Uuid,
		AssigneeUuid: user2.Uuid,
		Status:       0,
		SprintId:     0,
		Estimate:     0,
	})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	_, _ = s.Create(ctx, &issues.CreateIssueRequest{
		Title:        "test2",
		Description:  "this is a test",
		CreatorUuid:  user1.Uuid,
		AssigneeUuid: user2.Uuid,
		Status:       0,
		SprintId:     0,
		Estimate:     0,
	})

	// update
	issue, err = s.Update(ctx, &issues.UpdateIssueRequest{
		Uuid:         id,
		Title:        "test-updated",
		Description:  "this is a test",
		CreatorUuid:  user1.Uuid,
		AssigneeUuid: user2.Uuid,
		Status:       0,
		SprintId:     0,
		Estimate:     0,
	})
	assert.Nil(t, err)
	assert.Equal(t, "test-updated", issue.Title)
	_, err = s.Update(ctx, &issues.UpdateIssueRequest{
		Uuid:         "none",
		Title:        "test-updated",
		Description:  "this is a test",
		CreatorUuid:  user1.Uuid,
		AssigneeUuid: user2.Uuid,
		Status:       0,
		SprintId:     0,
		Estimate:     0,
	})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, &issues.UpdateIssueRequest{
		Uuid:         id,
		Title:        "",
		Description:  "this is a test",
		CreatorUuid:  user1.Uuid,
		AssigneeUuid: user2.Uuid,
		Status:       0,
		SprintId:     0,
		Estimate:     0,
	})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// unexpected error in update
	_, err = s.Update(ctx, &issues.UpdateIssueRequest{
		Uuid:         id,
		Title:        "error",
		Description:  "this is a test",
		CreatorUuid:  user1.Uuid,
		AssigneeUuid: user2.Uuid,
		Status:       0,
		SprintId:     0,
		Estimate:     0,
	})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	issue, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test-updated", issue.Title)
	assert.Equal(t, id, issue.Uuid)

	// query
	_issues, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, int(_issues.TotalCount))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	issue, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, issue.Uuid)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)
}

type mockRepository struct {
	items []IssueModel
}

func (m mockRepository) Get(ctx context.Context, id string) (IssueModel, error) {
	for _, item := range m.items {
		if item.UUID == id {
			return item, nil
		}
	}
	return IssueModel{}, pg.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int64) ([]IssueModel, int, error) {
	return m.items, len(m.items), nil
}

func (m *mockRepository) Create(ctx context.Context, issue IssueModel) error {
	if issue.Title == "error" {
		return errCRUD
	}
	m.items = append(m.items, issue)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, issue IssueModel) error {
	if issue.Title == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.UUID == issue.UUID {
			m.items[i] = issue
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	for i, item := range m.items {
		if item.UUID == id {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}
