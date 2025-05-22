package models

import "barcode-generator-be/config"

type SupplierRepository interface {
	FindAll() ([]Supplier, error)
	FindAllWithFilter(filter *SupplierFilter) ([]Supplier, int64, error)
	FindByID(id uint) (*Supplier, error)
	FindByCodeAndName(code, name string) (*Supplier, error)
	Create(supplier *Supplier) error
	Update(supplier *Supplier) error
	Delete(supplier *Supplier) error
}

type SupplierFilter struct {
	Code   string
	Name   string
	Offset int
	Limit  int
}

type GormSupplierRepository struct{}

var SupplierRepo SupplierRepository = &GormSupplierRepository{}

func SetSupplierRepository(repo SupplierRepository) {
	SupplierRepo = repo
}

func (r *GormSupplierRepository) FindAll() ([]Supplier, error) {
	var suppliers []Supplier
	err := config.DB.Order("name ASC").Find(&suppliers).Error
	return suppliers, err
}

func (r *GormSupplierRepository) FindAllWithFilter(filter *SupplierFilter) ([]Supplier, int64, error) {
	var suppliers []Supplier
	query := config.DB.Model(&Supplier{})
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
	if err := query.Order("name ASC").Offset(filter.Offset).Limit(filter.Limit).Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}
	return suppliers, total, nil
}

func (r *GormSupplierRepository) FindByID(id uint) (*Supplier, error) {
	var supplier Supplier
	err := config.DB.First(&supplier, "id = ?", id).Error
	return &supplier, err
}

func (r *GormSupplierRepository) FindByCodeAndName(code, name string) (*Supplier, error) {
	var supplier Supplier
	err := config.DB.First(&supplier, "code = ? AND name = ?", code, name).Error
	return &supplier, err
}

func (r *GormSupplierRepository) Create(supplier *Supplier) error {
	return config.DB.Create(supplier).Error
}

func (r *GormSupplierRepository) Update(supplier *Supplier) error {
	return config.DB.Model(supplier).Updates(supplier).Error
}

func (r *GormSupplierRepository) Delete(supplier *Supplier) error {
	return config.DB.Delete(supplier).Error
}
