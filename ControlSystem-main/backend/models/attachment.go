package models

import "gorm.io/gorm"

type Attachment struct {
	gorm.Model
	DefectID   *uint `gorm:"index"`
	Defect     Defect
	ProjectID  *uint `gorm:"index"`
	Project    Project
	FileName   string `gorm:"not null"`
	FilePath   string `gorm:"not null"`
	FileType   string `gorm:"not null"`
	UploadedBy uint   `gorm:"not null"`
	Uploader   User   `gorm:"foreignKey:UploadedBy"`
}
