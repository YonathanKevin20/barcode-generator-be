package models

import (
	"barcode-generator-be/config"
	"fmt"
	"strconv"
	"time"
)

type BarcodeRepository interface {
	FindAll() ([]Barcode, error)
	FindAllWithFilter(filter *BarcodeFilter) ([]BarcodeResult, int64, error)
	FindByID(id uint) (*BarcodeResult, error)
	FindActiveByIDs(ids []uint) ([]BarcodeResult, error)
	Create(barcode *Barcode) error
	Update(barcode *Barcode) error
	UpdateInactive(id uint, isInactive bool) error
	Delete(id uint) error
	GetNextProductCode(categoryID, supplierID uint) (string, error)
}

type BarcodeFilter struct {
	StatusID     uint
	CategoryID   uint
	CategoryName string
	SupplierID   uint
	SupplierName string
	ProductName  string
	Barcode      string
	Offset       int
	Limit        int
}

type BarcodeResult struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedByUser *string   `json:"created_by_user"`
	StatusName    string    `json:"status_name"`
	CategoryName  string    `json:"category_name"`
	SupplierName  string    `json:"supplier_name"`
	ProductName   string    `json:"product_name"`
	Barcode       string    `json:"barcode"`
	IsInactive    bool      `json:"is_inactive"`
}

type GormBarcodeRepository struct{}

var BarcodeRepo BarcodeRepository = &GormBarcodeRepository{}

func SetBarcodeRepository(repo BarcodeRepository) {
	BarcodeRepo = repo
}

func (r *GormBarcodeRepository) FindAll() ([]Barcode, error) {
	var barcodes []Barcode
	err := config.DB.Order("id DESC").Find(&barcodes).Error
	return barcodes, err
}

func (r *GormBarcodeRepository) FindAllWithFilter(filter *BarcodeFilter) ([]BarcodeResult, int64, error) {
	query := config.DB.Model(&Barcode{})
	if filter.StatusID != 0 {
		query = query.Where("status_id = ?", filter.StatusID)
	}
	if filter.CategoryID != 0 {
		query = query.Where("category_id = ?", filter.CategoryID)
	}
	if filter.CategoryName != "" {
		query = query.Joins("INNER JOIN categories ON barcodes.category_id = categories.id").
			Where("categories.name LIKE ?", "%"+filter.CategoryName+"%")
	}
	if filter.SupplierID != 0 {
		query = query.Where("supplier_id = ?", filter.SupplierID)
	}
	if filter.SupplierName != "" {
		query = query.Joins("INNER JOIN suppliers ON barcodes.supplier_id = suppliers.id").
			Where("suppliers.name LIKE ?", "%"+filter.SupplierName+"%")
	}
	if filter.ProductName != "" {
		query = query.Where("product_name LIKE ?", "%"+filter.ProductName+"%")
	}
	if filter.Barcode != "" {
		query = query.Where("barcode LIKE ?", "%"+filter.Barcode+"%")
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db := config.DB.Table("barcodes").
		Select("barcodes.id, barcodes.created_at, users.username AS created_by_user, statuses.name AS status_name, categories.name AS category_name, suppliers.name AS supplier_name, product_name, barcodes.barcode, barcodes.is_inactive").
		Joins("LEFT JOIN users ON barcodes.created_by = users.id").
		Joins("INNER JOIN statuses ON barcodes.status_id = statuses.id").
		Joins("INNER JOIN categories ON barcodes.category_id = categories.id").
		Joins("INNER JOIN suppliers ON barcodes.supplier_id = suppliers.id")
	if filter.StatusID != 0 {
		db = db.Where("barcodes.status_id = ?", filter.StatusID)
	}
	if filter.CategoryID != 0 {
		db = db.Where("barcodes.category_id = ?", filter.CategoryID)
	}
	if filter.CategoryName != "" {
		db = db.Where("categories.name LIKE ?", "%"+filter.CategoryName+"%")
	}
	if filter.SupplierID != 0 {
		db = db.Where("barcodes.supplier_id = ?", filter.SupplierID)
	}
	if filter.SupplierName != "" {
		db = db.Where("suppliers.name LIKE ?", "%"+filter.SupplierName+"%")
	}
	if filter.ProductName != "" {
		db = db.Where("barcodes.product_name LIKE ?", "%"+filter.ProductName+"%")
	}
	if filter.Barcode != "" {
		db = db.Where("barcodes.barcode LIKE ?", "%"+filter.Barcode+"%")
	}
	var results []BarcodeResult
	if err := db.Order("barcodes.created_at DESC, barcodes.id DESC").Offset(filter.Offset).Limit(filter.Limit).Find(&results).Error; err != nil {
		return nil, 0, err
	}
	return results, total, nil
}

func (r *GormBarcodeRepository) FindByID(id uint) (*BarcodeResult, error) {
	var result BarcodeResult
	err := config.DB.Table("barcodes").
		Select("barcodes.id, barcodes.created_at, users.username AS created_by_user, statuses.name AS status_name, categories.name AS category_name, suppliers.name AS supplier_name, product_name, barcodes.barcode, barcodes.is_inactive").
		Joins("LEFT JOIN users ON barcodes.created_by = users.id").
		Joins("INNER JOIN statuses ON barcodes.status_id = statuses.id").
		Joins("INNER JOIN categories ON barcodes.category_id = categories.id").
		Joins("INNER JOIN suppliers ON barcodes.supplier_id = suppliers.id").
		Where("barcodes.id = ?", id).
		First(&result).Error
	return &result, err
}

func (r *GormBarcodeRepository) FindActiveByIDs(ids []uint) ([]BarcodeResult, error) {
	var results []BarcodeResult
	err := config.DB.Table("barcodes").
		Select("barcodes.id, barcodes.created_at, users.username AS created_by_user, statuses.name AS status_name, categories.name AS category_name, suppliers.name AS supplier_name, product_name, barcodes.barcode, barcodes.is_inactive").
		Joins("LEFT JOIN users ON barcodes.created_by = users.id").
		Joins("INNER JOIN statuses ON barcodes.status_id = statuses.id").
		Joins("INNER JOIN categories ON barcodes.category_id = categories.id").
		Joins("INNER JOIN suppliers ON barcodes.supplier_id = suppliers.id").
		Where("barcodes.id IN (?) AND barcodes.is_inactive = false", ids).
		Order("barcodes.created_at DESC, barcodes.id DESC").
		Find(&results).Error
	return results, err
}

func (r *GormBarcodeRepository) Create(barcode *Barcode) error {
	return config.DB.Create(barcode).Error
}

func (r *GormBarcodeRepository) Update(barcode *Barcode) error {
	return config.DB.Model(barcode).Updates(barcode).Error
}

func (r *GormBarcodeRepository) UpdateInactive(id uint, isInactive bool) error {
	return config.DB.Model(&Barcode{}).Where("id = ?", id).Update("is_inactive", isInactive).Error
}

func (r *GormBarcodeRepository) Delete(id uint) error {
	return config.DB.Delete(&Barcode{}, id).Error
}

func (r *GormBarcodeRepository) GetNextProductCode(categoryID, supplierID uint) (string, error) {
	var maxProductCode string
	err := config.DB.Table("barcodes").
		Select("SUBSTRING(barcode, 10, 3) AS product_code").
		Where("category_id = ? AND supplier_id = ?", categoryID, supplierID).
		Order("product_code DESC").
		Limit(1).
		Scan(&maxProductCode).Error
	if err != nil {
		return "", err
	}
	nextProductCode := "001"
	if maxProductCode != "" {
		if n, err := strconv.Atoi(maxProductCode); err == nil {
			nextProductCode = fmt.Sprintf("%03d", n+1)
		}
	}
	return nextProductCode, nil
}
