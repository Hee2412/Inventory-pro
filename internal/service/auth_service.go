package service

import (
	"Inventory-pro/config"
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/repository"
	"Inventory-pro/pkg/jwt"
	password2 "Inventory-pro/pkg/password"
	"errors"
)

type AuthService interface {
	Login(username string, password string) (string, error)
	Register(creatorRole string, username string, password string, role string, storeName string) error
}
type authService struct {
	repo repository.UserRepository
	cfg  *config.Config
}

func NewAuthService(repo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{repo: repo, cfg: cfg}
}

func (s *authService) Login(username string, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("user not found")
	}

	err = password2.ComparePassword(user.Password, password)
	if err != nil {
		return "", errors.New("invalid password")
	}

	token, err := jwt.GenerateToken(user.ID, user.Role, s.cfg.JWTSecret, s.cfg.JWTExpires)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *authService) Register(
	creatorRole string,
	username string,
	password string,
	storeName string,
	role string) error {
	if creatorRole != "superadmin" && creatorRole != "admin" {
		return errors.New("invalid role")
	}
	user, err := s.repo.FindByUsername(username)
	if err == nil && user != nil {
		return errors.New("username existed")
	}
	hashedPassword, err := password2.HashPassword(password)
	newUser := &domain.User{
		Username:  username,
		Password:  hashedPassword,
		Role:      role,
		StoreName: storeName,
	}
	return s.repo.Create(newUser)
}
