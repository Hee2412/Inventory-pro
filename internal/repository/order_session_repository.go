package repository

import (
	"Inventory-pro/internal/domain"
	"gorm.io/gorm"
)

type OrderSessionRepository interface {
	Create(orderSession *domain.OrderSession) error
	FindById(id uint) (*domain.OrderSession, error)
	FindAll() ([]*domain.OrderSession, error)
	FindByStatus(status string) ([]*domain.OrderSession, error)
	Update(orderSession *domain.OrderSession) error
	Delete(id uint) error
}

type orderSessionRepository struct {
	db *gorm.DB
}

func NewOrderSessionRepository(db *gorm.DB) OrderSessionRepository {
	return &orderSessionRepository{db: db}
}

func (o *orderSessionRepository) Create(orderSession *domain.OrderSession) error {
	return o.db.Create(orderSession).Error
}

func (o *orderSessionRepository) FindById(id uint) (*domain.OrderSession, error) {
	var orderSession domain.OrderSession
	err := o.db.Where("id = ?", id).First(&orderSession).Error
	return &orderSession, err
}

func (o *orderSessionRepository) FindAll() ([]*domain.OrderSession, error) {
	var orderSessions []*domain.OrderSession
	err := o.db.Find(&orderSessions).Error
	return orderSessions, err
}

func (o *orderSessionRepository) FindByStatus(status string) ([]*domain.OrderSession, error) {
	var orderSessions []*domain.OrderSession
	err := o.db.Where("status = ?", status).Find(&orderSessions).Error
	return orderSessions, err
}

func (o *orderSessionRepository) Update(orderSession *domain.OrderSession) error {
	return o.db.Save(orderSession).Error
}

func (o *orderSessionRepository) Delete(id uint) error {
	return o.db.Delete(&domain.OrderSession{}, id).Error
}
