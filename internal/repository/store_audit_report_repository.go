package repository

import (
	"Inventory-pro/internal/domain"
	"gorm.io/gorm"
)

type StoreAuditReportRepository interface {
	Create(storeAuditReport ...*domain.StoreAuditReport) error
	FindById(id uint) (*domain.StoreAuditReport, error)
	FindByAuditSessionAndStore(storeId uint, auditSession uint) (*domain.StoreAuditReport, error)
	FindByAuditSessionID(sessionId uint) ([]*domain.StoreAuditReport, error)
	FindBySessionAndProduct(sessionId uint, productId uint) ([]*domain.StoreAuditReport, error)
	FindByStoreId(storeId uint) ([]*domain.StoreAuditReport, error)
	Update(storeAuditReport *domain.StoreAuditReport) error
	FindBySessionStoreAndProduct(sessionID uint, storeID uint, productID uint) (*domain.StoreAuditReport, error)
	Delete(id uint) error
}

type storeAuditRepository struct {
	db *gorm.DB
}

func NewStoreAuditRepository(db *gorm.DB) StoreAuditReportRepository {
	return &storeAuditRepository{db: db}
}

// Create Variadic
func (s *storeAuditRepository) Create(reports ...*domain.StoreAuditReport) error {
	if len(reports) == 0 {
		return nil
	}
	if len(reports) == 1 {
		return s.db.Create(reports[0]).Error
	}
	return s.db.CreateInBatches(reports, 100).Error
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

// FindByAuditSessionAndStore Response 1 record (report) by store x in session y
func (s *storeAuditRepository) FindByAuditSessionAndStore(storeId uint, sessionId uint) (*domain.StoreAuditReport, error) {
	var report domain.StoreAuditReport
	err := s.db.Where("store_id = ? AND session_id = ?", storeId, sessionId).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

// FindByAuditSessionID Find [] all store in 1 session
func (s *storeAuditRepository) FindByAuditSessionID(sessionId uint) ([]*domain.StoreAuditReport, error) {
	var auditSession []*domain.StoreAuditReport
	err := s.db.Where("session_id = ?", sessionId).Find(&auditSession).Error
	if err != nil {
		return nil, err
	}
	return auditSession, nil
}

func (s *storeAuditRepository) FindBySessionAndProduct(sessionID uint, productID uint) ([]*domain.StoreAuditReport, error) {
	var auditSession []*domain.StoreAuditReport
	err := s.db.Where("session_id = ? AND product_id = ?", sessionID, productID).Find(&auditSession).Error
	if err != nil {
		return nil, err
	}
	return auditSession, nil
}

// FindByStoreId Audit history of 1 store
func (s *storeAuditRepository) FindByStoreId(storeId uint) ([]*domain.StoreAuditReport, error) {
	var auditSession []*domain.StoreAuditReport
	err := s.db.Where("store_id = ?", storeId).Find(&auditSession).Error
	if err != nil {
		return nil, err
	}
	return auditSession, nil
}

func (s *storeAuditRepository) FindBySessionStoreAndProduct(sessionId uint, storeID uint, productID uint) (*domain.StoreAuditReport, error) {
	var report domain.StoreAuditReport
	err := s.db.Where("session_id = ? AND store_id= ? AND product_id = ?", sessionId, storeID, productID).First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (s *storeAuditRepository) Update(storeAuditReport *domain.StoreAuditReport) error {
	return s.db.Save(storeAuditReport).Error
}

func (s *storeAuditRepository) Delete(id uint) error {
	return s.db.Delete(&domain.StoreAuditReport{}, id).Error
}
