package users

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"proj/pkg/db"
	"testing"
	"time"
)

func TestRepository(t *testing.T) {
	database := db.NewForTest(t)
	db.ResetTables(t, database, "users")
	repo := NewRepository(database)

	ctx := context.Background()
	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	testUuid := uuid.New().String()
	// create
	now := time.Now()
	err = repo.Create(ctx, UserModel{
		UUID:      testUuid,
		Username:  "admin",
		Password:  "admin",
		Email:     "admin@admin.com",
		Enable:    true,
		CreatedAt: now,
		UpdatedAt: now,
	})
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, int64(1), count2-count)

	// get
	user, err := repo.Get(ctx, testUuid)
	assert.Nil(t, err)
	assert.Equal(t, "admin", user.Username)
	_, err = repo.Get(ctx, "test0")
	assert.Equal(t, sql.ErrNoRows, err)

	// update
	err = repo.Update(ctx, UserModel{
		ID:        user.ID,
		UUID:      testUuid,
		Username:  "owner",
		Email:     "admin@admin.com",
		Enable:    true,
		CreatedAt: now,
		UpdatedAt: now,
	})
	assert.Nil(t, err)
	user, _ = repo.Get(ctx, testUuid)
	assert.Equal(t, "owner", user.Username)

	// query
	users, err := repo.Query(ctx, 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, int64(len(users)))

	// delete
	err = repo.Delete(ctx, testUuid)
	assert.Nil(t, err)
	_, err = repo.Get(ctx, testUuid)
	assert.Equal(t, sql.ErrNoRows, err)
	err = repo.Delete(ctx, testUuid)
	assert.Equal(t, sql.ErrNoRows, err)
}
