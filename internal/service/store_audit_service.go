package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sort"
	"time"
)

type StoreAuditService interface {
	GetAuditReport(sessionID uint, storeID uint) (*response.AuditReportItemDetailResponse, error)
	UpdateAuditItem(sessionID uint, storeID uint, req request.UpdateAuditItemsRequest) error
	GetMyAuditReports(storeID uint) ([]*response.AuditReportItemDetailResponse, error)
}

type storeAuditService struct {
	auditSessionRepo repository.AuditSessionRepository
	storeAuditRepo   repository.StoreAuditReportRepository
	userRepo         repository.UserRepository
}

func NewStoreAuditService(
	auditSessionRepo repository.AuditSessionRepository,
	storeAuditRepo repository.StoreAuditReportRepository,
	userRepo repository.UserRepository,
) StoreAuditService {
	return &storeAuditService{
		auditSessionRepo: auditSessionRepo,
		storeAuditRepo:   storeAuditRepo,
		userRepo:         userRepo,
	}
}

// GetAuditReport GetOrCreateAuditReport Store get in session => found all products
func (s *storeAuditService) GetAuditReport(sessionID uint, storeID uint) (*response.AuditReportItemDetailResponse, error) {
	session, err := s.auditSessionRepo.FindById(sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%w: failed to fetch session data: %v", domain.ErrDatabase, err)
	}
	items, err := s.storeAuditRepo.FindByAuditSessionAndStore(storeID, sessionID)
	result := &response.AuditReportItemDetailResponse{
		SessionTitle: session.Title,
		TotalItems:   0,
		Items:        []*response.AuditItemsResponse{},
	}
	if err == nil && len(items) > 0 {
		result.TotalItems = len(items)
		result.Items = toAuditReportItemResponse(items)
	}
	return result, nil
}

func (s *storeAuditService) UpdateAuditItem(sessionID uint, storeID uint, req request.UpdateAuditItemsRequest) error {
	session, err := s.auditSessionRepo.FindById(sessionID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrSessionNotFound
		}
		return fmt.Errorf("%w: failed to fetch session data: %v", domain.ErrDatabase, err)
	}
	if session.Status == "CLOSED" {
		return errors.New("session is already closed")
	}
	if session.Status == "OPEN" && time.Now().After(session.EndDate) {
		session.Status = "CLOSED"
		err := s.auditSessionRepo.Update(session)
		if err != nil {
			return fmt.Errorf("%w: update status failed", domain.ErrDatabase)
		}
		return fmt.Errorf("%w: session is closed", domain.ErrSessionClosed)
	}
	for _, itemReq := range req.Items {
		item, err := s.storeAuditRepo.FindBySessionStoreAndProduct(sessionID, storeID, itemReq.ProductID)
		if err != nil {
			continue
		}
		if item == nil {
			continue
		}
		item.ActualStock = itemReq.ActualStock
		item.Variance = itemReq.ActualStock - item.SystemStock
		err = s.storeAuditRepo.Update(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *storeAuditService) GetMyAuditReports(storeID uint) ([]*response.AuditReportItemDetailResponse, error) {
	items, err := s.storeAuditRepo.FindByStoreId(storeID)
	if err != nil {
		return make([]*response.AuditReportItemDetailResponse, 0), err
	}
	sessionMap := make(map[uint][]*domain.StoreAuditReport)
	for _, item := range items {
		sessionMap[item.SessionID] = append(sessionMap[item.SessionID], item)
	}
	result := make([]*response.AuditReportItemDetailResponse, 0)
	for sessionID, sessionItems := range sessionMap {
		session, err := s.auditSessionRepo.FindById(sessionID)
		if err != nil {
			continue
		}
		result = append(result, &response.AuditReportItemDetailResponse{
			SessionTitle: session.Title,
			TotalItems:   len(sessionItems),
			Items:        toAuditReportItemResponse(sessionItems),
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].SessionTitle > result[j].SessionTitle
	})
	return result, nil
}
