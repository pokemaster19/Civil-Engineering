package models

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Setup() (*gorm.DB, error) {

	dsn := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Can not connect to the database:", err)
	}

	if err := db.AutoMigrate(
		&User{},
		&Project{},
		&Report{},
		&Defect{},
		&Comment{},
		&Attachment{},
	); err != nil {
		log.Fatal("Auto migration failed:", err)
	}

	var existingUser User
	if err := db.Where("email = ?", "admin@controlsystem.ru").First(&existingUser).Error; err == nil {
		log.Println("Admin user already exists, skipping creation")
	} else if err == gorm.ErrRecordNotFound {
		adminUser := &User{
			FirstName:  "Admin",
			LastName:   "Admin",
			MiddleName: "Admin",
			Email:      "admin@controlsystem.ru",
			Role:       5,
			Password:   os.Getenv("SUPERADMIN_PASSWORD"),
		}

		adminUser.Sanitize()

		if err := adminUser.HashPassword(); err != nil {
			log.Fatalf("Failed to hash admin password: %v", err)
		}
		if err := db.Create(adminUser).Error; err != nil {
			log.Fatalf("Failed to create admin user: %v", err)
		}
		log.Println("Admin user created successfully")
	} else {
		log.Fatalf("Failed to check existing user: %v", err)
	}

	log.Println("Successfully connected")
	return db, nil

}
