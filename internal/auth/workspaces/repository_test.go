package workspaces

import (
	"context"
	"testing"
	"time"

	"github.com/go-pg/pg"
	"github.com/google/uuid"
	"github.com/mirzakhany/pm/internal/entity"
	"github.com/mirzakhany/pm/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	database := db.NewForTest(t, []interface{}{(*entity.Workspace)(nil)})
	db.ResetTables(t, database, "workspaces")
	repo := NewRepository(database)

	ctx := context.Background()
	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	testUuid := uuid.New().String()
	// create
	now := time.Now()
	err = repo.Create(ctx, entity.Workspace{
		UUID:      testUuid,
		Title:     "example",
		Domain:    "example.com",
		CreatedAt: now,
		UpdatedAt: now,
	})
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, int64(1), count2-count)

	// get
	workspace, err := repo.Get(ctx, testUuid)
	assert.Nil(t, err)
	assert.Equal(t, "example", workspace.Title)
	_, err = repo.Get(ctx, "test0")
	assert.NotNil(t, err)
	assert.EqualError(t, pg.ErrNoRows, err.Error())

	// update
	err = repo.Update(ctx, entity.Workspace{
		ID:        workspace.ID,
		UUID:      testUuid,
		Title:     "example1",
		Domain:    "example.com",
		CreatedAt: now,
		UpdatedAt: now,
	})
	assert.Nil(t, err)
	workspace, _ = repo.Get(ctx, testUuid)
	assert.Equal(t, "example1", workspace.Title)

	// query
	_, count3, err := repo.Query(ctx, 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, int64(count3))

	// delete
	err = repo.Delete(ctx, testUuid)
	assert.Nil(t, err)
	_, err = repo.Get(ctx, testUuid)
	assert.NotNil(t, err)
	assert.EqualError(t, pg.ErrNoRows, err.Error())
	err = repo.Delete(ctx, testUuid)
	assert.NotNil(t, err)
	assert.EqualError(t, pg.ErrNoRows, err.Error())
}
