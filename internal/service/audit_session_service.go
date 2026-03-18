package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type AuditSessionService interface {
	CreateAuditSession(req request.CreateAuditSessionRequest, createdBy uint) (*response.AuditSessionResponse, error)
	GetAllAuditSessions() ([]*response.AuditSessionResponse, error)
	GetAuditSessionByID(sessionID uint) (*response.AuditSessionDetailsResponse, error)
	AddProductToAudit(req request.AddProductToAuditRequest) (*response.AddProductResponse, error)
	RemoveProductFromAudit(sessionID uint, productID uint) error
	CloseAuditSession(auditSessionID uint) error
	UpdateAuditSession(sessionID uint, req request.UpdateAuditSessionRequest) error
	AutoCloseExpiredSession() error
	GetAllSessionsPaginated(params request.AuditSessionSearchParams) ([]*response.AuditSessionResponse, int64, error)
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

func (a *auditSessionService) AutoCloseExpiredSession() error {
	sessions, _ := a.auditSessionRepo.FindByStatus("OPEN")
	for _, session := range sessions {
		if time.Now().After(session.EndDate) {
			session.Status = "CLOSED"
			err := a.auditSessionRepo.Update(session)
			if err != nil {
				return err
			}
			log.Printf("AutoClose expired session %d:", session.ID)
		}
	}
	return nil
}

func (a *auditSessionService) CreateAuditSession(req request.CreateAuditSessionRequest, createdBy uint) (*response.AuditSessionResponse, error) {
	if req.EndDate.Before(req.StartDate) {
		return nil, fmt.Errorf("%w: end date cannot before start date", domain.ErrInvalidInput)
	}
	if req.StartDate.Before(time.Now()) {
		return nil, fmt.Errorf("%w: start date is in the past", domain.ErrInvalidInput)
	}
	newAuditSession := &domain.AuditSession{
		Title:     req.Title,
		AuditType: req.AuditType,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Status:    "OPEN",
		CreatedBy: createdBy,
	}
	err := a.auditSessionRepo.Create(newAuditSession)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create session %v", domain.ErrDatabase, err)
	}
	result := toAuditSessionResponse(newAuditSession)
	return result, nil
}

func (a *auditSessionService) GetAllAuditSessions() ([]*response.AuditSessionResponse, error) {
	sessions, err := a.auditSessionRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to fetch session datas", domain.ErrDatabase)
	}
	result := make([]*response.AuditSessionResponse, 0)
	for _, session := range sessions {
		result = append(result, toAuditSessionResponse(session))
	}
	return result, nil
}

func (a *auditSessionService) GetAuditSessionByID(sessionID uint) (*response.AuditSessionDetailsResponse, error) {
	//var sessionId
	session, err := a.auditSessionRepo.FindById(sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%w: failed to fetch session data: %v", domain.ErrDatabase, err)
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%w: failed to fetch session data: %v", domain.ErrDatabase, err)
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
		product, err := a.productRepo.FindById(productID)
		if err != nil {
			errs = append(errs, fmt.Sprintf("product %d not found", productID))
			continue
		}
		//create report item for each store

		for _, store := range stores {
			//check if report exists
			existing, _ := a.storeAuditRepo.FindBySessionStoreAndProduct(
				session.ID,
				store.ID,
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
				StoreName:   store.StoreName,
				ProductName: product.ProductName,
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
			return nil, fmt.Errorf("%w: add product failed", domain.ErrDatabase)
		}
	}
	return &response.AddProductResponse{
		Added:   added,
		Skipped: len(errs),
		Errs:    errs,
	}, nil
}

func (a *auditSessionService) RemoveProductFromAudit(sessionID uint, productID uint) error {
	//check session
	session, err := a.auditSessionRepo.FindById(sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrSessionNotFound
		}
		return fmt.Errorf("%w: failed to fetch session data: %v", domain.ErrDatabase, err)
	}
	//check status
	if session.Status == "CLOSED" {
		return fmt.Errorf("%w: session is closed", domain.ErrSessionClosed)
	}
	//find all reports that have productID in session
	reports, err := a.storeAuditRepo.FindBySessionAndProduct(
		sessionID, productID)
	if err != nil || len(reports) == 0 {
		return fmt.Errorf("%w: product is not asigned in this session: %v", domain.ErrDatabase, err)
	}
	//delete report in session
	for _, report := range reports {
		if report.Status != "DRAFT" {
			return fmt.Errorf("%w: report %s is not DRAFT", domain.ErrDatabase, report.ProductName)
		}
		err = a.storeAuditRepo.Delete(report.ID)
		if err != nil {
			return fmt.Errorf("%w: failed to remove report %s", domain.ErrDatabase, report.ProductName)
		}
	}
	return nil
}

func (a *auditSessionService) CloseAuditSession(auditSessionID uint) error {
	session, err := a.auditSessionRepo.FindById(auditSessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrSessionNotFound
		}
		return fmt.Errorf("%w: failed to fetch session data: %v", domain.ErrDatabase, err)
	}
	if session.Status == "CLOSED" {
		return fmt.Errorf("%w: session is already closed", domain.ErrSessionClosed)
	}
	session.Status = "CLOSED"
	err = a.auditSessionRepo.Update(session)
	if err != nil {
		return fmt.Errorf("%w: failed to close session data: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (a *auditSessionService) UpdateAuditSession(sessionID uint, req request.UpdateAuditSessionRequest) error {
	//check session
	session, err := a.auditSessionRepo.FindById(sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrSessionNotFound
		}
		return fmt.Errorf("%w: failed to fetch session data: %v", domain.ErrDatabase, err)
	}
	if session.Status == "CLOSED" {
		return fmt.Errorf("%w: session is already closed", domain.ErrSessionClosed)
	}
	//call service
	if req.Title != nil {
		session.Title = *req.Title
	}
	if req.Status != nil {
		session.Status = *req.Status
	}
	if req.StartDate != nil {
		session.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		session.EndDate = *req.EndDate
	}
	if req.AuditType != nil {
		session.AuditType = *req.AuditType
	}
	if req.EndDate != nil && req.StartDate != nil {
		if req.EndDate.Before(*req.StartDate) {
			return fmt.Errorf("%w: end date cannot be before start date", domain.ErrInvalidInput)
		}
	}
	err = a.auditSessionRepo.Update(session)
	if err != nil {
		return fmt.Errorf("%w: failed to update session data: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (a *auditSessionService) GetAllSessionsPaginated(params request.AuditSessionSearchParams) ([]*response.AuditSessionResponse, int64, error) {
	sessions, total, err := a.auditSessionRepo.FindAllPaginated(params)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*response.AuditSessionResponse, 0, len(sessions))
	for _, sess := range sessions {
		result = append(result, toAuditSessionResponse(sess))
	}
	return result, total, nil
}
