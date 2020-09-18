package domains

import (
	"projectmanager/internal/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const TableName = "domains"

// Domain is db model for single domain
type Domain struct {
	models.BaseTable
	UUID    string `json:"uuid" gorm:"index;unique"`
	Title   string `json:"title"`
	Address string `json:"address"  gorm:"index;unique"`
}

func (Domain) TableName() string {
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

func (r *Repository) Create(domain DomainRequest) (*Domain, error) {

	newItem := &Domain{
		UUID:    uuid.New().String(),
		Title:   domain.Title,
		Address: domain.Address,
	}

	result := r.DB.Create(newItem)

	if result.Error != nil {
		return nil, result.Error
	}
	return newItem, nil
}

func (r *Repository) Update(uuid string, domain DomainRequest) (*Domain, error) {
	var item Domain
	result := r.DB.Where("uuid = ?", uuid).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}

	result = r.DB.Model(item).Updates(Domain{
		Title:   domain.Title,
		Address: domain.Address,
	})

	if result.Error != nil {
		return nil, result.Error
	}

	var savedItem Domain
	result = r.DB.First(&savedItem, item.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &savedItem, nil
}

func (r *Repository) Retrieve(offset, limit int) ([]Domain, int64, error) {
	var items []Domain
	var total int64
	result := r.DB.Model(Domain{}).Count(&total).Limit(limit).Offset(offset).Find(&items)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return items, total, nil
}

func (r *Repository) GetByUUID(uuid string) (*Domain, error) {
	return r.Get("uuid = ?", uuid)
}

func (r *Repository) Get(query string, args ...interface{}) (*Domain, error) {
	var item Domain
	result := r.DB.Where(query, args).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *Repository) Delete(uuid string) error {
	var item Domain
	result := r.DB.Where("uuid = ?", uuid).First(&item)
	if result.Error != nil {
		return result.Error
	}
	return r.DB.Where("uuid = ?", uuid).Delete(&Domain{}).Error
}
