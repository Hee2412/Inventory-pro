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
	UpdateAuditItem(sessionID uint, storeID uint, req request.UpdateAuditItem) error
	SubmitAuditReport(sessionID uint, storeID uint) error
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
	if session.Status != "OPEN" {
		return nil, errors.New("session is not already open")
	}
	//find items in session
	items, err := s.storeAuditRepo.FindByAuditSessionAndStore(storeID, sessionID)
	if err != nil || len(items) == 0 {
		return &response.AuditReportItemDetailResponse{
			SessionTitle: session.Title,
			StoreName:    "",
			TotalItems:   0,
			Items:        []response.AuditItemsResponse{},
		}, nil
	}
	result := &response.AuditReportItemDetailResponse{
		SessionTitle: session.Title,
		StoreName:    items[0].Store.StoreName,
		TotalItems:   len(items),
		Items:        toAuditReportItemResponse(items),
	}
	return result, nil
}

func (s *storeAuditService) UpdateAuditItem(sessionID uint, storeID uint, req request.UpdateAuditItem) error {
	//find 1 item
	item, err := s.storeAuditRepo.FindBySessionStoreAndProduct(sessionID, storeID, req.ProductID)
	if err != nil {
		return errors.New("item not found")
	}
	if item.Status != "DRAFT" {
		return errors.New("can only update draft item")
	}
	item.ActualStock = req.ActualStock
	item.Variance = req.ActualStock - item.SystemStock
	return s.storeAuditRepo.Update(item)
}

func (s *storeAuditService) SubmitAuditReport(sessionID uint, storeID uint) error {
	//get all items
	items, err := s.storeAuditRepo.FindByAuditSessionAndStore(storeID, sessionID)
	if err != nil || len(items) == 0 {
		return errors.New("no item to submit")
	}

	//check null
	for _, item := range items {
		if item.Status != "DRAFT" {
			return errors.New("can only submit draft item")
		}
	}
	now := time.Now()
	for _, item := range items {
		item.Status = "SUBMITTED"
		item.SubmittedAt = &now

		err := s.storeAuditRepo.Update(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *storeAuditService) GetMyAuditReports(storeID uint) ([]response.AuditReportItemDetailResponse, error) {
	//find by storeID
	items, err := s.storeAuditRepo.FindByStoreId(storeID)
	if err != nil {
		return make([]response.AuditReportItemDetailResponse, 0), err
	}
	sessionMap := make(map[uint][]*domain.StoreAuditReport)
	for _, item := range items {
		sessionMap[item.SessionID] = append(sessionMap[item.SessionID], item)
	}
	result := make([]response.AuditReportItemDetailResponse, 0)
	for sessionID, sessionItems := range sessionMap {
		session, _ := s.auditSessionRepo.FindById(sessionID)
		result = append(result, response.AuditReportItemDetailResponse{
			SessionTitle: session.Title,
			StoreName:    sessionItems[0].Store.StoreName,
			TotalItems:   len(sessionItems),
			Items:        toAuditReportItemResponse(sessionItems),
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].SessionTitle > result[j].SessionTitle
	})
	return result, nil
}
