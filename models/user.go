package models

import (
	"database/sql/driver"
	"errors"
	"time"
)

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleOperator Role = "operator"
)

func (r *Role) Scan(value any) error {
	strValue, ok := value.(string)
	if !ok {
		return errors.New("invalid role value")
	}
	*r = Role(strValue)
	return nil
}

func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	Role      Role      `gorm:"type:user_role;default:operator;not null" json:"role"`
}
