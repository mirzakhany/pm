package roles

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
	db.ResetTables(t, database, "roles")
	repo := NewRepository(database)

	ctx := context.Background()
	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	testUuid := uuid.New().String()
	// create
	now := time.Now()
	err = repo.Create(ctx, RoleModel{
		UUID:      testUuid,
		Title:     "admin",
		CreatedAt: now,
		UpdatedAt: now,
	})
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, int64(1), count2-count)

	// get
	role, err := repo.Get(ctx, testUuid)
	assert.Nil(t, err)
	assert.Equal(t, "admin", role.Title)
	_, err = repo.Get(ctx, "test0")
	assert.Equal(t, sql.ErrNoRows, err)

	// update
	err = repo.Update(ctx, RoleModel{
		ID: 	   role.ID,
		UUID:      testUuid,
		Title:     "manager",
		CreatedAt: now,
		UpdatedAt: now,
	})
	assert.Nil(t, err)
	role, _ = repo.Get(ctx, testUuid)
	assert.Equal(t, "manager", role.Title)

	// query
	roles, err := repo.Query(ctx, 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, int64(len(roles)))

	// delete
	err = repo.Delete(ctx, testUuid)
	assert.Nil(t, err)
	_, err = repo.Get(ctx, testUuid)
	assert.Equal(t, sql.ErrNoRows, err)
	err = repo.Delete(ctx, testUuid)
	assert.Equal(t, sql.ErrNoRows, err)
}
