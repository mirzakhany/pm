package issues

import (
	"context"
	"testing"
	"time"

	"github.com/mirzakhany/pm/internal/entity"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/mirzakhany/pm/pkg/db"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	database := db.NewForTest(t, []interface{}{(*entity.User)(nil), (*entity.Cycle)(nil), (*entity.Issue)(nil), (*entity.IssueStatus)(nil)})
	db.ResetTables(t, database, "issues", "users", "cycles")
	repo := NewRepository(database)

	ctx := context.Background()
	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	testUuid := uuid.New().String()
	// create
	now := time.Now()
	err = repo.Create(ctx, entity.Issue{
		UUID:        testUuid,
		Title:       "issue1",
		Description: "issue1",
		StatusID:    0,
		CycleID:     0,
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
	issue, err := repo.Get(ctx, testUuid)
	assert.Nil(t, err)
	assert.Equal(t, "issue1", issue.Title)
	_, err = repo.Get(ctx, "test0")
	assert.NotNil(t, err)
	assert.EqualError(t, pg.ErrNoRows, err.Error())

	// update
	updatedIssue := entity.Issue{
		ID:          issue.ID,
		UUID:        issue.UUID,
		Title:       "issue2",
		Description: issue.Description,
		Status:      issue.Status,
		CycleID:     issue.CycleID,
		Estimate:    issue.Estimate,
		AssigneeID:  issue.AssigneeID,
		CreatorID:   issue.CreatorID,
		CreatedAt:   issue.CreatedAt,
		UpdatedAt:   now,
	}
	err = repo.Update(ctx, updatedIssue)
	assert.Nil(t, err)
	issue2, _ := repo.Get(ctx, issue.UUID)
	assert.Equal(t, "issue2", issue2.Title)

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
