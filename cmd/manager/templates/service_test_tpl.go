package templates

const ServiceTestTmpl = `
package {{ .Pkg.NamePlural | lower }}

import (
	"context"
	"errors"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/mirzakhany/pm/internal/entity"
	"github.com/mirzakhany/pm/protobuf/{{ .Pkg.NamePlural | lower }}"
	"github.com/stretchr/testify/assert"
)

var errCRUD = errors.New("error crud")

func TestCreate{{ .Pkg.Name }}Request_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     {{ .Pkg.NamePlural | lower }}.Create{{ .Pkg.Name }}Request
		wantError bool
	}{
		{"success", {{ .Pkg.NamePlural | lower }}.Create{{ .Pkg.Name }}Request{Title: "test"}, false},
		{"required", {{ .Pkg.NamePlural | lower }}.Create{{ .Pkg.Name }}Request{Title: ""}, true},
		{"too long", {{ .Pkg.NamePlural | lower }}.Create{{ .Pkg.Name }}Request{Title: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdate{{ .Pkg.Name }}Request_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     {{ .Pkg.NamePlural | lower }}.Update{{ .Pkg.Name }}Request
		wantError bool
	}{
		{"success", {{ .Pkg.NamePlural | lower }}.Update{{ .Pkg.Name }}Request{Title: "test"}, false},
		{"required", {{ .Pkg.NamePlural | lower }}.Update{{ .Pkg.Name }}Request{Title: ""}, true},
		{"too long", {{ .Pkg.NamePlural | lower }}.Update{{ .Pkg.Name }}Request{Title: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUpdateRequest(&tt.model)
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func Test_service_CRUD(t *testing.T) {
	s := NewService(&mockRepository{})
	ctx := context.Background()

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, int64(0), count)

	// successful creation
	{{ .Pkg.Name | lower }}, err := s.Create(ctx, &{{ .Pkg.NamePlural | lower }}.Create{{ .Pkg.Name }}Request{Title: "test"})
	assert.Nil(t, err)
	assert.NotEmpty(t, {{ .Pkg.Name | lower }}.Uuid)
	id := {{ .Pkg.Name | lower }}.Uuid
	assert.Equal(t, "test", {{ .Pkg.Name | lower }}.Title)
	assert.NotEmpty(t, {{ .Pkg.Name | lower }}.CreatedAt)
	assert.NotEmpty(t, {{ .Pkg.Name | lower }}.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// validation error in creation
	_, err = s.Create(ctx, &{{ .Pkg.NamePlural | lower }}.Create{{ .Pkg.Name }}Request{Title: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	// unexpected error in creation
	_, err = s.Create(ctx, &{{ .Pkg.NamePlural | lower }}.Create{{ .Pkg.Name }}Request{Title: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)

	_, _ = s.Create(ctx, &{{ .Pkg.NamePlural | lower }}.Create{{ .Pkg.Name }}Request{Title: "test2"})

	// update
	{{ .Pkg.Name | lower }}, err = s.Update(ctx, &{{ .Pkg.NamePlural | lower }}.Update{{ .Pkg.Name }}Request{Title: "test updated", Uuid: id})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", {{ .Pkg.Name | lower }}.Title)
	_, err = s.Update(ctx, &{{ .Pkg.NamePlural | lower }}.Update{{ .Pkg.Name }}Request{Title: "test updated", Uuid: "none"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, &{{ .Pkg.NamePlural | lower }}.Update{{ .Pkg.Name }}Request{Title: "", Uuid: id})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// unexpected error in update
	_, err = s.Update(ctx, &{{ .Pkg.NamePlural | lower }}.Update{{ .Pkg.Name }}Request{Title: "error", Uuid: id})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(2), count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	{{ .Pkg.Name | lower }}, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", {{ .Pkg.Name | lower }}.Title)
	assert.Equal(t, id, {{ .Pkg.Name | lower }}.Uuid)

	// query
	_{{ .Pkg.NamePlural | lower }}, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, int(_{{ .Pkg.NamePlural | lower }}.TotalCount))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	{{ .Pkg.Name | lower }}, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, {{ .Pkg.Name | lower }}.Uuid)
	count, _ = s.Count(ctx)
	assert.Equal(t, int64(1), count)
}

type mockRepository struct {
	items []entity.{{ .Pkg.Name }}
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.{{ .Pkg.Name }}, error) {
	for _, item := range m.items {
		if item.UUID == id {
			return item, nil
		}
	}
	return entity.{{ .Pkg.Name }}{}, pg.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int64, error) {
	return int64(len(m.items)), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int64) ([]entity.{{ .Pkg.Name }}, int, error) {
	return m.items, len(m.items), nil
}

func (m *mockRepository) Create(ctx context.Context, {{ .Pkg.Name | lower }} entity.{{ .Pkg.Name }}) error {
	if {{ .Pkg.Name | lower }}.Title == "error" {
		return errCRUD
	}
	m.items = append(m.items, {{ .Pkg.Name | lower }})
	return nil
}

func (m *mockRepository) Update(ctx context.Context, {{ .Pkg.Name | lower }} entity.{{ .Pkg.Name }}) error {
	if {{ .Pkg.Name | lower }}.Title == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.UUID == {{ .Pkg.Name | lower }}.UUID {
			m.items[i] = {{ .Pkg.Name | lower }}
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

`
