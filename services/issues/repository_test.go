package issues

import (
	"context"
	"testing"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/mirzakhany/pm/pkg/db"
	issuesProto "github.com/mirzakhany/pm/services/issues/proto"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	database := db.NewForTest(t, []interface{}{(*IssueModel)(nil)})
	db.ResetTables(t, database, "issues")
	repo := NewRepository(database)

	ctx := context.Background()
	// initial count
	count, err := repo.Count(ctx)
	assert.Nil(t, err)

	testUuid := uuid.New().String()
	// create
	now := time.Now()
	err = repo.Create(ctx, IssueModel{
		UUID:        testUuid,
		Title:       "issue1",
		Description: "issue1",
		Status:      issuesProto.IssueStatus_IN_BACKLOG,
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
	issue, err := repo.Get(ctx, testUuid)
	assert.Nil(t, err)
	assert.Equal(t, "issue1", issue.Title)
	_, err = repo.Get(ctx, "test0")
	assert.EqualError(t, pg.ErrNoRows, err.Error())

	// update
	updatedIssue := IssueModel{
		ID:          issue.ID,
		UUID:        issue.UUID,
		Title:       "issue2",
		Description: issue.Description,
		Status:      issue.Status,
		SprintID:    issue.SprintID,
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
	assert.EqualError(t, pg.ErrNoRows, err.Error())
	err = repo.Delete(ctx, testUuid)
	assert.EqualError(t, pg.ErrNoRows, err.Error())
}
