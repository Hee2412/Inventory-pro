package repository

import (
	"Inventory-pro/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindById(id uint) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	FindAll() ([]*domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
	HardDelete(id uint) error

	FindByRole(role string) ([]*domain.User, error)
	FindActiveUsers() ([]*domain.User, error)
	CountByRole(role string) (int64, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Create(user *domain.User) error {
	return u.db.Create(user).Error
}

func (u userRepository) FindById(id uint) (*domain.User, error) {
	var user domain.User
	err := u.db.First(&user, id).Error
	return &user, err
}

func (u userRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := u.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (u userRepository) FindAll() ([]*domain.User, error) {
	var users []*domain.User
	err := u.db.Find(&users).Error
	return users, err
}

func (u userRepository) Update(user *domain.User) error {
	return u.db.Save(user).Error
}

func (u userRepository) Delete(id uint) error {
	return u.db.Delete(&domain.User{}, id).Error
}

func (u userRepository) HardDelete(id uint) error {
	return u.db.Unscoped().Delete(&domain.User{}, id).Error
}

func (u userRepository) FindByRole(role string) ([]*domain.User, error) {
	var user []*domain.User
	err := u.db.Where("role = ?", role).Find(&user).Error
	return user, err
}

func (u userRepository) FindActiveUsers() ([]*domain.User, error) {
	var user []*domain.User
	err := u.db.Where("is_active = true", true).Find(&user).Error
	return user, err
}

func (u userRepository) CountByRole(role string) (int64, error) {
	var count int64
	err := u.db.Table("users").Where("role = ?", role).Count(&count).Error
	return count, err
}
