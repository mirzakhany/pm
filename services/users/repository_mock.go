package users

import (
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
)

var errCRUD = errors.New("error crud")

type mockRepository struct {
	items []UserModel
}

func (m mockRepository) Get(ctx context.Context, id string) (UserModel, error) {
	for _, item := range m.items {
		if item.UUID == id {
			return item, nil
		}
	}
	return UserModel{}, pg.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int64) ([]UserModel, int, error) {
	return m.items, len(m.items), nil
}

func (m *mockRepository) Create(ctx context.Context, user UserModel) error {
	if user.Username == "error" {
		return errCRUD
	}
	m.items = append(m.items, user)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, user UserModel) error {
	if user.Username == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.UUID == user.UUID {
			m.items[i] = user
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
