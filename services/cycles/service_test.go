package cycles

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mirzakhany/pm/pkg/auth"

	usersProto "github.com/mirzakhany/pm/services/users/proto"

	userSrv "github.com/mirzakhany/pm/services/users"

	"github.com/go-pg/pg/v10"
	cycles "github.com/mirzakhany/pm/services/cycles/proto"
	"github.com/stretchr/testify/assert"
)

var errCRUD = errors.New("error crud")

func TestCreateCycleRequest_Validate(t *testing.T) {
	now := timestamppb.Now()
	tests := []struct {
		name      string
		model     cycles.CreateCycleRequest
		wantError bool
	}{
		{"success", cycles.CreateCycleRequest{Title: "test", Description: "test", StartAt: now, EndAt: now, Active: true}, false},
		{"required", cycles.CreateCycleRequest{Title: "", Description: "test", StartAt: now, EndAt: now, Active: true}, true},
		{"too long", cycles.CreateCycleRequest{Description: "test", StartAt: now, EndAt: now, Active: true, Title: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdateCycleRequest_Validate(t *testing.T) {
	now := timestamppb.Now()
	tests := []struct {
		name      string
		model     cycles.UpdateCycleRequest
		wantError bool
	}{
		{"success", cycles.UpdateCycleRequest{Title: "test", Description: "test", StartAt: now, EndAt: now, Active: true}, false},
		{"required", cycles.UpdateCycleRequest{Title: "", Description: "test", StartAt: now, EndAt: now, Active: true}, true},
		{"too long", cycles.UpdateCycleRequest{Description: "test", StartAt: now, EndAt: now, Active: true, Title: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
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

	ctx = auth.ContextWithUser(ctx, user1)
	now := timestamppb.Now()
	// successful creation
	cycle, err := s.Create(ctx, &cycles.CreateCycleRequest{Title: "test", Description: "test", StartAt: now, EndAt: now, Active: true})
	assert.Nil(t, err)
	assert.NotEmpty(t, cycle.Uuid)
	id := cycle.Uuid
	assert.Equal(t, "test", cycle.Title)
	assert.NotEmpty(t, cycle.CreatedAt)
	assert.NotEmpty(t, cycle.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// validation error in creation
	_, err = s.Create(ctx, &cycles.CreateCycleRequest{Title: "", Description: "test", StartAt: now, EndAt: now, Active: true})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// unexpected error in creation
	_, err = s.Create(ctx, &cycles.CreateCycleRequest{Title: "error", Description: "test", StartAt: now, EndAt: now, Active: true})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	_, _ = s.Create(ctx, &cycles.CreateCycleRequest{Title: "test2", Description: "test", StartAt: now, EndAt: now, Active: true})

	// update
	cycle, err = s.Update(ctx, &cycles.UpdateCycleRequest{Title: "test updated", Description: "test", StartAt: now, EndAt: now, Active: true, Uuid: id})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", cycle.Title)
	_, err = s.Update(ctx, &cycles.UpdateCycleRequest{Title: "test updated", Description: "test", StartAt: now, EndAt: now, Active: true, Uuid: "none"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, &cycles.UpdateCycleRequest{Title: "", Description: "test", StartAt: now, EndAt: now, Active: true, Uuid: id})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// unexpected error in update
	_, err = s.Update(ctx, &cycles.UpdateCycleRequest{Title: "error", Description: "test", StartAt: now, EndAt: now, Active: true, Uuid: id})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	cycle, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", cycle.Title)
	assert.Equal(t, id, cycle.Uuid)

	// query
	_cycles, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, int(_cycles.TotalCount))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	cycle, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, cycle.Uuid)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)
}

type mockRepository struct {
	items []CycleModel
}

func (m mockRepository) Get(ctx context.Context, id string) (CycleModel, error) {
	for _, item := range m.items {
		if item.UUID == id {
			return item, nil
		}
	}
	return CycleModel{}, pg.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int64) ([]CycleModel, int, error) {
	return m.items, len(m.items), nil
}

func (m *mockRepository) Create(ctx context.Context, cycle CycleModel) error {
	if cycle.Title == "error" {
		return errCRUD
	}
	m.items = append(m.items, cycle)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, cycle CycleModel) error {
	if cycle.Title == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.UUID == cycle.UUID {
			m.items[i] = cycle
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
