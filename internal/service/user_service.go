package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"Inventory-pro/pkg/password"
	"errors"
)

type UserService interface {
	GetAllUsers() ([]response.UserResponse, error)
	GetUserById(userId uint) (*response.UserResponse, error)
	UpdateUser(userId uint, req request.UpdateUserRequest) error
	DeactivateUser(userId uint) error
	ActivateUser(userId uint) error

	DeleteUser(userId uint) error
	HardDeleteUser(userId uint) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}
func toUserResponse(user *domain.User) response.UserResponse {
	return response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		StoreName: user.StoreName,
		StoreCode: user.StoreCode,
		IsActive:  user.IsActive,
		CreateAt:  user.CreatedAt,
		LastLogin: user.LastLogin,
	}
}
func (s *userService) GetAllUsers() ([]response.UserResponse, error) {
	user, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	var result []response.UserResponse
	for _, user := range user {
		result = append(result, toUserResponse(user))
	}
	return result, nil
}
func (s *userService) GetUserById(userId uint) (*response.UserResponse, error) {
	user, err := s.repo.FindById(userId)
	if err != nil {
		return nil, err
	}
	result := toUserResponse(user)
	return &result, nil
}
func (s *userService) UpdateUser(userId uint, req request.UpdateUserRequest) error {
	user, err := s.repo.FindById(userId)
	if err != nil {
		return errors.New("user not found")
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
			return err
		}
		user.Password = hashedPassword
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	return s.repo.Update(user)
}
func (s *userService) DeactivateUser(userId uint) error {
	user, err := s.repo.FindById(userId)
	if err != nil {
		return errors.New("user not found")
	}
	if user.IsActive != true {
		return errors.New("user is not active")
	}
	user.IsActive = false
	return s.repo.Update(user)
}

func (s *userService) ActivateUser(userId uint) error {
	user, err := s.repo.FindById(userId)
	if err != nil {
		return errors.New("user not found")
	}
	if user.IsActive == true {
		return errors.New("user is active")
	}
	user.IsActive = true
	return s.repo.Update(user)
}

func (s *userService) DeleteUser(userId uint) error {
	_, err := s.repo.FindById(userId)
	if err != nil {
		return errors.New("user not found")
	}
	return s.repo.Delete(userId)
}
func (s *userService) HardDeleteUser(userId uint) error {
	_, err := s.repo.FindById(userId)
	if err != nil {
		return errors.New("user not found")
	}
	return s.repo.HardDelete(userId)
}
