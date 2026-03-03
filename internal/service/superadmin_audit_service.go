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
	ApproveStoreReport(sessionID uint, storeID uint, adminID uint) error
	DeclineStoreReport(sessionID uint, storeID uint, reason string, adminID uint) error
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
		return nil, err
	}
	//group by storeID
	storeMap := make(map[uint]*response.StoreAuditSummaryResponse)
	for _, report := range reports {
		if _, exist := storeMap[report.StoreID]; !exist {
			storeName := "unknown"
			if report.Store != nil {
				storeName = report.Store.StoreName
			}
			storeMap[report.StoreID] = &response.StoreAuditSummaryResponse{
				StoreId:   report.StoreID,
				StoreName: storeName,
				Status:    "SUBMITTED",
			}
		}
		if report.Status == "DRAFT" {
			storeMap[report.StoreID].Status = "DRAFT"
		}
	}
	var result []response.StoreAuditSummaryResponse
	for _, store := range storeMap {
		result = append(result, *store)
	}
	return result, nil
}

func (s *superAdminAuditService) GetReportDetail(sessionID uint, storeID uint) (*response.AuditReportItemDetailResponse, error) {
	//get session info
	session, err := s.auditSessionRepo.FindById(sessionID)
	if err != nil {
		return nil, errors.New("session not found")
	}
	//get items in session
	items, err := s.storeAuditRepo.FindByAuditSessionAndStore(sessionID, storeID)
	if err != nil {
		return nil, err
	}
	// convert in to response
	result := &response.AuditReportItemDetailResponse{
		SessionTitle: session.Title,
		StoreName:    items[0].Store.StoreName,
		TotalItems:   len(items),
		Items:        toAuditReportItemResponse(items),
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
	//create mapping

	storeMap := make(map[uint]bool)
	storeName := make(map[uint]string)
	productSet := make(map[uint]bool)

	for _, report := range reports {
		productSet[report.StoreID] = true
		if _, exists := storeMap[report.StoreID]; !exists {
			storeName[report.StoreID] = report.Store.StoreName
			storeMap[report.StoreID] = false
		}
		if report.Status == "DRAFT" {
			storeMap[report.StoreID] = true
		}
	}
	var storeSubmitted int
	var issues []response.StoreIssue

	for id, isDraft := range storeMap {
		if isDraft {
			issues = append(issues, response.StoreIssue{
				StoreID:   id,
				StoreName: storeName[id],
			})
		} else {
			storeSubmitted++
		}
	}
	return &response.AuditSummaryResponse{
		SessionTitle:     session.Title,
		TotalStores:      len(storeMap),
		StoresSubmitted:  storeSubmitted,
		TotalProducts:    len(productSet),
		StoresWithIssues: issues,
	}, nil
}

func (s *superAdminAuditService) ApproveStoreReport(adminID uint, sessionID uint, storeID uint) error {
	//get all items
	items, err := s.storeAuditRepo.FindByAuditSessionAndStore(sessionID, storeID)
	if err != nil {
		return errors.New("cant get items")
	}
	if len(items) == 0 {
		return errors.New("no record found")
	}
	status := getStoreStatus(items)
	if status == "DRAFT" {
		return errors.New("store is DRAFT")
	}
	if status == "APPROVED" {
		return errors.New("store is APPROVED")
	}
	updateData := map[string]interface{}{
		"status":      "APPROVED",
		"approved_at": time.Now(),
		"approved_by": adminID,
	}
	return s.storeAuditRepo.UpdateStatusByStore(sessionID, storeID, updateData)
}

func (s *superAdminAuditService) DeclineStoreReport(sessionID uint, storeID uint, reason string, adminID uint) error {
	//get all items
	items, err := s.storeAuditRepo.FindByAuditSessionAndStore(sessionID, storeID)
	if err != nil {
		return errors.New("cant get items")
	}
	if len(items) == 0 {
		return errors.New("no record found")
	}
	status := getStoreStatus(items)
	if status == "DRAFT" {
		return errors.New("store is DRAFT")
	}
	if status == "APPROVED" {
		return errors.New("store is APPROVED")
	}
	updateData := map[string]interface{}{
		"status":      "DECLINED",
		"approved_at": time.Now(),
		"approved_by": adminID,
		"reason":      reason,
	}
	return s.storeAuditRepo.UpdateStatusByStore(sessionID, storeID, updateData)
}
