package tasks

import (
	"context"
	"time"

	"github.com/mirzakhany/pm/pkg/db"
	tasksProto "github.com/mirzakhany/pm/services/tasks/proto"
	"github.com/mirzakhany/pm/services/users"
)

type TaskModel struct {
	tableName   struct{} `pg:"tasks,alias:t"` //nolint
	ID          uint64   `pg:",pk"`
	UUID        string
	Title       string
	Description string
	Status      tasksProto.TaskStatus
	SprintID    uint64
	Estimate    uint64
	AssigneeID  uint64           `pg:",pk"`
	Assignee    *users.UserModel `pg:"rel:has-one"`
	CreatorID   uint64           `pg:",pk"`
	Creator     *users.UserModel `pg:"rel:has-one"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Repository encapsulates the logic to access tasks from the data source.
type Repository interface {
	// Get returns the task with the specified task UUID.
	Get(ctx context.Context, uuid string) (TaskModel, error)
	// Count returns the number of tasks.
	Count(ctx context.Context) (int64, error)
	// Query returns the list of tasks with the given offset and limit.
	Query(ctx context.Context, offset, limit int64) ([]TaskModel, int, error)
	// Create saves a new task in the storage.
	Create(ctx context.Context, task TaskModel) error
	// Update updates the task with given UUID in the storage.
	Update(ctx context.Context, task TaskModel) error
	// Delete removes the task with given UUID from the storage.
	Delete(ctx context.Context, uuid string) error
}

// repository persists tasks in database
type repository struct {
	db *db.DB
}

// NewRepository creates a new task repository
func NewRepository(db *db.DB) Repository {
	return repository{db}
}

// Get reads the task with the specified ID from the database.
func (r repository) Get(ctx context.Context, uuid string) (TaskModel, error) {
	var task TaskModel
	err := r.db.With(ctx).Model(&task).
		Relation("Assignee").Relation("Creator").
		Where("t.uuid = ?", uuid).First()

	return task, err
}

// Create saves a new task record in the database.
// It returns the ID of the newly inserted task record.
func (r repository) Create(ctx context.Context, task TaskModel) error {
	_, err := r.db.With(ctx).Model(&task).Insert()
	return err
}

// Update saves the changes to an task in the database.
func (r repository) Update(ctx context.Context, task TaskModel) error {
	_, err := r.db.With(ctx).Model(&task).WherePK().Update()
	return err
}

// Delete deletes an task with the specified ID from the database.
func (r repository) Delete(ctx context.Context, uuid string) error {
	task, err := r.Get(ctx, uuid)
	if err != nil {
		return err
	}
	_, err = r.db.With(ctx).Model(&task).WherePK().Delete()
	return err
}

// Count returns the number of the task records in the database.
func (r repository) Count(ctx context.Context) (int64, error) {
	var count int
	count, err := r.db.With(ctx).Model((*TaskModel)(nil)).Count()
	return int64(count), err
}

// Query retrieves the task records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int64) ([]TaskModel, int, error) {
	var _tasks []TaskModel
	count, err := r.db.With(ctx).Model(&_tasks).
		Relation("Assignee").Relation("Creator").
		Order("id ASC").Limit(int(limit)).Offset(int(offset)).SelectAndCount()
	return _tasks, count, err
}
