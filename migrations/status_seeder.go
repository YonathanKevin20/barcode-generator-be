package migrations

import (
	"barcode-generator-be/models"

	"gorm.io/gorm"
)

func StatusSeeder(db *gorm.DB) error {
	statuses := []string{"In Stock", "Konsinyasi"}
	for _, name := range statuses {
		var count int64
		db.Model(&models.Status{}).Where("name = ?", name).Count(&count)
		if count == 0 {
			if err := db.Create(&models.Status{Name: name}).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
