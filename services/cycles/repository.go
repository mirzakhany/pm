package cycles

import (
	"context"
	"time"

	"github.com/mirzakhany/pm/services/users"

	"github.com/mirzakhany/pm/pkg/db"
)

type CycleModel struct {
	tableName   struct{} `pg:"cycles,alias:c"` //nolint
	ID          uint64   `pg:",pk"`
	UUID        string
	Title       string
	Description string
	Active      bool
	CreatorID   uint64           `pg:",pk"`
	Creator     *users.UserModel `pg:"rel:has-one"`
	StartAt     time.Time
	EndAt       time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Repository encapsulates the logic to access cycles from the data source.
type Repository interface {
	// Get returns the cycle with the specified cycle UUID.
	Get(ctx context.Context, uuid string) (CycleModel, error)
	// Count returns the number of cycles.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of cycles with the given offset and limit.
	Query(ctx context.Context, offset, limit int64) ([]CycleModel, int, error)
	// Create saves a new cycle in the storage.
	Create(ctx context.Context, cycle CycleModel) error
	// Update updates the cycle with given UUID in the storage.
	Update(ctx context.Context, cycle CycleModel) error
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
func (r repository) Get(ctx context.Context, uuid string) (CycleModel, error) {
	var cycle CycleModel
	err := r.db.With(ctx).Model(&cycle).Where("uuid = ?", uuid).First()
	return cycle, err
}

// Create saves a new cycle record in the database.
// It returns the ID of the newly inserted cycle record.
func (r repository) Create(ctx context.Context, cycle CycleModel) error {
	_, err := r.db.With(ctx).Model(&cycle).Insert()
	return err
}

// Update saves the changes to an cycle in the database.
func (r repository) Update(ctx context.Context, cycle CycleModel) error {
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
	count, err := r.db.With(ctx).Model((*CycleModel)(nil)).Count()
	return int64(count), err
}

// Query retrieves the cycle records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int64) ([]CycleModel, int, error) {
	var _cycles []CycleModel
	count, err := r.db.With(ctx).Model(&_cycles).Relation("Creator").
		Order("id ASC").Limit(int(limit)).
		Offset(int(offset)).SelectAndCount()
	return _cycles, count, err
}
