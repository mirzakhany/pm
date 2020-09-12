package rules

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"projectmanager/internal/models"
)

const TableName = "rules"

// Rule is db model for single rule
type Rule struct {
	models.BaseTable
	Title string `json:"title"`

}

func (Rule) TableName() string {
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

func (r *Repository) Create(rule RuleRequest) (*Rule, error) {

	newItem := &Rule{
		Title: rule.Title,
	}

	result := r.DB.Create(newItem)
	if result.Error != nil {
		return nil, result.Error
	}
	return newItem, nil
}

func (r *Repository) Update(uuid string, rule RuleRequest) (*Rule, error) {
	var item Rule
	result := r.DB.Where("uuid = ?", uuid).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}

	result = r.DB.Model(item).Updates(Rule{
		Title: rule.Title,
	})

	if result.Error != nil {
		return nil, result.Error
	}

	var savedItem Rule
	result = r.DB.First(&savedItem, item.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &savedItem, nil
}

func (r *Repository) Retrieve(offset, limit int) ([]Rule, int64, error) {
	var items []Rule
	var total int64
	result := r.DB.Model(Rule{}).Count(&total).Limit(limit).Offset(offset).Find(&items)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return items, total, nil
}

func (r *Repository) GetByUUID(uuid string) (*Rule, error) {
	return r.Get("uuid = ?", uuid)
}

func (r *Repository) Get(query string, args ...interface{}) (*Rule, error) {
	var item Rule
	result := r.DB.Where(query, args).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *Repository) Delete(uuid string) error {
	return r.DB.Where("uuid = ?", uuid).Delete(&Rule{}).Error
}
