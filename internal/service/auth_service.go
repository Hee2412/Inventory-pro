package service

import (
	"Inventory-pro/config"
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/repository"
	"Inventory-pro/pkg/jwt"
	password2 "Inventory-pro/pkg/password"
	"errors"
)

type AuthService interface {
	Login(username string, password string) (string, *domain.User, error)
	Register(creatorRole string, req request.RegisterRequest) error
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
		return "", nil, errors.New("user not found")
	}
	if !user.IsActive {
		return "", nil, errors.New("your account is inactive")
	}
	if !password2.ComparePassword(user.Password, password) {
		return "", nil, errors.New("invalid username or password")
	}
	token, err := jwt.GenerateToken(user.ID, user.Role, s.cfg.JWTSecret, s.cfg.JWTExpires)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *authService) Register(creatorRole string, req request.RegisterRequest) error {
	if existingUser, _ := s.userRepo.FindByUsername(req.Username); existingUser != nil {
		return errors.New("user already exists")
	}
	hashedPassword, err := password2.HashPassword(req.Password)
	if err != nil {
		return err
	}
	newUser := &domain.User{
		Username: req.Username,
		Password: hashedPassword,
	}

	switch req.Role {
	case "admin":
		if creatorRole != "super_admin" {
			return errors.New("invalid role")
		}
		newUser.Role = "admin"
	case "store":
		if creatorRole != "super_admin" && creatorRole != "admin" {
			return errors.New("invalid role")
		}
		if req.StoreName == "" {
			return errors.New("invalid store name")
		}
		newUser.Role = "store"
		newUser.StoreName = req.StoreName
	default:
		return errors.New("invalid role")
	}
	return s.userRepo.Create(newUser)
}

// GetProfile GET /api/me
func (s *authService) GetProfile(userId uint) (*domain.User, error) {
	user, err := s.userRepo.FindById(userId)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if !user.IsActive {
		return nil, errors.New("your account is inactive")
	}
	return user, nil
}
