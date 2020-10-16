package templates

const RepositoryTmpl = `
package {{ .Pkg.NamePlural | lower }}

import (
	"context"

	"github.com/mirzakhany/pm/internal/entity"

	"github.com/mirzakhany/pm/pkg/db"
)

// Repository encapsulates the logic to access {{ .Pkg.NamePlural | lower }} from the data source.
type Repository interface {
	// Get returns the {{ .Pkg.Name | lower }} with the specified {{ .Pkg.Name | lower }} UUID.
	Get(ctx context.Context, uuid string) (entity.{{ .Pkg.Name }}, error)
	// Count returns the number of {{ .Pkg.NamePlural | lower }}.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of {{ .Pkg.NamePlural | lower }} with the given offset and limit.
	Query(ctx context.Context, offset, limit int64) ([]entity.{{ .Pkg.Name }}, int, error)
	// Create saves a new {{ .Pkg.Name | lower }} in the storage.
	Create(ctx context.Context, {{ .Pkg.Name | lower }} entity.{{ .Pkg.Name }}) error
	// Update updates the {{ .Pkg.Name | lower }} with given UUID in the storage.
	Update(ctx context.Context, {{ .Pkg.Name | lower }} entity.{{ .Pkg.Name }}) error
	// Delete removes the {{ .Pkg.Name | lower }} with given UUID from the storage.
	Delete(ctx context.Context, uuid string) error
}

// repository persists {{ .Pkg.NamePlural | lower }} in database
type repository struct {
	db *db.DB
}

// NewRepository creates a new {{ .Pkg.Name | lower }} repository
func NewRepository(db *db.DB) Repository {
	return repository{db}
}

// Get reads the {{ .Pkg.Name | lower }} with the specified ID from the database.
func (r repository) Get(ctx context.Context, uuid string) (entity.{{ .Pkg.Name }}, error) {
	var {{ .Pkg.Name | lower }} entity.{{ .Pkg.Name }}
	err := r.db.With(ctx).Model(&{{ .Pkg.Name | lower }}).Where("uuid = ?", uuid).First()
	return {{ .Pkg.Name | lower }}, err
}

// Create saves a new {{ .Pkg.Name | lower }} record in the database.
// It returns the ID of the newly inserted {{ .Pkg.Name | lower }} record.
func (r repository) Create(ctx context.Context, {{ .Pkg.Name | lower }} entity.{{ .Pkg.Name }}) error {
	_, err := r.db.With(ctx).Model(&{{ .Pkg.Name | lower }}).Insert()
	return err
}

// Update saves the changes to an {{ .Pkg.Name | lower }} in the database.
func (r repository) Update(ctx context.Context, {{ .Pkg.Name | lower }} entity.{{ .Pkg.Name }}) error {
	_, err := r.db.With(ctx).Model(&{{ .Pkg.Name | lower }}).WherePK().Update()
	return err
}

// Delete deletes an {{ .Pkg.Name | lower }} with the specified ID from the database.
func (r repository) Delete(ctx context.Context, uuid string) error {
	{{ .Pkg.Name | lower }}, err := r.Get(ctx, uuid)
	if err != nil {
		return err
	}
	_, err = r.db.With(ctx).Model(&{{ .Pkg.Name | lower }}).WherePK().Delete()
	return err
}

// Count returns the number of the {{ .Pkg.Name | lower }} records in the database.
func (r repository) Count(ctx context.Context) (int64, error) {
	var count int
	count, err := r.db.With(ctx).Model((*entity.{{ .Pkg.Name }})(nil)).Count()
	return int64(count), err
}

// Query retrieves the {{ .Pkg.Name | lower }} records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int64) ([]entity.{{ .Pkg.Name }}, int, error) {
	var _{{ .Pkg.NamePlural | lower }} []entity.{{ .Pkg.Name }}
	count, err := r.db.With(ctx).Model(&_{{ .Pkg.NamePlural | lower }}).
		Order("id ASC").Limit(int(limit)).
		Offset(int(offset)).SelectAndCount()
	return _{{ .Pkg.NamePlural | lower }}, count, err
}

`
