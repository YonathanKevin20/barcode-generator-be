package models

import (
	"time"
)

type Barcode struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"-"`
	StatusID    uint      `gorm:"not null;index" json:"status_id"`
	Status      *Status   `gorm:"foreignKey:StatusID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"status,omitempty"`
	CategoryID  uint      `gorm:"not null;index" json:"category_id"`
	Category    *Category `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"category,omitempty"`
	SupplierID  uint      `gorm:"not null;index" json:"supplier_id"`
	Supplier    *Supplier `gorm:"foreignKey:SupplierID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"supplier,omitempty"`
	ProductName string    `gorm:"size:100;not null;index" json:"product_name"`
	Barcode     string    `gorm:"size:13;uniqueIndex;not null" json:"barcode"`
}
