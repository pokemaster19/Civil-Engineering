package repositories

import (
	"ControlSystem/models"

	"gorm.io/gorm"
)

type DefectRepository struct {
	db *gorm.DB
}

func NewDefectRepository(db *gorm.DB) *DefectRepository {
	return &DefectRepository{db: db}
}

func (r *DefectRepository) GetByProjectID(projectID uint) ([]models.Defect, error) {
	var defects []models.Defect
	err := r.db.Where("project_id = ?", projectID).Find(&defects).Error
	return defects, err
}

func (r *DefectRepository) GetByID(id uint) (*models.Defect, error) {
	var defect models.Defect
	err := r.db.First(&defect, id).Error
	return &defect, err
}
