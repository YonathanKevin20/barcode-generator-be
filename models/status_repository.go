package models

import "barcode-generator-be/config"

type StatusRepository interface {
	FindAll() ([]Status, error)
}

type GormStatusRepository struct{}

func (r *GormStatusRepository) FindAll() ([]Status, error) {
	var statuses []Status
	err := config.DB.Find(&statuses).Error
	return statuses, err
}
