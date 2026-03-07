package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
	"time"
)

type SuperAdminAuditService interface {
	GetAllReportsInSession(sessionID uint) ([]response.StoreAuditSummaryResponse, error)
	GetReportDetail(sessionID uint, storeID uint) (*response.AuditReportItemDetailResponse, error)
	GetAuditSummary(sessionID uint) (*response.AuditSummaryResponse, error)
	ApproveStoreReport(storeID uint, sessionID uint, adminID uint) error
	DeclineStoreReport(storeID uint, sessionID uint, reason string, adminID uint) error
}

type superAdminAuditService struct {
	auditSessionRepo repository.AuditSessionRepository
	storeAuditRepo   repository.StoreAuditReportRepository
	userRepo         repository.UserRepository
}

func NewSuperAdminAuditService(
	auditSessionRepo repository.AuditSessionRepository,
	storeAuditRepo repository.StoreAuditReportRepository,
	userRepo repository.UserRepository,
) SuperAdminAuditService {
	return &superAdminAuditService{
		auditSessionRepo: auditSessionRepo,
		storeAuditRepo:   storeAuditRepo,
		userRepo:         userRepo,
	}
}

func getStoreStatus(items []*domain.StoreAuditReport) string {
	if len(items) == 0 {
		return "DRAFT"
	}

	allSubmitted := true
	allApproved := true

	for _, item := range items {
		if item.Status != "SUBMITTED" {
			allSubmitted = false
		}
		if item.Status != "APPROVED" {
			allApproved = false
		}
	}

	if allApproved {
		return "APPROVED"
	}
	if allSubmitted {
		return "SUBMITTED"
	}
	return "DRAFT"
}

func (s *superAdminAuditService) GetAllReportsInSession(sessionID uint) ([]response.StoreAuditSummaryResponse, error) {
	// check session
	session, err := s.auditSessionRepo.FindById(sessionID)
	if err != nil {
		return nil, errors.New("session not found")
	}
	//get all reports
	reports, err := s.storeAuditRepo.FindByAuditSessionID(session.ID)
	if err != nil {
		return make([]response.StoreAuditSummaryResponse, 0), nil
	}
	//group by storeID
	storeMap := make(map[uint][]*domain.StoreAuditReport)
	for _, report := range reports {
		storeMap[report.StoreID] = append(storeMap[report.StoreID], report)
	}
	var result []response.StoreAuditSummaryResponse
	for storeID, items := range storeMap {
		storeName := "Unknown"
		if len(items) > 0 && items[0].Store != nil {
			storeName = items[0].Store.StoreName
		}
		result = append(result, response.StoreAuditSummaryResponse{
			StoreId:   storeID,
			StoreName: storeName,
			Status:    getStoreStatus(items),
		})
	}
	return result, nil
}

func (s *superAdminAuditService) GetReportDetail(sessionID uint, storeID uint) (*response.AuditReportItemDetailResponse, error) {
	//get session info
	session, err := s.auditSessionRepo.FindById(sessionID)
	if err != nil {
		return nil, errors.New("session not found")
	}
	store, err := s.userRepo.FindById(storeID)
	if err != nil {
		return nil, errors.New("store not found")
	}
	//get items in session
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

func (s *superAdminAuditService) GetAuditSummary(sessionID uint) (*response.AuditSummaryResponse, error) {
	reports, err := s.storeAuditRepo.FindByAuditSessionID(sessionID)
	if err != nil {
		return nil, errors.New("no record found")
	}
	session, err := s.auditSessionRepo.FindById(sessionID)
	if err != nil {
		return nil, errors.New("session not found")
	}

	//create mapping, group items by store
	storeItemsMap := make(map[uint][]*domain.StoreAuditReport)
	productSet := make(map[uint]bool)
	totalVariance := 0.0

	for _, report := range reports {
		storeItemsMap[report.StoreID] = append(storeItemsMap[report.StoreID], report)
		productSet[report.ProductID] = true
		totalVariance += report.Variance
	}

	//calculate summary
	var storeApproved int
	var storeDraft int
	var issues []response.StoreIssue

	for storeID, items := range storeItemsMap {
		status := getStoreStatus(items)

		if status == "APPROVED" {
			storeApproved++
		} else {
			storeDraft++
		}
		//calculate store variance
		storeVariance := 0.0
		for _, item := range items {
			storeVariance += item.Variance
		}
		//add to issue if variance < -10
		if storeVariance < -10 {
			storeName := "Unknown"
			if len(items) > 0 && items[0].Store != nil {
				storeName = items[0].Store.StoreName
			}
			issues = append(issues, response.StoreIssue{
				StoreID:   storeID,
				StoreName: storeName,
				Status:    status,
				Variance:  storeVariance,
			})
		}
	}

	return &response.AuditSummaryResponse{
		SessionTitle:     session.Title,
		TotalStores:      len(storeItemsMap),
		StoresApproved:   storeApproved,
		StoreDraft:       storeDraft,
		TotalProducts:    len(productSet),
		TotalVariance:    totalVariance,
		StoresWithIssues: issues,
	}, nil
}

func (s *superAdminAuditService) ApproveStoreReport(storeID uint, sessionID uint, adminID uint) error {
	//get all items
	items, err := s.storeAuditRepo.FindByAuditSessionAndStore(storeID, sessionID)
	if err != nil {
		return errors.New("cant get items")
	}
	if len(items) == 0 {
		return errors.New("no record found")
	}
	status := getStoreStatus(items)
	if status != "DRAFT" {
		return errors.New("can only approve submitted reports")
	}

	updateData := map[string]interface{}{
		"status":      "APPROVED",
		"approved_at": time.Now(),
		"approved_by": adminID,
	}
	return s.storeAuditRepo.UpdateStatusByStore(storeID, sessionID, updateData)
}

func (s *superAdminAuditService) DeclineStoreReport(storeID uint, sessionID uint, reason string, adminID uint) error {
	//get all items
	items, err := s.storeAuditRepo.FindByAuditSessionAndStore(storeID, sessionID)
	if err != nil {
		return errors.New("cant get items")
	}
	if len(items) == 0 {
		return errors.New("no record found")
	}
	status := getStoreStatus(items)
	if status != "DRAFT" {
		return errors.New("can only decline store reports")
	}

	updateData := map[string]interface{}{
		"status":      "DECLINED",
		"approved_at": time.Now(),
		"approved_by": adminID,
		"reason":      reason,
	}
	return s.storeAuditRepo.UpdateStatusByStore(storeID, sessionID, updateData)
}
