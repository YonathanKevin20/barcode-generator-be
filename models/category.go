package models

import (
	"time"
)

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Code      string    `gorm:"size:4;uniqueIndex;not null" json:"code"`
	Name      string    `gorm:"size:100;not null" json:"name"`
}
