package models

import (
	"fmt"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	RoleEngineer uint = iota + 1
	RoleManager
	RoleObserver
	RoleAdmin
	RoleSuperAdmin
)

type User struct {
	gorm.Model
	FirstName  string    `gorm:"not null"`
	MiddleName string    `gorm:"not null"`
	LastName   string    `gorm:"not null"`
	Email      string    `gorm:"unique, uniqueIndex, not null"`
	Password   string    `gorm:"not null"`
	Role       uint      `gorm:"not null"` // 1 - engineer, 2 - manager, 3 - observer, 4 - admin, 5 - superadmin
	IsEnabled  bool      `gorm:"default:true"`
	Projects   []Project `gorm:"many2many:user_project;"`
}

func (user *User) Sanitize() {
	user.Email = html.EscapeString(strings.TrimSpace(user.Email))
	user.FirstName = html.EscapeString(strings.TrimSpace(user.FirstName))
	user.LastName = html.EscapeString(strings.TrimSpace(user.LastName))
	user.MiddleName = html.EscapeString(strings.TrimSpace(user.MiddleName))
}

func (user *User) HashPassword() error {
	user.Password = strings.TrimSpace(user.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	user.Password = string(hashedPassword)
	return nil

}

func (user *User) VerifyPassword(password string) error {
	password = strings.TrimSpace(password)
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
