package cycles

import (
	"context"
	"errors"

	"github.com/mirzakhany/pm/internal/entity"

	"github.com/mirzakhany/pm/internal/auth/users"

	"github.com/go-pg/pg"
)

var errCRUD = errors.New("error crud")

// NewServiceForTest creates a new user service for test.
func NewServiceForTest(userSrv users.Service) Service {
	return NewService(&mockRepository{}, userSrv)
}

type mockRepository struct {
	items []entity.Cycle
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.Cycle, error) {
	for _, item := range m.items {
		if item.UUID == id {
			return item, nil
		}
	}
	return entity.Cycle{}, pg.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int64) ([]entity.Cycle, int, error) {
	return m.items, len(m.items), nil
}

func (m *mockRepository) Create(ctx context.Context, cycle entity.Cycle) error {
	if cycle.Title == "error" {
		return errCRUD
	}
	m.items = append(m.items, cycle)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, cycle entity.Cycle) error {
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
