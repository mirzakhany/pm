package sprints

import (
	"projectmanager/internal/models"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const TableName = "sprints"

// Sprint is db model for single sprint
type Sprint struct {
	models.BaseTable
	UUID    string    `json:"uuid"`
	Title   string    `json:"title"`
	Status  string    `json:"status"`
	StartAt time.Time `json:"startAt"`
	EndAt   time.Time `json:"endAt"`
}

func (Sprint) TableName() string {
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

func (r *Repository) Create(sprint SprintRequest) (*Sprint, error) {

	newItem := &Sprint{
		UUID:    uuid.New().String(),
		Title:   sprint.Title,
		Status:  sprint.Status,
		StartAt: sprint.StartAt,
		EndAt:   sprint.EndAt,
	}

	result := r.DB.Create(newItem)

	if result.Error != nil {
		return nil, result.Error
	}
	return newItem, nil
}

func (r *Repository) Update(uuid string, sprint SprintRequest) (*Sprint, error) {
	var item Sprint
	result := r.DB.Where("uuid = ?", uuid).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}

	result = r.DB.Model(item).Updates(Sprint{
		Title:   sprint.Title,
		Status:  sprint.Status,
		StartAt: sprint.StartAt,
		EndAt:   sprint.EndAt,
	})
	if result.Error != nil {
		return nil, result.Error
	}

	var savedItem Sprint
	result = r.DB.First(&savedItem, item.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &savedItem, nil
}

func (r *Repository) FindSprints(query string, args ...interface{}) ([]Sprint, error) {
	var items []Sprint
	result := r.DB.Find(&items, query, args)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *Repository) Retrieve(offset, limit int) ([]Sprint, int64, error) {
	var items []Sprint
	var total int64
	result := r.DB.Model(Sprint{}).Count(&total).Limit(limit).Offset(offset).Find(&items)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return items, total, nil
}

func (r *Repository) GetByUUID(uuid string) (*Sprint, error) {
	return r.Get("uuid = ?", uuid)
}

func (r *Repository) Get(query string, args ...interface{}) (*Sprint, error) {
	var item Sprint
	result := r.DB.Where(query, args).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *Repository) DeleteSprint(uuid string) error {
	return r.DB.Where("uuid = ?", uuid).Delete(&Sprint{}).Error
}
