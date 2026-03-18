package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type SuperAdminAuditService interface {
	GetAllReportsInSession(sessionID uint) ([]response.StoreAuditSummaryResponse, error)
	GetReportDetail(sessionID uint, storeID uint) (*response.AuditReportItemDetailResponse, error)
	GetAuditSummary(sessionID uint) (*response.AuditSummaryResponse, error)
	ApproveStoreReport(storeID uint, sessionID uint, adminID uint) error
	DeclineStoreReport(storeID uint, sessionID uint, reason string, adminID uint) error
	GetIncompleteAudit(sessionID uint) (*response.AuditTrackingResponse, error)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%w: failed to fetch session: %v", domain.ErrDatabase, err)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%w: failed to fetch session: %v", domain.ErrDatabase, err)
	}
	//get items in session
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

func (s *superAdminAuditService) GetAuditSummary(sessionID uint) (*response.AuditSummaryResponse, error) {
	reports, err := s.storeAuditRepo.FindByAuditSessionID(sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w: failed to fetch reports: %v", domain.ErrDatabase, err)
	}
	session, err := s.auditSessionRepo.FindById(sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%w: failed to fetch session: %v", domain.ErrDatabase, err)
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
		return fmt.Errorf("%w: failed to fetch store audit for store %d: %v", domain.ErrDatabase, storeID, err)
	}
	if len(items) == 0 {
		return domain.ErrNotFound
	}
	status := getStoreStatus(items)
	if status != "DRAFT" {
		return fmt.Errorf("%w: only reports in DRAFT status can be approved", domain.ErrInvalidInput)
	}

	updateData := map[string]interface{}{
		"status":      "APPROVED",
		"approved_at": time.Now(),
		"approved_by": adminID,
	}
	err = s.storeAuditRepo.UpdateStatusByStore(storeID, sessionID, updateData)
	if err != nil {
		return fmt.Errorf("%w: failed to update report status: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (s *superAdminAuditService) DeclineStoreReport(storeID uint, sessionID uint, reason string, adminID uint) error {
	//get all items
	items, err := s.storeAuditRepo.FindByAuditSessionAndStore(storeID, sessionID)
	if err != nil {
		return fmt.Errorf("%w: failed to fetch store audit for store %d: %v", domain.ErrDatabase, storeID, err)
	}
	if len(items) == 0 {
		return domain.ErrNotFound
	}
	status := getStoreStatus(items)
	if status != "DRAFT" {
		return fmt.Errorf("%w: only reports in DRAFT status can be approved", domain.ErrInvalidInput)
	}

	updateData := map[string]interface{}{
		"status":      "DECLINED",
		"approved_at": time.Now(),
		"approved_by": adminID,
		"reason":      reason,
	}
	err = s.storeAuditRepo.UpdateStatusByStore(storeID, sessionID, updateData)
	if err != nil {
		return fmt.Errorf("%w: failed to update report status: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (s *superAdminAuditService) GetIncompleteAudit(sessionID uint) (*response.AuditTrackingResponse, error) {
	stores, err := s.userRepo.FindByRoleAndActive("store", true)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to fetch data %v", domain.ErrDatabase, err)
	}
	session, err := s.auditSessionRepo.FindById(sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%w: failed to fetch session: %v", domain.ErrDatabase, err)
	}
	allReports, err := s.storeAuditRepo.FindByAuditSessionID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to fetch data%v", domain.ErrInternalServer, err)
	}
	reportMap := make(map[uint][]*domain.StoreAuditReport)
	for _, r := range allReports {
		reportMap[r.StoreID] = append(reportMap[r.StoreID], r)
	}

	var incompleteStores []*response.StoreTracking
	completedCount := 0

	for _, store := range stores {
		reports := reportMap[store.ID]
		if len(reports) == 0 {
			continue
		}

		hasUpdated := false
		for _, report := range reports {
			if report.UpdatedAt.After(session.CreatedAt) {
				hasUpdated = true
				break
			}
		}
		if hasUpdated {
			completedCount++
		} else {
			incompleteStores = append(incompleteStores, &response.StoreTracking{
				StoreID:   store.ID,
				StoreName: store.StoreName,
			})
		}
	}
	return &response.AuditTrackingResponse{
		Completed:       completedCount,
		Incomplete:      len(incompleteStores),
		IncompleteStore: incompleteStores,
	}, nil
}
