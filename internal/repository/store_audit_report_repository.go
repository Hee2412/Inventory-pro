package repository

import (
	"Inventory-pro/internal/domain"
	"gorm.io/gorm"
)

type StoreAuditReportRepository interface {
	Create(storeAuditReport *domain.StoreAuditReport) error
	FindById(id uint) (*domain.StoreAuditReport, error)
	FindByAuditSessionAndStore(storeId uint, auditSession uint) ([]*domain.StoreAuditReport, error)
	FindByAuditSessionID(sessionId uint) ([]*domain.StoreAuditReport, error)
	FindByStoreId(storeId uint) ([]*domain.StoreAuditReport, error)
	Update(storeAuditReport *domain.StoreAuditReport) error
}

type storeAuditRepository struct {
	db *gorm.DB
}

func NewStoreAuditRepository(db *gorm.DB) StoreAuditReportRepository {
	return &storeAuditRepository{db: db}
}

func (s *storeAuditRepository) Create(storeAuditReport *domain.StoreAuditReport) error {
	return s.db.Create(storeAuditReport).Error
}

// FindById Searching for exact 1 line, 1 record
func (s *storeAuditRepository) FindById(id uint) (*domain.StoreAuditReport, error) {
	var storeAuditReport domain.StoreAuditReport
	err := s.db.Where("id = ?", id).First(&storeAuditReport).Error
	if err != nil {
		return nil, err
	}
	return &storeAuditReport, nil
}

// FindByAuditSessionAndStore Response [] all products by store x in session y
func (s *storeAuditRepository) FindByAuditSessionAndStore(storeId uint, sessionId uint) ([]*domain.StoreAuditReport, error) {
	var storeAuditReport []*domain.StoreAuditReport
	err := s.db.Where("store_id = ? AND session_id = ?", storeId, sessionId).First(&storeAuditReport).Error
	if err != nil {
		return nil, err
	}
	return storeAuditReport, nil
}

// FindByAuditSessionID Find [] all store in 1 session
func (s *storeAuditRepository) FindByAuditSessionID(sessionId uint) ([]*domain.StoreAuditReport, error) {
	var auditSession []*domain.StoreAuditReport
	err := s.db.Where("session_id = ?", sessionId).First(&auditSession).Error
	if err != nil {
		return nil, err
	}
	return auditSession, nil
}

// FindByStoreId Audit history of 1 store
func (s *storeAuditRepository) FindByStoreId(storeId uint) ([]*domain.StoreAuditReport, error) {
	var auditSession []*domain.StoreAuditReport
	err := s.db.Where("store_id = ?", storeId).First(&auditSession).Error
	if err != nil {
		return nil, err
	}
	return auditSession, nil
}

func (s *storeAuditRepository) Update(storeAuditReport *domain.StoreAuditReport) error {
	return s.db.Save(storeAuditReport).Error
}
