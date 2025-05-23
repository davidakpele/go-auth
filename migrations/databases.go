package migrations

import (
	"log"
	"api-service/models"
	"gorm.io/gorm"
)

// MigrateModels is an exported function to handle database migrations
func MigrateModels(db *gorm.DB) error {
	log.Println("Starting database migration...")
	err := db.AutoMigrate(
		&models.User{}, 
		&models.AccountVerification{},
	)
	if err == nil {
		log.Println("Database migrated successfully")
	}
	return err
}