package users

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"projectmanager/internal/models"
)

const TableName = "users"

// User is db model for single user
type User struct {
	models.BaseTable
	UUID     string        `json:"uuid"`
	Username string        `json:"username"`
	Password string        `json:"password"`
	Enable   bool          `json:"enable"`
	Email    string        `json:"email"`
}

func (User) TableName() string {
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

func (r *Repository) Create(user UserRequest) (*User, error) {

	newItem := &User{
		UUID:     uuid.New().String(),
		Username: user.Username,
		Password: user.Password,
		Enable:   user.Enable,
	}

	result := r.DB.Create(newItem)

	if result.Error != nil {
		return nil, result.Error
	}
	return newItem, nil
}

func (r *Repository) Update(uuid string, user UserRequest) (*User, error) {
	var item User
	result := r.DB.Where("uuid = ?", uuid).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}

	result = r.DB.Model(item).Updates(User{
		Username: user.Username,
		Password: user.Password,
		Enable:   user.Enable,
	})

	if result.Error != nil {
		return nil, result.Error
	}

	var savedItem User
	result = r.DB.First(&savedItem, item.ID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &savedItem, nil
}

func (r *Repository) FindSprints(query string, args ...interface{}) ([]User, error) {
	var items []User
	result := r.DB.Find(&items, query, args)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *Repository) Retrieve(offset, limit int) ([]User, int64, error) {
	var items []User
	var total int64
	result := r.DB.Model(User{}).Count(&total).Limit(limit).Offset(offset).Find(&items)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return items, total, nil
}

func (r *Repository) GetByUUID(uuid string) (*User, error) {
	return r.Get("uuid = ?", uuid)
}

func (r *Repository) Get(query string, args ...interface{}) (*User, error) {
	var item User
	result := r.DB.Where(query, args).First(&item)
	if result.Error != nil {
		return nil, result.Error
	}
	return &item, nil
}

func (r *Repository) Delete(uuid string) error {
	return r.DB.Where("uuid = ?", uuid).Delete(&User{}).Error
}
