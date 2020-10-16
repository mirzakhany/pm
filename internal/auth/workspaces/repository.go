package workspaces

import (
	"context"

	"github.com/mirzakhany/pm/internal/entity"

	"github.com/mirzakhany/pm/pkg/db"
)

// Repository encapsulates the logic to access workspaces from the data source.
type Repository interface {
	// Get returns the workspace with the specified workspace UUID.
	Get(ctx context.Context, uuid string) (entity.Workspace, error)
	// Count returns the number of workspaces.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of workspaces with the given offset and limit.
	Query(ctx context.Context, offset, limit int64) ([]entity.Workspace, int, error)
	// Create saves a new workspace in the storage.
	Create(ctx context.Context, workspace entity.Workspace) error
	// Update updates the workspace with given UUID in the storage.
	Update(ctx context.Context, workspace entity.Workspace) error
	// Delete removes the workspace with given UUID from the storage.
	Delete(ctx context.Context, uuid string) error
}

// repository persists workspaces in database
type repository struct {
	db *db.DB
}

// NewRepository creates a new workspace repository
func NewRepository(db *db.DB) Repository {
	return repository{db}
}

// Get reads the workspace with the specified ID from the database.
func (r repository) Get(ctx context.Context, uuid string) (entity.Workspace, error) {
	var workspace entity.Workspace
	err := r.db.With(ctx).Model(&workspace).Where("uuid = ?", uuid).First()
	return workspace, err
}

// Create saves a new workspace record in the database.
// It returns the ID of the newly inserted workspace record.
func (r repository) Create(ctx context.Context, workspace entity.Workspace) error {
	_, err := r.db.With(ctx).Model(&workspace).Insert()
	return err
}

// Update saves the changes to an workspace in the database.
func (r repository) Update(ctx context.Context, workspace entity.Workspace) error {
	_, err := r.db.With(ctx).Model(&workspace).WherePK().Update()
	return err
}

// Delete deletes an workspace with the specified ID from the database.
func (r repository) Delete(ctx context.Context, uuid string) error {
	workspace, err := r.Get(ctx, uuid)
	if err != nil {
		return err
	}
	_, err = r.db.With(ctx).Model(&workspace).WherePK().Delete()
	return err
}

// Count returns the number of the workspace records in the database.
func (r repository) Count(ctx context.Context) (int64, error) {
	var count int
	count, err := r.db.With(ctx).Model((*entity.Workspace)(nil)).Count()
	return int64(count), err
}

// Query retrieves the workspace records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int64) ([]entity.Workspace, int, error) {
	var _workspaces []entity.Workspace
	count, err := r.db.With(ctx).Model(&_workspaces).
		Order("id ASC").Limit(int(limit)).
		Offset(int(offset)).SelectAndCount()
	return _workspaces, count, err
}
