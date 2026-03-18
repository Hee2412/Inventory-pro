package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"Inventory-pro/pkg/password"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type UserService interface {
	GetAllUsers() ([]*response.UserResponse, error)
	GetUserById(userId uint) (*response.UserResponse, error)
	UpdateUser(userId uint, req request.UpdateUserRequest) error
	DeactivateUser(userId uint) error
	ActivateUser(userId uint) error

	DeleteUser(userId uint) error
	HardDeleteUser(userId uint) error
	SearchAndFilter(params request.UserSearchParams) ([]*response.UserResponse, int64, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{userRepo: repo}
}

func (s *userService) GetAllUsers() ([]*response.UserResponse, error) {
	users, err := s.userRepo.FindAll()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: failed to get user: %v", domain.ErrDatabase, err)
	}
	var result []*response.UserResponse
	for _, user := range users {
		result = append(result, toUserResponse(user))
	}
	return result, nil
}
func (s *userService) GetUserById(userId uint) (*response.UserResponse, error) {
	user, err := s.userRepo.FindById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: failed to get user: %v", domain.ErrDatabase, err)
	}
	result := toUserResponse(user)
	return result, nil
}
func (s *userService) UpdateUser(userId uint, req request.UpdateUserRequest) error {
	user, err := s.userRepo.FindById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("%w: failed to find user: %v", domain.ErrDatabase, err)
	}
	if req.StoreName != nil {
		user.StoreName = *req.StoreName
	}
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Password != nil {
		hashedPassword, err := password.HashPassword(*req.Password)
		if err != nil {
			return fmt.Errorf("%w: failed to hash password: %v", domain.ErrInternalServer, err)
		}
		user.Password = hashedPassword
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("%w: failed to update user: %v", domain.ErrDatabase, err)
	}
	return nil
}
func (s *userService) DeactivateUser(userId uint) error {
	user, err := s.userRepo.FindById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("%w: failed to find user", domain.ErrDatabase)
	}
	if user.IsActive != true {
		return fmt.Errorf("%w: user is inactive", domain.ErrUserInactive)
	}
	user.IsActive = false
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("%w: failed to update user: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (s *userService) ActivateUser(userId uint) error {
	user, err := s.userRepo.FindById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("%w: failed to find user", domain.ErrDatabase)
	}
	if user.IsActive == true {
		return fmt.Errorf("%w: user is already active", domain.ErrInternalServer)
	}
	user.IsActive = true
	return s.userRepo.Update(user)
}

func (s *userService) DeleteUser(userId uint) error {
	user, err := s.userRepo.FindById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("%w: failed to find user", domain.ErrDatabase)
	}
	if user.DeletedAt.Valid == true {
		return fmt.Errorf("%v: user is already deleted", domain.ErrInvalidInput)
	}
	user.IsActive = false
	err = s.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("%w: failed to update user: %v", domain.ErrDatabase, err)
	}
	err = s.userRepo.Delete(userId)
	if err != nil {
		return fmt.Errorf("%w: failed to delete user: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (s *userService) HardDeleteUser(userId uint) error {
	user, err := s.userRepo.FindById(userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("%w: failed to find user: %v", domain.ErrDatabase, err)
	}
	if user.DeletedAt.Valid == false {
		return fmt.Errorf("%v: soft delete user first", domain.ErrInvalidInput)
	}
	err = s.userRepo.HardDelete(userId)
	if err != nil {
		return fmt.Errorf("%w: failed to delete user: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (s *userService) SearchAndFilter(params request.UserSearchParams) ([]*response.UserResponse, int64, error) {
	users, total, err := s.userRepo.SearchAndFilter(params)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: failed to fetch user data %v", domain.ErrDatabase, err)
	}
	result := make([]*response.UserResponse, 0, len(users))
	for _, user := range users {
		result = append(result, toUserResponse(user))
	}
	return result, total, nil
}
