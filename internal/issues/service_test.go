package issues

import (
	"context"
	"errors"
	"testing"

	"github.com/mirzakhany/pm/internal/entity"

	"github.com/mirzakhany/pm/internal/cycles"

	"github.com/google/uuid"

	"github.com/go-pg/pg/v10"
	userSrv "github.com/mirzakhany/pm/internal/auth/users"
	issues "github.com/mirzakhany/pm/protobuf/issues"
	usersProto "github.com/mirzakhany/pm/protobuf/users"
	"github.com/stretchr/testify/assert"
)

var errCRUD = errors.New("error crud")

func TestCreateIssueRequest_Validate(t *testing.T) {
	Uuid := uuid.New().String()
	tests := []struct {
		name      string
		model     issues.CreateIssueRequest
		wantError bool
	}{
		{"success", issues.CreateIssueRequest{
			Title:       "test",
			Description: "this is a test",
			StatusUuid:  Uuid,
			CycleUuid:   Uuid,
			Estimate:    0,
		}, false},
		{"required", issues.CreateIssueRequest{
			Title:       "",
			Description: "this is a test",
			StatusUuid:  Uuid,
			CycleUuid:   Uuid,
			Estimate:    0,
		}, true},
		{"too long", issues.CreateIssueRequest{
			Description: "this is a test",
			StatusUuid:  Uuid,
			CycleUuid:   Uuid,
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
	Uuid := uuid.New().String()
	tests := []struct {
		name      string
		model     issues.UpdateIssueRequest
		wantError bool
	}{
		{"success", issues.UpdateIssueRequest{
			Title:       "test",
			Description: "this is a test",
			StatusUuid:  Uuid,
			CycleUuid:   Uuid,
			Estimate:    0,
		}, false},
		{"required", issues.UpdateIssueRequest{
			Title:       "",
			Description: "this is a test",
			StatusUuid:  Uuid,
			CycleUuid:   Uuid,
			Estimate:    0,
		}, true},
		{"too long", issues.UpdateIssueRequest{
			Description: "this is a test",
			StatusUuid:  Uuid,
			CycleUuid:   Uuid,
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
	Uuid := uuid.New().String()
	userServices := userSrv.NewServiceForTest()
	cycleService := cycles.NewServiceForTest(userServices)
	s := NewService(&mockRepository{}, userServices, cycleService)
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
		StatusUuid:   Uuid,
		CycleUuid:    Uuid,
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
		StatusUuid:  Uuid,
		CycleUuid:   Uuid,
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
		StatusUuid:   Uuid,
		CycleUuid:    Uuid,
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
		StatusUuid:   Uuid,
		CycleUuid:    Uuid,
		Estimate:     0,
	})

	// update
	issue, err = s.Update(ctx, &issues.UpdateIssueRequest{
		Uuid:         id,
		Title:        "test-updated",
		Description:  "this is a test",
		CreatorUuid:  user1.Uuid,
		AssigneeUuid: user2.Uuid,
		StatusUuid:   Uuid,
		CycleUuid:    Uuid,
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
		StatusUuid:   Uuid,
		CycleUuid:    Uuid,
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
		StatusUuid:   Uuid,
		CycleUuid:    Uuid,
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
		StatusUuid:   Uuid,
		CycleUuid:    Uuid,
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
	items       []entity.Issue
	statusItems []entity.IssueStatus
}

func (m mockRepository) GetStatus(ctx context.Context, id string) (entity.IssueStatus, error) {
	for _, item := range m.statusItems {
		if item.UUID == id {
			return item, nil
		}
	}
	return entity.IssueStatus{}, pg.ErrNoRows
}

func (m mockRepository) CountStatus(ctx context.Context) (int64, error) {
	return int64(len(m.statusItems)), nil
}

func (m mockRepository) QueryStatus(ctx context.Context, offset, limit int64) ([]entity.IssueStatus, int, error) {
	return m.statusItems, len(m.statusItems), nil
}

func (m mockRepository) CreateStatus(ctx context.Context, issueStatus entity.IssueStatus) error {
	if issueStatus.Title == "error" {
		return errCRUD
	}
	m.statusItems = append(m.statusItems, issueStatus)
	return nil
}

func (m mockRepository) UpdateStatus(ctx context.Context, issueStatus entity.IssueStatus) error {
	if issueStatus.Title == "error" {
		return errCRUD
	}
	for i, item := range m.statusItems {
		if item.UUID == issueStatus.UUID {
			m.statusItems[i] = issueStatus
			break
		}
	}
	return nil
}

func (m mockRepository) DeleteStatus(ctx context.Context, id string) error {
	for i, item := range m.statusItems {
		if item.UUID == id {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.Issue, error) {
	for _, item := range m.items {
		if item.UUID == id {
			return item, nil
		}
	}
	return entity.Issue{}, pg.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int64) ([]entity.Issue, int, error) {
	return m.items, len(m.items), nil
}

func (m *mockRepository) Create(ctx context.Context, issue entity.Issue) error {
	if issue.Title == "error" {
		return errCRUD
	}
	m.items = append(m.items, issue)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, issue entity.Issue) error {
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
