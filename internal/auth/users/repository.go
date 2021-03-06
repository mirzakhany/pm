package users

import (
	"context"

	"github.com/mirzakhany/pm/internal/entity"
	"github.com/mirzakhany/pm/pkg/db"
)

// Repository encapsulates the logic to access users from the data source.
type Repository interface {
	// Get returns the user with the specified user UUID.
	Get(ctx context.Context, uuid string) (entity.User, error)
	// Count returns the number of users.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of users with the given offset and limit.
	Query(ctx context.Context, offset, limit int64) ([]entity.User, int, error)
	// Create saves a new user in the storage.
	Create(ctx context.Context, user entity.User) error
	// Update updates the user with given UUID in the storage.
	Update(ctx context.Context, user entity.User) error
	// Delete removes the user with given UUID from the storage.
	Delete(ctx context.Context, uuid string) error
	// Where returns the list of users with the given condition
	Where(ctx context.Context, condition string, params ...interface{}) ([]entity.User, int, error)
	// WhereOne returns the one of users with the given condition
	WhereOne(ctx context.Context, condition string, params ...interface{}) (entity.User, error)
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
func (r repository) Get(ctx context.Context, uuid string) (entity.User, error) {
	var user entity.User
	err := r.db.With(ctx).Model(&user).Where("uuid = ?", uuid).First()
	return user, err
}

// Create saves a new user record in the database.
// It returns the ID of the newly inserted user record.
func (r repository) Create(ctx context.Context, user entity.User) error {
	_, err := r.db.With(ctx).Model(&user).Insert()
	return err
}

// Update saves the changes to an user in the database.
func (r repository) Update(ctx context.Context, user entity.User) error {
	_, err := r.db.With(ctx).Model(&user).WherePK().Update()
	return err
}

// Delete deletes an user with the specified ID from the database.
func (r repository) Delete(ctx context.Context, uuid string) error {
	user, err := r.Get(ctx, uuid)
	if err != nil {
		return err
	}
	_, err = r.db.With(ctx).Model(&user).WherePK().Delete()
	return err
}

// Count returns the number of the user records in the database.
func (r repository) Count(ctx context.Context) (int64, error) {
	var count int
	count, err := r.db.With(ctx).Model((*entity.User)(nil)).Count()
	return int64(count), err
}

// Query retrieves the user records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int64) ([]entity.User, int, error) {
	var _users []entity.User
	count, err := r.db.With(ctx).
		Model(&_users).
		Order("id ASC").
		Limit(int(limit)).
		Offset(int(offset)).
		SelectAndCount()
	return _users, count, err
}

// Where returns the list of users with the given condition
func (r repository) Where(ctx context.Context, condition string, params ...interface{}) ([]entity.User, int, error) {
	var _users []entity.User
	count, err := r.db.With(ctx).Model(&_users).Where(condition, params).SelectAndCount()
	return _users, count, err
}

// WhereOne returns the one of users with the given condition
func (r repository) WhereOne(ctx context.Context, condition string, params ...interface{}) (entity.User, error) {
	var user entity.User
	err := r.db.With(ctx).Model(&user).Where(condition, params...).First()
	return user, err
}
