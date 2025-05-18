package models

import "barcode-generator-be/config"

type UserRepository interface {
	FindAll() ([]User, error)
	FindAllWithFilter(filter *UserFilter) ([]User, int64, error)
	FindByID(id uint) (*User, error)
	FindByUsername(username string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	Delete(user *User) error
}

type UserFilter struct {
	Username string
	Role     string
	Offset   int
	Limit    int
}

type GormUserRepository struct{}

func (r *GormUserRepository) FindAll() ([]User, error) {
	var users []User
	err := config.DB.Find(&users).Error
	return users, err
}

func (r *GormUserRepository) FindAllWithFilter(filter *UserFilter) ([]User, int64, error) {
	var users []User
	query := config.DB.Model(&User{})
	if filter.Username != "" {
		query = query.Where("username ILIKE ?", "%"+filter.Username+"%")
	}
	if filter.Role != "" && filter.Role != "all" {
		query = query.Where("role = ?", filter.Role)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(filter.Offset).Limit(filter.Limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *GormUserRepository) FindByID(id uint) (*User, error) {
	var user User
	err := config.DB.First(&user, "id = ?", id).Error
	return &user, err
}

func (r *GormUserRepository) FindByUsername(username string) (*User, error) {
	var user User
	err := config.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *GormUserRepository) Create(user *User) error {
	return config.DB.Create(user).Error
}

func (r *GormUserRepository) Update(user *User) error {
	return config.DB.Model(user).Updates(user).Error
}

func (r *GormUserRepository) Delete(user *User) error {
	return config.DB.Delete(user).Error
}
