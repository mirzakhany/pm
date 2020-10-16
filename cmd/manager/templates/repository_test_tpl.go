package templates

const RepositoryTestTmpl = `
package {{ .Pkg.NamePlural | lower }}

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
	database := db.NewForTest(t, []interface{}{(*entity.{{ .Pkg.Name }})(nil)})
	db.ResetTables(t, database, "{{ .Pkg.NamePlural | lower }}")
	repo := NewRepository(database)

	ctx := context.Background()
	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	testUuid := uuid.New().String()
	// create
	now := time.Now()
	err = repo.Create(ctx, entity.{{ .Pkg.Name }}{
		UUID:      testUuid,
		Title:     "admin",
		CreatedAt: now,
		UpdatedAt: now,
	})
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, int64(1), count2-count)

	// get
	{{ .Pkg.Name | lower }}, err := repo.Get(ctx, testUuid)
	assert.Nil(t, err)
	assert.Equal(t, "admin", {{ .Pkg.Name | lower }}.Title)
	_, err = repo.Get(ctx, "test0")
	assert.NotNil(t, err)
	assert.EqualError(t, pg.ErrNoRows, err.Error())

	// update
	err = repo.Update(ctx, entity.{{ .Pkg.Name }}{
		ID:        {{ .Pkg.Name | lower }}.ID,
		UUID:      testUuid,
		Title:     "manager",
		CreatedAt: now,
		UpdatedAt: now,
	})
	assert.Nil(t, err)
	{{ .Pkg.Name | lower }}, _ = repo.Get(ctx, testUuid)
	assert.Equal(t, "manager", {{ .Pkg.Name | lower }}.Title)

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

`
