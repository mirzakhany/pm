package roles

import (
	"context"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"proj/pkg/db"
	"time"
)

const TableName = "roles"

type RoleModel struct {
	ID        uint64
	UUID      string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r RoleModel) TableName() string {
	return TableName
}

// Repository encapsulates the logic to access roles from the data source.
type Repository interface {
	// Get returns the role with the specified role UUID.
	Get(ctx context.Context, uuid string) (RoleModel, error)
	// Count returns the number of roles.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of roles with the given offset and limit.
	Query(ctx context.Context, offset, limit int64) ([]RoleModel, error)
	// Create saves a new role in the storage.
	Create(ctx context.Context, role RoleModel) error
	// Update updates the role with given UUID in the storage.
	Update(ctx context.Context, role RoleModel) error
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
func (r repository) Get(ctx context.Context, uuid string) (RoleModel, error) {
	var role RoleModel
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"uuid": uuid}).One(&role)
	return role, err
}

// Create saves a new role record in the database.
// It returns the ID of the newly inserted role record.
func (r repository) Create(ctx context.Context, role RoleModel) error {
	return r.db.With(ctx).Model(&role).Insert()
}

// Update saves the changes to an role in the database.
func (r repository) Update(ctx context.Context, role RoleModel) error {
	return r.db.With(ctx).Model(&role).Update()
}

// Delete deletes an role with the specified ID from the database.
func (r repository) Delete(ctx context.Context, uuid string) error {
	role, err := r.Get(ctx, uuid)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&role).Delete()
}

// Count returns the number of the role records in the database.
func (r repository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.With(ctx).Select("COUNT(*)").From(TableName).Row(&count)
	return count, err
}

// Query retrieves the role records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int64) ([]RoleModel, error) {
	var _roles []RoleModel
	err := r.db.With(ctx).
		Select().
		OrderBy("id").
		Offset(offset).
		Limit(limit).
		All(&_roles)
	return _roles, err
}
