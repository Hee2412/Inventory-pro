package repository

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"gorm.io/gorm"
)

type AuditSessionRepository interface {
	Create(auditSession *domain.AuditSession) error
	FindById(id uint) (*domain.AuditSession, error)
	FindAll() ([]*domain.AuditSession, error)
	FindByStatus(status string) ([]*domain.AuditSession, error)
	Update(auditSession *domain.AuditSession) error
	DeleteById(id uint) error
	HardDelete(id uint) error

	AdminFindById(id uint) (*domain.AuditSession, error)
	AdminFindAll() ([]*domain.AuditSession, error)
	FindAllPaginated(params request.AuditSessionSearchParams) ([]*domain.AuditSession, int64, error)
}

type auditSessionRepository struct {
	db *gorm.DB
}

func NewAuditSessionRepository(db *gorm.DB) AuditSessionRepository {
	return &auditSessionRepository{db: db}
}

func (a *auditSessionRepository) Create(auditSession *domain.AuditSession) error {
	return a.db.Create(auditSession).Error
}

func (a *auditSessionRepository) FindById(id uint) (*domain.AuditSession, error) {
	var auditSession domain.AuditSession
	err := a.db.Where("id = ?", id).First(&auditSession).Error
	if err != nil {
		return nil, err
	}
	return &auditSession, nil
}

func (a *auditSessionRepository) FindAll() ([]*domain.AuditSession, error) {
	var auditSessions []*domain.AuditSession
	err := a.db.Find(&auditSessions).Error
	if err != nil {
		return nil, err
	}
	return auditSessions, nil
}

func (a *auditSessionRepository) FindByStatus(status string) ([]*domain.AuditSession, error) {
	var auditSessions []*domain.AuditSession
	err := a.db.Where("status = ?", status).Find(&auditSessions).Error
	if err != nil {
		return nil, err
	}
	return auditSessions, nil
}

func (a *auditSessionRepository) Update(auditSession *domain.AuditSession) error {
	return a.db.Save(auditSession).Error
}

func (a *auditSessionRepository) DeleteById(id uint) error {
	return a.db.Delete(&domain.AuditSession{}, id).Error
}

func (a *auditSessionRepository) HardDelete(id uint) error {
	return a.db.Unscoped().Delete(&domain.AuditSession{}, id).Error
}

func (a *auditSessionRepository) AdminFindById(id uint) (*domain.AuditSession, error) {
	var auditSession domain.AuditSession
	err := a.db.Unscoped().Where("id = ?", id).First(&auditSession).Error
	if err != nil {
		return nil, err
	}
	return &auditSession, nil
}

func (a *auditSessionRepository) AdminFindAll() ([]*domain.AuditSession, error) {
	var auditSessions []*domain.AuditSession
	err := a.db.Unscoped().Find(&auditSessions).Error
	if err != nil {
		return nil, err
	}
	return auditSessions, nil
}

func (a *auditSessionRepository) FindAllPaginated(params request.AuditSessionSearchParams) ([]*domain.AuditSession, int64, error) {
	var sessions []*domain.AuditSession
	var total int64
	query := a.db.Model(&domain.AuditSession{})

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.ToDate != nil {
		query = query.Where("created_at >= ?", params.ToDate)
	}
	if params.FromDate != nil {
		query = query.Where("created_at <= ?", params.FromDate)
	}
	query.Count(&total)
	offset := (params.Page - 1) * params.Limit
	err := query.Offset(offset).Limit(params.Limit).
		Order("created_at DESC").
		Find(&sessions).Error
	return sessions, total, err
}
