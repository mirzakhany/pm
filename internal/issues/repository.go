package issues

import (
	"context"

	"github.com/mirzakhany/pm/internal/entity"

	"github.com/mirzakhany/pm/pkg/db"
)

// Repository encapsulates the logic to access issues from the data source.
type Repository interface {
	// Get returns the issue with the specified issue UUID.
	Get(ctx context.Context, uuid string) (entity.Issue, error)
	// Count returns the number of issues.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of issues with the given offset and limit.
	Query(ctx context.Context, offset, limit int64) ([]entity.Issue, int, error)
	// Create saves a new issue in the storage.
	Create(ctx context.Context, issue entity.Issue) error
	// Update updates the issue with given UUID in the storage.
	Update(ctx context.Context, issue entity.Issue) error
	// Delete removes the issue with given UUID from the storage.
	Delete(ctx context.Context, uuid string) error

	//IssueStatus

	// GetStatus returns the status with the specified issue UUID.
	GetStatus(ctx context.Context, uuid string) (entity.IssueStatus, error)
	// CountStatus returns the number of status.
	CountStatus(ctx context.Context) (int64, error)
	// QueryStatus returns the list of status with the given offset and limit.
	QueryStatus(ctx context.Context, offset, limit int64) ([]entity.IssueStatus, int, error)
	// CreateStatus saves a new status in the storage.
	CreateStatus(ctx context.Context, issueStatus entity.IssueStatus) error
	// UpdateStatus updates the status with given UUID in the storage.
	UpdateStatus(ctx context.Context, issueStatus entity.IssueStatus) error
	// DeleteStatus removes the status with given UUID from the storage.
	DeleteStatus(ctx context.Context, uuid string) error
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
func (r repository) Get(ctx context.Context, uuid string) (entity.Issue, error) {
	var issue entity.Issue
	err := r.db.With(ctx).Model(&issue).
		Relation("Assignee").
		Relation("Creator").
		Relation("Cycle").
		Where("i.uuid = ?", uuid).First()

	return issue, err
}

// Create saves a new issue record in the database.
func (r repository) Create(ctx context.Context, issue entity.Issue) error {
	_, err := r.db.With(ctx).Model(&issue).Insert()
	return err
}

// Update saves the changes to an issue in the database.
func (r repository) Update(ctx context.Context, issue entity.Issue) error {
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
	count, err := r.db.With(ctx).Model((*entity.Issue)(nil)).Count()
	return int64(count), err
}

// Query retrieves the issue records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int64) ([]entity.Issue, int, error) {
	var _issues []entity.Issue
	count, err := r.db.With(ctx).Model(&_issues).
		Relation("Assignee").
		Relation("Creator").
		Relation("Cycle").
		Order("id ASC").
		Limit(int(limit)).
		Offset(int(offset)).
		SelectAndCount()
	return _issues, count, err
}

func (r repository) GetStatus(ctx context.Context, uuid string) (entity.IssueStatus, error) {
	var issueStatus entity.IssueStatus
	err := r.db.With(ctx).Model(&issueStatus).Where("i.uuid = ?", uuid).First()
	return issueStatus, err
}

func (r repository) CountStatus(ctx context.Context) (int64, error) {
	var count int
	count, err := r.db.With(ctx).Model((*entity.IssueStatus)(nil)).Count()
	return int64(count), err
}

func (r repository) QueryStatus(ctx context.Context, offset, limit int64) ([]entity.IssueStatus, int, error) {
	var _issueStatus []entity.IssueStatus
	count, err := r.db.With(ctx).Model(&_issueStatus).
		Limit(int(limit)).
		Offset(int(offset)).
		SelectAndCount()
	return _issueStatus, count, err
}

// CreateStatus saves a new status record in the database.
func (r repository) CreateStatus(ctx context.Context, issueStatus entity.IssueStatus) error {
	_, err := r.db.With(ctx).Model(&issueStatus).Insert()
	return err
}

func (r repository) UpdateStatus(ctx context.Context, issueStatus entity.IssueStatus) error {
	_, err := r.db.With(ctx).Model(&issueStatus).WherePK().Update()
	return err
}

func (r repository) DeleteStatus(ctx context.Context, uuid string) error {
	issueStatus, err := r.Get(ctx, uuid)
	if err != nil {
		return err
	}
	_, err = r.db.With(ctx).Model(&issueStatus).WherePK().Delete()
	return err
}
