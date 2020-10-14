package cycles

import (
	"context"

	"github.com/mirzakhany/pm/internal/entity"

	"github.com/mirzakhany/pm/pkg/db"
)

// Repository encapsulates the logic to access cycles from the data source.
type Repository interface {
	// Get returns the cycle with the specified cycle UUID.
	Get(ctx context.Context, uuid string) (entity.Cycle, error)
	// Count returns the number of cycles.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of cycles with the given offset and limit.
	Query(ctx context.Context, offset, limit int64) ([]entity.Cycle, int, error)
	// Create saves a new cycle in the storage.
	Create(ctx context.Context, cycle entity.Cycle) error
	// Update updates the cycle with given UUID in the storage.
	Update(ctx context.Context, cycle entity.Cycle) error
	// Delete removes the cycle with given UUID from the storage.
	Delete(ctx context.Context, uuid string) error
}

// repository persists cycles in database
type repository struct {
	db *db.DB
}

// NewRepository creates a new cycle repository
func NewRepository(db *db.DB) Repository {
	return repository{db}
}

// Get reads the cycle with the specified ID from the database.
func (r repository) Get(ctx context.Context, uuid string) (entity.Cycle, error) {
	var cycle entity.Cycle
	err := r.db.With(ctx).Model(&cycle).Where("i.uuid = ?", uuid).First()
	return cycle, err
}

// Create saves a new cycle record in the database.
// It returns the ID of the newly inserted cycle record.
func (r repository) Create(ctx context.Context, cycle entity.Cycle) error {
	_, err := r.db.With(ctx).Model(cycle).Insert()
	return err
}

// Update saves the changes to an cycle in the database.
func (r repository) Update(ctx context.Context, cycle entity.Cycle) error {
	_, err := r.db.With(ctx).Model(&cycle).WherePK().Update()
	return err
}

// Delete deletes an cycle with the specified ID from the database.
func (r repository) Delete(ctx context.Context, uuid string) error {
	cycle, err := r.Get(ctx, uuid)
	if err != nil {
		return err
	}
	_, err = r.db.With(ctx).Model(&cycle).WherePK().Delete()
	return err
}

// Count returns the number of the cycle records in the database.
func (r repository) Count(ctx context.Context) (int64, error) {
	var count int
	count, err := r.db.With(ctx).Model((*entity.Cycle)(nil)).Count()
	return int64(count), err
}

// Query retrieves the cycle records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int64) ([]entity.Cycle, int, error) {
	var _cycles []entity.Cycle
	count, err := r.db.With(ctx).Model(&_cycles).
		Order("id ASC").
		Limit(int(limit)).
		Offset(int(offset)).
		SelectAndCount()
	return _cycles, count, err
}
