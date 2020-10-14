package cycles

import (
	"context"
	"testing"
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/mirzakhany/pm/internal/entity"
	"github.com/mirzakhany/pm/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	database := db.NewForTest(t, []interface{}{&entity.User{}, &entity.Cycle{}})
	db.ResetTables(t, database, "users", "cycles")
	repo := NewRepository(database)

	ctx := context.Background()
	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	testUuid := uuid.New().String()
	// create
	now := time.Now()
	err = repo.Create(ctx, entity.Cycle{
		UUID:        testUuid,
		Title:       "test",
		Description: "test",
		StartAt:     now,
		EndAt:       now,
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, int64(1), count2-count)

	// get
	cycle, err := repo.Get(ctx, testUuid)
	assert.Nil(t, err)
	assert.Equal(t, "test", cycle.Title)
	_, err = repo.Get(ctx, "test0")
	assert.NotNil(t, err)
	assert.EqualError(t, gorm.ErrRecordNotFound, err.Error())

	// update
	err = repo.Update(ctx, entity.Cycle{
		ID:          cycle.ID,
		UUID:        cycle.UUID,
		Title:       "test2",
		Description: "test",
		StartAt:     cycle.StartAt,
		EndAt:       cycle.EndAt,
		Active:      cycle.Active,
		CreatedAt:   cycle.CreatedAt,
		UpdatedAt:   cycle.UpdatedAt,
	})
	assert.Nil(t, err)
	cycle, _ = repo.Get(ctx, cycle.UUID)
	assert.Equal(t, "test2", cycle.Title)

	// query
	_, count3, err := repo.Query(ctx, 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, count3)

	// delete
	err = repo.Delete(ctx, testUuid)
	assert.Nil(t, err)
	_, err = repo.Get(ctx, testUuid)
	assert.NotNil(t, err)
	assert.EqualError(t, gorm.ErrRecordNotFound, err.Error())
	err = repo.Delete(ctx, testUuid)
	assert.NotNil(t, err)
	assert.EqualError(t, gorm.ErrRecordNotFound, err.Error())
}
