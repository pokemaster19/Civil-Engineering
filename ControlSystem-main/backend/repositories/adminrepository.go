package repositories

import (
	"ControlSystem/models"

	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (r *AdminRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *AdminRepository) DeleteUser(user *models.User) error {
	return r.db.Delete(&user).Error
}
