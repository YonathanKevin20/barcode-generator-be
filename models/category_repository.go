package models

import "barcode-generator-be/config"

type CategoryRepository interface {
	FindAll() ([]Category, error)
	FindAllWithFilter(filter *CategoryFilter) ([]Category, int64, error)
	FindByID(id uint) (*Category, error)
	FindByCode(code string) (*Category, error)
	FindByCodeAndName(code, name string) (*Category, error)
	Create(category *Category) error
	Update(category *Category) error
	Delete(category *Category) error
}

type CategoryFilter struct {
	Code   string
	Name   string
	Offset int
	Limit  int
}

type GormCategoryRepository struct{}

func (r *GormCategoryRepository) FindAll() ([]Category, error) {
	var categories []Category
	err := config.DB.Order("name ASC").Find(&categories).Error
	return categories, err
}

func (r *GormCategoryRepository) FindAllWithFilter(filter *CategoryFilter) ([]Category, int64, error) {
	var categories []Category
	query := config.DB.Model(&Category{})
	if filter.Code != "" {
		query = query.Where("code = ?", filter.Code)
	}
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("name ASC").Offset(filter.Offset).Limit(filter.Limit).Find(&categories).Error; err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

func (r *GormCategoryRepository) FindByID(id uint) (*Category, error) {
	var category Category
	err := config.DB.First(&category, "id = ?", id).Error
	return &category, err
}

func (r *GormCategoryRepository) FindByCode(code string) (*Category, error) {
	var category Category
	err := config.DB.First(&category, "code = ?", code).Error
	return &category, err
}

func (r *GormCategoryRepository) FindByCodeAndName(code, name string) (*Category, error) {
	var category Category
	err := config.DB.First(&category, "code = ? AND name = ?", code, name).Error
	return &category, err
}

func (r *GormCategoryRepository) Create(category *Category) error {
	return config.DB.Create(category).Error
}

func (r *GormCategoryRepository) Update(category *Category) error {
	return config.DB.Model(category).Updates(category).Error
}

func (r *GormCategoryRepository) Delete(category *Category) error {
	return config.DB.Delete(category).Error
}
