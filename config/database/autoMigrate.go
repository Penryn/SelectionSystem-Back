package database

import (
	"SelectionSystem-Back/app/models"

	"gorm.io/gorm"
)

func autoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Teacher{},
		&models.Student{},
		&models.Reason{},
		&models.DDL{},
		&models.Conversation{},
	)

	return err
}
