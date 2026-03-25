package service

import (
	"Inventory-pro/config"
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/repository"
	"Inventory-pro/pkg/jwt"
	password2 "Inventory-pro/pkg/password"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type AuthService interface {
	Login(username string, password string) (string, *domain.User, error)
	Register(creatorID *uint, creatorRole string, req request.RegisterRequest) error
	GetProfile(userId uint) (*domain.User, error)
}
type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(repo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{userRepo: repo, cfg: cfg}
}

func (s *authService) Login(username string, password string) (string, *domain.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, domain.ErrUserNotFound
		}
		return "", nil, fmt.Errorf("%w: context: %v", domain.ErrDatabase, err)
	}
	if !user.IsActive {
		return "", nil, fmt.Errorf("%w: user is inactive %v", domain.ErrUserInactive, user.ID)
	}
	if !password2.ComparePassword(user.Password, password) {
		return "", nil, fmt.Errorf("%w: username/password invalid", domain.ErrInvalidCredentials)
	}
	token, err := jwt.GenerateToken(user.ID, user.Role, s.cfg.JWTSecret, s.cfg.JWTExpires)
	if err != nil {
		return "", nil, fmt.Errorf("%w: failed to generate token: %v", domain.ErrInternalServer, err)
	}
	lastLogin := time.Now()
	user.LastLogin = &lastLogin
	err = s.userRepo.Update(user)
	if err != nil {
		return "", nil, fmt.Errorf("%w: failed to update user: %v", domain.ErrInternalServer, err)
	}
	return token, user, nil
}

func (s *authService) Register(creatorID *uint, creatorRole string, req request.RegisterRequest) error {
	if existingUser, _ := s.userRepo.FindByUsername(req.Username); existingUser != nil {
		return fmt.Errorf("%w: username exists: %v", domain.ErrAlreadyExists, req.Username)
	}
	hashedPassword, err := password2.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("%w: failed to hash password: %v", domain.ErrInternalServer, err)
	}

	newUser := &domain.User{
		Username:  req.Username,
		Password:  hashedPassword,
		CreatedBy: creatorID,
	}

	switch req.Role {
	case "admin":
		if creatorRole != "super_admin" {
			return fmt.Errorf("%w: invalid role", domain.ErrForbidden)
		}
		newUser.Role = "admin"
	case "store":
		if creatorRole != "super_admin" && creatorRole != "admin" {
			return fmt.Errorf("%w: invalid role", domain.ErrForbidden)
		}
		if req.StoreName == "" {
			return fmt.Errorf("%w: invalid store name", domain.ErrInvalidInput)
		}
		newUser.Role = "store"
		newUser.StoreName = req.StoreName
	default:
		return fmt.Errorf("%w: invalid role", domain.ErrForbidden)
	}
	err = s.userRepo.Create(newUser)
	if err != nil {
		return fmt.Errorf("%w: failed to create user: %v", domain.ErrInternalServer, err)
	}
	return nil
}

func (s *authService) GetProfile(userId uint) (*domain.User, error) {
	user, err := s.userRepo.FindById(userId)
	if err != nil {
		return nil, fmt.Errorf("%w: user not found %v", domain.ErrNotFound, userId)
	}
	if !user.IsActive {
		return nil, fmt.Errorf("%w: user is inactive %v", domain.ErrUserInactive, userId)
	}
	return user, nil
}
