package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
	"sort"
	"time"
)

type StoreAuditService interface {
	GetAuditReport(sessionID uint, storeID uint) (*response.AuditReportItemDetailResponse, error)
	UpdateAuditItem(sessionID uint, storeID uint, req request.UpdateAuditItemsRequest) error
	GetMyAuditReports(storeID uint) ([]response.AuditReportItemDetailResponse, error)
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

func toAuditReportItemResponse(items []*domain.StoreAuditReport) []response.AuditItemsResponse {
	result := make([]response.AuditItemsResponse, 0)
	for _, item := range items {
		result = append(result, response.AuditItemsResponse{
			ProductID:   item.ProductID,
			ProductName: item.Product.ProductName,
			SystemStock: item.SystemStock,
			ActualStock: item.ActualStock,
			Variance:    item.Variance,
		})
	}
	return result
}

// GetAuditReport GetOrCreateAuditReport Store get in session => found all products
func (s *storeAuditService) GetAuditReport(sessionID uint, storeID uint) (*response.AuditReportItemDetailResponse, error) {
	//check status session
	session, err := s.auditSessionRepo.FindById(sessionID)
	if err != nil {
		return nil, errors.New("session not found")
	}
	store, err := s.userRepo.FindById(storeID)
	if err != nil {
		return nil, errors.New("store not found")
	}
	//find items in session
	items, err := s.storeAuditRepo.FindByAuditSessionAndStore(storeID, sessionID)
	result := &response.AuditReportItemDetailResponse{
		SessionTitle: session.Title,
		StoreName:    store.StoreName,
		TotalItems:   0,
		Items:        []response.AuditItemsResponse{},
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
		return errors.New("session not found")
	}
	if session.Status == "OPEN" && time.Now().After(session.EndDate) {
		session.Status = "CLOSED"
		err := s.auditSessionRepo.Update(session)
		if err != nil {
			return err
		}
		return errors.New("session has expired and been closed")
	}
	if session.Status == "CLOSED" {
		return errors.New("session is already closed")
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

func (s *storeAuditService) GetMyAuditReports(storeID uint) ([]response.AuditReportItemDetailResponse, error) {
	//find reports by storeID
	items, err := s.storeAuditRepo.FindByStoreId(storeID)
	if err != nil {
		return make([]response.AuditReportItemDetailResponse, 0), err
	}
	store, err := s.userRepo.FindById(storeID)
	if err != nil {
		return nil, errors.New("store not found")
	}
	sessionMap := make(map[uint][]*domain.StoreAuditReport)
	for _, item := range items {
		sessionMap[item.SessionID] = append(sessionMap[item.SessionID], item)
	}
	result := make([]response.AuditReportItemDetailResponse, 0)
	for sessionID, sessionItems := range sessionMap {
		session, err := s.auditSessionRepo.FindById(sessionID)
		if err != nil {
			continue
		}
		result = append(result, response.AuditReportItemDetailResponse{
			SessionTitle: session.Title,
			StoreName:    store.StoreName,
			TotalItems:   len(sessionItems),
			Items:        toAuditReportItemResponse(sessionItems),
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].SessionTitle > result[j].SessionTitle
	})
	return result, nil
}
