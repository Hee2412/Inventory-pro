package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
	"fmt"
	"time"
)

type AuditSessionService interface {
	CreateAuditSession(req request.CreateAuditSessionRequest) (*response.AuditSessionResponse, error)
	GetAllAuditSessions() ([]response.AuditSessionResponse, error)
	GetAuditSessionByID(sessionID uint) (*response.AuditSessionDetailsResponse, error)
	AddProductToAudit(req request.AddProductToAuditRequest) (*response.AddProductResponse, error)
	RemoveProductFromAudit(auditSessionID uint, productID uint) error
	CloseAuditSession(auditSessionID uint) error
}

type auditSessionService struct {
	auditSessionRepo repository.AuditSessionRepository
	storeAuditRepo   repository.StoreAuditReportRepository
	userRepo         repository.UserRepository
	productRepo      repository.ProductRepository
}

func NewAuditSessionService(
	auditSessionRepo repository.AuditSessionRepository,
	storeAuditRepo repository.StoreAuditReportRepository,
	userRepo repository.UserRepository,
	productRepo repository.ProductRepository) AuditSessionService {
	return &auditSessionService{
		auditSessionRepo: auditSessionRepo,
		storeAuditRepo:   storeAuditRepo,
		userRepo:         userRepo,
		productRepo:      productRepo,
	}
}
func toAuditSessionResponse(session *domain.AuditSession) response.AuditSessionResponse {
	return response.AuditSessionResponse{
		SessionID: session.ID,
		Title:     session.Title,
		AuditType: session.AuditType,
		StartDate: session.StartDate,
		EndDate:   session.EndDate,
		Status:    session.Status,
	}
}

func (a *auditSessionService) CreateAuditSession(req request.CreateAuditSessionRequest) (*response.AuditSessionResponse, error) {
	if req.EndDate.Before(req.StartDate) {
		return nil, errors.New("the end date cannot be before the start date")
	}
	if req.StartDate.Before(time.Now()) {
		return nil, errors.New("the start date is in the past")
	}
	newAuditSession := &domain.AuditSession{
		Title:     req.Title,
		AuditType: req.AuditType,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Status:    "OPEN",
	}
	err := a.auditSessionRepo.Create(newAuditSession)
	if err != nil {
		return nil, err
	}
	result := toAuditSessionResponse(newAuditSession)
	return &result, nil
}

func (a *auditSessionService) GetAllAuditSessions() ([]response.AuditSessionResponse, error) {
	sessions, err := a.auditSessionRepo.FindAll()
	if err != nil {
		return nil, err
	}
	result := make([]response.AuditSessionResponse, 0)
	for _, session := range sessions {
		result = append(result, toAuditSessionResponse(session))
	}
	return result, nil
}

func (a *auditSessionService) GetAuditSessionByID(sessionID uint) (*response.AuditSessionDetailsResponse, error) {
	//var sessionId
	session, err := a.auditSessionRepo.FindById(sessionID)
	if err != nil {
		return nil, err
	}
	//check store audit in session
	reports, err := a.storeAuditRepo.FindByAuditSessionID(sessionID)
	if err != nil {
		return &response.AuditSessionDetailsResponse{
			SessionInfo: toAuditSessionResponse(session),
			Report:      []response.StoreAuditSummaryResponse{},
		}, nil
	}
	//mapping
	storeMap := make(map[uint]response.StoreAuditSummaryResponse)
	for _, r := range reports {
		if _, exist := storeMap[r.StoreID]; !exist {
			storeName := "unknown"
			if r.Store != nil {
				storeName = r.Store.StoreName
			}
			storeMap[r.StoreID] = response.StoreAuditSummaryResponse{
				StoreId:   r.StoreID,
				StoreName: storeName,
				Status:    r.Status,
			}
		}
	}
	//convert map to slice
	var storeResponse []response.StoreAuditSummaryResponse
	for _, store := range storeMap {
		storeResponse = append(storeResponse, store)
	}
	return &response.AuditSessionDetailsResponse{
		SessionInfo: toAuditSessionResponse(session),
		Report:      storeResponse,
	}, nil
}

func (a *auditSessionService) AddProductToAudit(req request.AddProductToAuditRequest) (*response.AddProductResponse, error) {
	//check session/status
	session, err := a.auditSessionRepo.FindById(req.AuditSessionID)
	if err != nil {
		return nil, errors.New("session not found")
	}
	if session.Status != "OPEN" {
		return nil, errors.New("session is not OPEN")
	}
	//get all stores FindByRole
	stores, err := a.userRepo.FindByRole("store")
	if err != nil {
		return nil, err
	}

	var added int
	var errs []string
	var reports []*domain.StoreAuditReport

	//for each product -> create report/StoreAuditReport for each store
	for _, productID := range req.ProductID {
		//check product exists
		_, err := a.productRepo.FindById(productID)
		if err != nil {
			errs = append(errs, fmt.Sprintf("product %d not found", productID))
			continue
		}
		//create report item for each store

		for _, store := range stores {
			//check if report exists
			existing, _ := a.storeAuditRepo.FindBySessionStoreAndProduct(
				store.ID,
				session.ID,
				productID)
			if existing != nil {
				//it does -> add item to report
				errs = append(errs,
					fmt.Sprintf("product %d already exists for store %d",
						productID, store.ID))
				continue
			}
			reports = append(reports, &domain.StoreAuditReport{
				SessionID:   session.ID,
				StoreID:     store.ID,
				ProductID:   productID,
				SystemStock: 0,
				ActualStock: 0,
				Variance:    0,
				Status:      "DRAFT",
			})
			added++
		}
	}
	if len(reports) > 0 {
		err := a.storeAuditRepo.Create(reports...)
		if err != nil {
			return nil, err
		}
	}
	return &response.AddProductResponse{
		Added:   added,
		Skipped: len(errs),
		Errs:    errs,
	}, nil
}

func (a *auditSessionService) RemoveProductFromAudit(auditSessionID uint, productID uint) error {
	//check session
	session, err := a.auditSessionRepo.FindById(auditSessionID)
	if err != nil {
		return errors.New("session not found")
	}
	//check status
	if session.Status == "CLOSED" {
		return errors.New("session is CLOSE")
	}
	//check product
	_, err = a.productRepo.FindById(productID)
	if err != nil {
		return errors.New("product not found")
	}
	//find all reports that have productID in session
	reports, err := a.storeAuditRepo.FindBySessionAndProduct(
		auditSessionID, productID)
	if err != nil || len(reports) == 0 {
		return errors.New("product not found in this session")
	}
	//delete report in session
	for _, report := range reports {
		if report.Status != "DRAFT" {
			return errors.New("cannot remove product from submitted reports")
		}
		err = a.storeAuditRepo.Delete(report.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *auditSessionService) CloseAuditSession(auditSessionID uint) error {
	session, err := a.auditSessionRepo.FindById(auditSessionID)
	if err != nil {
		return errors.New("session not found")
	}
	if session.Status == "CLOSED" {
		return errors.New("session is CLOSE")
	}
	if session.Status != "OPEN" {
		return errors.New("can only close the open session")
	}
	session.Status = "CLOSED"
	return a.auditSessionRepo.Update(session)
}
