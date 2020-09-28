package tasks

import (
	"context"
	"testing"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/mirzakhany/pm/pkg/db"
	tasksProto "github.com/mirzakhany/pm/services/tasks/proto"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	database := db.NewForTest(t, []interface{}{(*TaskModel)(nil)})
	db.ResetTables(t, database, "tasks")
	repo := NewRepository(database)

	ctx := context.Background()
	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	testUuid := uuid.New().String()
	// create
	now := time.Now()
	err = repo.Create(ctx, TaskModel{
		UUID:        testUuid,
		Title:       "task1",
		Description: "task1",
		Status:      tasksProto.TaskStatus_IN_BACKLOG,
		SprintID:    0,
		Estimate:    4,
		AssigneeID:  0,
		CreatorID:   0,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	assert.Nil(t, err)
	count2, _ := repo.Count(ctx)
	assert.Equal(t, int64(1), count2-count)

	// get
	task, err := repo.Get(ctx, testUuid)
	assert.Nil(t, err)
	assert.Equal(t, "task1", task.Title)
	_, err = repo.Get(ctx, "test0")
	assert.EqualError(t, pg.ErrNoRows, err.Error())

	// update
	updatedTask := TaskModel{
		ID:          task.ID,
		UUID:        task.UUID,
		Title:       "task2",
		Description: task.Description,
		Status:      task.Status,
		SprintID:    task.SprintID,
		Estimate:    task.Estimate,
		AssigneeID:  task.AssigneeID,
		CreatorID:   task.CreatorID,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   now,
	}
	err = repo.Update(ctx, updatedTask)
	assert.Nil(t, err)
	task2, _ := repo.Get(ctx, task.UUID)
	assert.Equal(t, "task2", task2.Title)

	// query
	_, count3, err := repo.Query(ctx, 0, count2)
	assert.Nil(t, err)
	assert.Equal(t, count2, int64(count3))

	// delete
	err = repo.Delete(ctx, testUuid)
	assert.Nil(t, err)
	_, err = repo.Get(ctx, testUuid)
	assert.EqualError(t, pg.ErrNoRows, err.Error())
	err = repo.Delete(ctx, testUuid)
	assert.EqualError(t, pg.ErrNoRows, err.Error())
}
