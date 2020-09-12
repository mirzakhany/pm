package tasks

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"projectmanager/internal/models"
	"projectmanager/services/projects/sprints"
)

const TableName = "tasks"

// Task is db model for single task
type Task struct {
	models.BaseTable
	UUID     string         `json:"uuid"`
	Title    string         `json:"title"`
	SprintID int            `json:"sprintId" gorm:"index"`
	Sprint   sprints.Sprint `json:"sprint" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Estimate string         `json:"estimate"`
	Status   string         `json:"status"`
	Assignee string         `json:"assignee"`
}

func (Task) TableName() string {
	return TableName
}

type Repository struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewRepository(db *gorm.DB) *Repository {
	repo := &Repository{DB: db}
	return repo
}

func (r *Repository) Create(task TaskRequest) (*Task, error) {

	newItem := &Task{
		UUID:     uuid.New().String(),
		Title:    task.Title,
		SprintID: task.SprintID,
		Estimate: task.Estimate,
		Status:   task.Status,
		Assignee: task.Assignee,
	}

	result := r.DB.Create(newItem)
	if result.Error != nil {
		return nil, result.Error
	}

	var savedTask Task
	result = r.DB.Preload(clause.Associations).First(&savedTask, newItem.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &savedTask, nil
}

func (r *Repository) Update(uuid string, task TaskRequest) (*Task, error) {

	var item Task
	result := r.DB.Where("uuid = ?", uuid).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}

	result = r.DB.Model(&item).Updates(Task{
		Title:    task.Title,
		SprintID: task.SprintID,
		Estimate: task.Estimate,
		Status:   task.Status,
		Assignee: task.Assignee,
	})
	if result.Error != nil {
		return nil, result.Error
	}

	var savedTask Task
	result = r.DB.Preload(clause.Associations).First(&savedTask, item.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &savedTask, nil
}

func (r *Repository) FindTasks(query string, args ...interface{}) ([]Task, error) {
	var items []Task
	result := r.DB.Preload(clause.Associations).Find(items, query, args)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *Repository) GetTasks(offset, limit int) ([]Task, int64, error) {
	var items []Task
	var total int64
	result := r.DB.Model(Task{}).Count(&total).Preload(clause.Associations).Limit(limit).Offset(offset).Find(&items)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return items, total, nil
}

func (r *Repository) GetByUUID(uuid string) (*Task, error) {
	return r.Get("uuid = ?", uuid)
}

func (r *Repository) Get(query string, args ...interface{}) (*Task, error) {
	var item Task
	result := r.DB.Preload(clause.Associations).Where(query, args).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *Repository) Delete(uuid string) error {
	return r.DB.Where("uuid = ?", uuid).Delete(&Task{}).Error
}
