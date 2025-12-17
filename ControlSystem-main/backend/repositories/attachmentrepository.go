package repositories

import (
	"ControlSystem/models"

	"gorm.io/gorm"
)

type AttachmentRepository struct {
	db *gorm.DB
}

func NewAttachmentRepository(db *gorm.DB) *AttachmentRepository {
	return &AttachmentRepository{db: db}
}

func (r *AttachmentRepository) GetByProjectIDs(projectIDs []uint) ([]models.Attachment, error) {
	var attachments []models.Attachment
	err := r.db.Where("project_id IN ?", projectIDs).Find(&attachments).Error
	return attachments, err
}

func (r *AttachmentRepository) GetByDefectIDs(defectIDs []uint) ([]models.Attachment, error) {
	var attachments []models.Attachment
	err := r.db.Where("defect_id IN ?", defectIDs).Find(&attachments).Error
	return attachments, err
}

func (r *AttachmentRepository) Create(attachment *models.Attachment) error {
	return r.db.Create(attachment).Error
}
