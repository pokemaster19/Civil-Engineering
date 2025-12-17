package models

import (
	"time"

	"gorm.io/gorm"
)

type Defect struct {
	gorm.Model
	ProjectID   uint `gorm:"not null"`
	Project     Project
	Title       string `gorm:"not null"`
	Description string
	Priority    uint `gorm:"not null"` // 1 - low, 2 - medium, 3 - high
	Status      uint `gorm:"not null"` // 1 - open, 2 - in progress, 3 - resolved, 4 - overdue
	AssignedTo  uint `gorm:"index"`
	Assignee    User `gorm:"foreignKey:AssignedTo"`
	CreatedBy   uint `gorm:"not null, index"`
	Creator     User `gorm:"foreignKey:CreatedBy"`
	Deadline    *time.Time
}

type Comment struct {
	gorm.Model
	DefectID  uint `gorm:"not null, index"`
	Defect    Defect
	Content   string `gorm:"not null"`
	CreatedBy uint   `gorm:"not null"`
	Creator   User   `gorm:"foreignKey:CreatedBy"`
}
