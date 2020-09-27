package users

import (
	"context"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"proj/pkg/db"
	"time"
)

const TableName = "users"

type UserModel struct {
	ID        uint64
	UUID      string
	Username  string
	Password  string
	Email     string
	Enable    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r UserModel) TableName() string {
	return TableName
}

// Repository encapsulates the logic to access users from the data source.
type Repository interface {
	// Get returns the user with the specified user UUID.
	Get(ctx context.Context, uuid string) (UserModel, error)
	// Count returns the number of users.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of users with the given offset and limit.
	Query(ctx context.Context, offset, limit int64) ([]UserModel, error)
	// Create saves a new user in the storage.
	Create(ctx context.Context, user UserModel) error
	// Update updates the user with given UUID in the storage.
	Update(ctx context.Context, user UserModel) error
	// Delete removes the user with given UUID from the storage.
	Delete(ctx context.Context, uuid string) error
}

// repository persists users in database
type repository struct {
	db *db.DB
}

// NewRepository creates a new user repository
func NewRepository(db *db.DB) Repository {
	return repository{db}
}

// Get reads the user with the specified ID from the database.
func (r repository) Get(ctx context.Context, uuid string) (UserModel, error) {
	var user UserModel
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"uuid": uuid}).One(&user)
	return user, err
}

// Create saves a new user record in the database.
// It returns the ID of the newly inserted user record.
func (r repository) Create(ctx context.Context, user UserModel) error {
	return r.db.With(ctx).Model(&user).Insert()
}

// Update saves the changes to an user in the database.
func (r repository) Update(ctx context.Context, user UserModel) error {
	return r.db.With(ctx).Model(&user).Update()
}

// Delete deletes an user with the specified ID from the database.
func (r repository) Delete(ctx context.Context, uuid string) error {
	user, err := r.Get(ctx, uuid)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&user).Delete()
}

// Count returns the number of the user records in the database.
func (r repository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.With(ctx).Select("COUNT(*)").From(TableName).Row(&count)
	return count, err
}

// Query retrieves the user records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int64) ([]UserModel, error) {
	var _users []UserModel
	err := r.db.With(ctx).
		Select().
		OrderBy("id").
		Offset(offset).
		Limit(limit).
		All(&_users)
	return _users, err
}
