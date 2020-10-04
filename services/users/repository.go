package users

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes"

	users "github.com/mirzakhany/pm/services/users/proto"

	"github.com/mirzakhany/pm/pkg/db"
)

type UserModel struct {
	tableName struct{} `pg:"users,alias:u"` //nolint
	ID        uint64   `pg:",pk"`
	UUID      string
	Username  string `pg:",unique"`
	Password  string
	Email     string `pg:",unique"`
	Enable    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Repository encapsulates the logic to access users from the data source.
type Repository interface {
	// Get returns the user with the specified user UUID.
	Get(ctx context.Context, uuid string) (UserModel, error)
	// Count returns the number of users.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of users with the given offset and limit.
	Query(ctx context.Context, offset, limit int64) ([]UserModel, int, error)
	// Create saves a new user in the storage.
	Create(ctx context.Context, user UserModel) error
	// Update updates the user with given UUID in the storage.
	Update(ctx context.Context, user UserModel) error
	// Delete removes the user with given UUID from the storage.
	Delete(ctx context.Context, uuid string) error
	// Where returns the list of users with the given condition
	Where(ctx context.Context, condition string, params ...interface{}) ([]UserModel, int, error)
	// WhereOne returns the one of users with the given condition
	WhereOne(ctx context.Context, condition string, params ...interface{}) (UserModel, error)
}

func (um UserModel) ToProto(secure bool) *users.User {
	c, _ := ptypes.TimestampProto(um.CreatedAt)
	u, _ := ptypes.TimestampProto(um.UpdatedAt)

	user := &users.User{
		Id:        um.ID,
		Uuid:      um.UUID,
		Username:  um.Username,
		Password:  um.Password,
		Email:     um.Email,
		Enable:    um.Enable,
		CreatedAt: c,
		UpdatedAt: u,
	}
	if secure {
		user.Id = 0
		user.Password = ""
	}
	return user
}

func ToProtoList(uml []UserModel, secure bool) []*users.User {
	var u []*users.User
	for _, i := range uml {
		u = append(u, i.ToProto(secure))
	}
	return u
}

func FromProto(user *users.User) UserModel {
	c, _ := ptypes.Timestamp(user.CreatedAt)
	u, _ := ptypes.Timestamp(user.UpdatedAt)
	return UserModel{
		ID:        user.Id,
		UUID:      user.Uuid,
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		Enable:    user.Enable,
		CreatedAt: c,
		UpdatedAt: u,
	}
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
	err := r.db.With(ctx).Model(&user).Where("uuid = ?", uuid).First()
	return user, err
}

// Create saves a new user record in the database.
// It returns the ID of the newly inserted user record.
func (r repository) Create(ctx context.Context, user UserModel) error {
	_, err := r.db.With(ctx).Model(&user).Insert()
	return err
}

// Update saves the changes to an user in the database.
func (r repository) Update(ctx context.Context, user UserModel) error {
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
	count, err := r.db.With(ctx).Model((*UserModel)(nil)).Count()
	return int64(count), err
}

// Query retrieves the user records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int64) ([]UserModel, int, error) {
	var _users []UserModel
	count, err := r.db.With(ctx).
		Model(&_users).
		Order("id ASC").
		Limit(int(limit)).
		Offset(int(offset)).
		SelectAndCount()
	return _users, count, err
}

// Where returns the list of users with the given condition
func (r repository) Where(ctx context.Context, condition string, params ...interface{}) ([]UserModel, int, error) {
	var _users []UserModel
	count, err := r.db.With(ctx).Model(&_users).Where(condition, params).SelectAndCount()
	return _users, count, err
}

// WhereOne returns the one of users with the given condition
func (r repository) WhereOne(ctx context.Context, condition string, params ...interface{}) (UserModel, error) {
	var user UserModel
	err := r.db.With(ctx).Model(&user).Where(condition, params...).First()
	return user, err
}
