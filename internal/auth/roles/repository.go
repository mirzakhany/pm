package roles

import (
	"context"

	"github.com/mirzakhany/pm/internal/entity"

	"github.com/mirzakhany/pm/pkg/db"
)

// Repository encapsulates the logic to access roles from the data source.
type Repository interface {
	// Get returns the role with the specified role UUID.
	Get(ctx context.Context, uuid string) (entity.Role, error)
	// Count returns the number of roles.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of roles with the given offset and limit.
	Query(ctx context.Context, offset, limit int64) ([]entity.Role, int, error)
	// Create saves a new role in the storage.
	Create(ctx context.Context, role entity.Role) error
	// Update updates the role with given UUID in the storage.
	Update(ctx context.Context, role entity.Role) error
	// Delete removes the role with given UUID from the storage.
	Delete(ctx context.Context, uuid string) error
}

// repository persists roles in database
type repository struct {
	db *db.DB
}

// NewRepository creates a new role repository
func NewRepository(db *db.DB) Repository {
	return repository{db}
}

// Get reads the role with the specified ID from the database.
func (r repository) Get(ctx context.Context, uuid string) (entity.Role, error) {
	var role entity.Role
	err := r.db.With(ctx).Model(&role).Where("uuid = ?", uuid).First()
	return role, err
}

// Create saves a new role record in the database.
// It returns the ID of the newly inserted role record.
func (r repository) Create(ctx context.Context, role entity.Role) error {
	_, err := r.db.With(ctx).Model(&role).Insert()
	return err
}

// Update saves the changes to an role in the database.
func (r repository) Update(ctx context.Context, role entity.Role) error {
	_, err := r.db.With(ctx).Model(&role).WherePK().Update()
	return err
}

// Delete deletes an role with the specified ID from the database.
func (r repository) Delete(ctx context.Context, uuid string) error {
	role, err := r.Get(ctx, uuid)
	if err != nil {
		return err
	}
	_, err = r.db.With(ctx).Model(&role).WherePK().Delete()
	return err
}

// Count returns the number of the role records in the database.
func (r repository) Count(ctx context.Context) (int64, error) {
	var count int
	count, err := r.db.With(ctx).Model((*entity.Role)(nil)).Count()
	return int64(count), err
}

// Query retrieves the role records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int64) ([]entity.Role, int, error) {
	var _roles []entity.Role
	count, err := r.db.With(ctx).Model(&_roles).
		Order("id ASC").Limit(int(limit)).
		Offset(int(offset)).SelectAndCount()
	return _roles, count, err
}
