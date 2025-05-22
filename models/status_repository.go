package models

import "barcode-generator-be/config"

type StatusRepository interface {
	FindAll() ([]Status, error)
	FindByID(id uint) (*Status, error)
}

type GormStatusRepository struct{}

var StatusRepo StatusRepository = &GormStatusRepository{}

func SetStatusRepository(repo StatusRepository) {
	StatusRepo = repo
}

func (r *GormStatusRepository) FindAll() ([]Status, error) {
	var statuses []Status
	err := config.DB.Find(&statuses).Error
	return statuses, err
}

func (r *GormStatusRepository) FindByID(id uint) (*Status, error) {
	var status Status
	err := config.DB.First(&status, "id = ?", id).Error
	return &status, err
}
