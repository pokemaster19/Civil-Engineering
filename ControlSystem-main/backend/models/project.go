package models

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	Name        string `gorm:"not null; unique"`
	Description string
	Status      uint   `gorm:"not null;default:1"` // 1 - active, 2 - completed, 3 - archived
	Users       []User `gorm:"many2many:user_project;"`
	Defects     []Defect
	Reports     []Report
}

type Report struct {
	gorm.Model
	ProjectID   uint `gorm:"not null"`
	Project     Project
	GeneratedBy uint   `gorm:"not null"`
	Creator     User   `gorm:"foreignKey:GeneratedBy"`
	Title       string `gorm:"not null"`
}
