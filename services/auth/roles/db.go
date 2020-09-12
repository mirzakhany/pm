package roles

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"projectmanager/internal/models"
)

const TableName = "roles"

// Role is db model for single role
type Role struct {
	models.BaseTable
	UUID  string `json:"uuid"`
	Title string `json:"title"`
}

func (Role) TableName() string {
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

func (r *Repository) Create(role RoleRequest) (*Role, error) {

	newItem := &Role{
		UUID:  uuid.New().String(),
		Title: role.Title,
	}

	result := r.DB.Create(newItem)

	if result.Error != nil {
		return nil, result.Error
	}
	return newItem, nil
}

func (r *Repository) Update(uuid string, role RoleRequest) (*Role, error) {
	var item Role
	result := r.DB.Where("uuid = ?", uuid).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}

	result = r.DB.Model(item).Updates(Role{
		Title: role.Title,
	})

	if result.Error != nil {
		return nil, result.Error
	}

	var savedItem Role
	result = r.DB.First(&savedItem, result.RowsAffected)
	if result.Error != nil {
		return nil, result.Error
	}
	return &savedItem, nil
}

func (r *Repository) Retrieve(offset, limit int) ([]Role, int64, error) {
	var items []Role
	var total int64
	result := r.DB.Model(Role{}).Count(&total).Limit(limit).Offset(offset).Find(&items)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return items, total, nil
}

func (r *Repository) GetByUUID(uuid string) (*Role, error) {
	return r.Get("uuid = ?", uuid)
}

func (r *Repository) Get(query string, args ...interface{}) (*Role, error) {
	var item Role
	result := r.DB.Where(query, args).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *Repository) Delete(uuid string) error {
	return r.DB.Where("uuid = ?", uuid).Delete(&Role{}).Error
}
