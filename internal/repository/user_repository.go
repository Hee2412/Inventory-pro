package repository

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/pkg/pagination"
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
	FindAllPaginated(page, limit int) ([]*domain.User, int64, error)
	SearchAndFilter(params request.UserSearchParams) ([]*domain.User, int64, error)
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

func (u *userRepository) FindById(id uint) (*domain.User, error) {
	var user domain.User
	err := u.db.Where("id=?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) FindByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := u.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) FindAll() ([]*domain.User, error) {
	var users []*domain.User
	err := u.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userRepository) Update(user *domain.User) error {
	return u.db.Save(user).Error
}

func (u *userRepository) Delete(id uint) error {
	return u.db.Delete(&domain.User{}, id).Error
}

func (u *userRepository) HardDelete(id uint) error {
	return u.db.Unscoped().Delete(&domain.User{}, id).Error
}

func (u *userRepository) FindByRole(role string) ([]*domain.User, error) {
	var user []*domain.User
	err := u.db.Where("role = ?", role).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) FindActiveUsers() ([]*domain.User, error) {
	var user []*domain.User
	err := u.db.Where("is_active = true").Find(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userRepository) CountByRole(role string) (int64, error) {
	var count int64
	err := u.db.Table("users").Where("role = ?", role).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (u *userRepository) FindAllPaginated(page, limit int) ([]*domain.User, int64, error) {
	var users []*domain.User
	var total int64
	// count total
	if err := u.db.Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	//get paginated data
	err := u.db.Scopes(pagination.Paginate(page, limit)).
		Find(&users).Error

	return users, total, err
}

func (u *userRepository) SearchAndFilter(params request.UserSearchParams) ([]*domain.User, int64, error) {
	var users []*domain.User
	var total int64

	query := u.db.Model(&domain.User{})
	//search by email || store name
	if params.Search != "" {
		searchTerm := "%" + params.Search + "%"
		query = query.Where(
			"email LIKE ? OR store_name LIKE ?",
			searchTerm, searchTerm,
		)
	}
	//filter by role
	if params.Role != "" {
		query = query.Where("role = ?", params.Role)
	}
	//filter by is_active
	if params.IsActive != nil {
		query = query.Where("is_active = ?", params.IsActive)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (params.Page - 1) * params.Limit
	err := query.Offset(offset).Limit(params.Limit).
		Order("created_at DESC").
		Find(&users).Error
	return users, total, err
}
