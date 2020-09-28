package issues

import (
	"context"
	"time"

	"github.com/mirzakhany/pm/pkg/db"
	issuesProto "github.com/mirzakhany/pm/services/issues/proto"
	"github.com/mirzakhany/pm/services/users"
)

type IssueModel struct {
	tableName   struct{} `pg:"issues,alias:i"` //nolint
	ID          uint64   `pg:",pk"`
	UUID        string
	Title       string
	Description string
	Status      issuesProto.IssueStatus
	SprintID    uint64
	Estimate    uint64
	AssigneeID  uint64           `pg:",pk"`
	Assignee    *users.UserModel `pg:"rel:has-one"`
	CreatorID   uint64           `pg:",pk"`
	Creator     *users.UserModel `pg:"rel:has-one"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Repository encapsulates the logic to access issues from the data source.
type Repository interface {
	// Get returns the issue with the specified issue UUID.
	Get(ctx context.Context, uuid string) (IssueModel, error)
	// Count returns the number of issues.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of issues with the given offset and limit.
	Query(ctx context.Context, offset, limit int64) ([]IssueModel, int, error)
	// Create saves a new issue in the storage.
	Create(ctx context.Context, issue IssueModel) error
	// Update updates the issue with given UUID in the storage.
	Update(ctx context.Context, issue IssueModel) error
	// Delete removes the issue with given UUID from the storage.
	Delete(ctx context.Context, uuid string) error
}

// repository persists issues in database
type repository struct {
	db *db.DB
}

// NewRepository creates a new issue repository
func NewRepository(db *db.DB) Repository {
	return repository{db}
}

// Get reads the issue with the specified ID from the database.
func (r repository) Get(ctx context.Context, uuid string) (IssueModel, error) {
	var issue IssueModel
	err := r.db.With(ctx).Model(&issue).
		Relation("Assignee").Relation("Creator").
		Where("i.uuid = ?", uuid).First()

	return issue, err
}

// Create saves a new issue record in the database.
// It returns the ID of the newly inserted issue record.
func (r repository) Create(ctx context.Context, issue IssueModel) error {
	_, err := r.db.With(ctx).Model(&issue).Insert()
	return err
}

// Update saves the changes to an issue in the database.
func (r repository) Update(ctx context.Context, issue IssueModel) error {
	_, err := r.db.With(ctx).Model(&issue).WherePK().Update()
	return err
}

// Delete deletes an issue with the specified ID from the database.
func (r repository) Delete(ctx context.Context, uuid string) error {
	issue, err := r.Get(ctx, uuid)
	if err != nil {
		return err
	}
	_, err = r.db.With(ctx).Model(&issue).WherePK().Delete()
	return err
}

// Count returns the number of the issue records in the database.
func (r repository) Count(ctx context.Context) (int64, error) {
	var count int
	count, err := r.db.With(ctx).Model((*IssueModel)(nil)).Count()
	return int64(count), err
}

// Query retrieves the issue records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int64) ([]IssueModel, int, error) {
	var _issues []IssueModel
	count, err := r.db.With(ctx).Model(&_issues).
		Relation("Assignee").Relation("Creator").
		Order("id ASC").Limit(int(limit)).Offset(int(offset)).SelectAndCount()
	return _issues, count, err
}
