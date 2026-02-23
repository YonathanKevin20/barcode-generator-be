package migrations

import (
	"barcode-generator-be/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	// Run manual migrations
	if err := createUserRoleEnum(db); err != nil {
		return err
	}

	// Run AutoMigrate for all models
	err := db.AutoMigrate(
		&models.User{},
		&models.Status{},
		&models.Category{},
		&models.Supplier{},
		&models.Barcode{},
		// Add any other models here
	)
	if err != nil {
		return err
	}

	// Run seeders
	if err := StatusSeeder(db); err != nil {
		return err
	}

	return nil
}

func createUserRoleEnum(db *gorm.DB) error {
	var exists bool
	err := db.Raw("SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role')").Scan(&exists).Error
	if err != nil {
		return err
	}

	if !exists {
		if err := db.Exec("CREATE TYPE user_role AS ENUM ('admin', 'operator')").Error; err != nil {
			return err
		}
	}

	return nil
}
