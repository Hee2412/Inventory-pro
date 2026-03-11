package repository

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"gorm.io/gorm"
)

type OrderSessionRepository interface {
	Create(orderSession *domain.OrderSession) error
	FindById(id uint) (*domain.OrderSession, error)
	FindAll() ([]*domain.OrderSession, error)
	FindByStatus(status string) ([]*domain.OrderSession, error)
	Update(orderSession *domain.OrderSession) error
	Delete(id uint) error
	FindAllPaginated(params request.SessionSearchParams) ([]*domain.OrderSession, int64, error)
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
	if err != nil {
		return nil, err
	}
	return &orderSession, nil
}

func (o *orderSessionRepository) FindAll() ([]*domain.OrderSession, error) {
	var orderSessions []*domain.OrderSession
	err := o.db.Find(&orderSessions).Error
	if err != nil {
		return nil, err
	}
	return orderSessions, nil
}

func (o *orderSessionRepository) FindByStatus(status string) ([]*domain.OrderSession, error) {
	var orderSessions []*domain.OrderSession
	err := o.db.Where("status = ?", status).Find(&orderSessions).Error
	if err != nil {
		return nil, err
	}
	return orderSessions, nil
}

func (o *orderSessionRepository) Update(orderSession *domain.OrderSession) error {
	return o.db.Save(orderSession).Error
}

func (o *orderSessionRepository) Delete(id uint) error {
	return o.db.Delete(&domain.OrderSession{}, id).Error
}

func (o *orderSessionRepository) FindAllPaginated(params request.SessionSearchParams) ([]*domain.OrderSession, int64, error) {
	var orderSessions []*domain.OrderSession
	var total int64
	query := o.db.Model(&domain.OrderSession{})

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.FromDate != nil {
		query = query.Where("created_at >= ?", params.FromDate)
	}
	if params.ToDate != nil {
		query = query.Where("created_at <= ?", params.ToDate)
	}
	query.Count(&total)
	offset := (params.Page - 1) * params.Limit
	err := query.Offset(offset).Limit(params.Limit).
		Order("created_at DESC").
		Find(&orderSessions).Error
	return orderSessions, total, err
}
