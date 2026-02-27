package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
	"time"
)

type AuditSessionService interface {
	CreateAuditSession(req request.CreateAuditSessionRequest, creatBy uint) (*response.AuditSessionResponse, error)
	GetAllAuditSessions() ([]response.AuditSessionResponse, error)
	GetAuditSessionByID(sessionID uint) (*response.AuditSessionDetailsResponse, error)
	AddProductToAudit(req request.AddProductToAuditRequest, creatBy uint) (*response.AuditSessionResponse, error)
	RemoveProductFromAudit(auditSessionID uint, productID uint) error
	CloseAuditSession(auditSessionID uint) error
}

type auditSessionService struct {
	auditSessionRepo repository.AuditSessionRepository
	storeAuditRepo   repository.StoreAuditReportRepository
}

func NewAuditSessionService(
	auditSessionRepo repository.AuditSessionRepository,
	storeAuditRepo repository.StoreAuditReportRepository) AuditSessionService {
	return &auditSessionService{
		auditSessionRepo: auditSessionRepo,
		storeAuditRepo:   storeAuditRepo,
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

func (a *auditSessionService) CreateAuditSession(req request.CreateAuditSessionRequest, creatBy uint) (*response.AuditSessionResponse, error) {
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
		return nil, err
	}
	//mapping
	storeMap := make(map[uint]response.StoreAuditReportResponse)
	for _, r := range reports {
		if _, exist := storeMap[r.StoreID]; !exist {
			storeName := "unknown"
			if r.Store != nil {
				storeName = r.Store.StoreName
			}
			storeMap[r.StoreID] = response.StoreAuditReportResponse{
				StoreId:   r.StoreID,
				StoreName: storeName,
				Status:    r.Status,
			}
		}
	}
	//response
	var storeResponse []response.StoreAuditReportResponse
	for _, store := range storeMap {
		storeResponse = append(storeResponse, store)
	}
	return &response.AuditSessionDetailsResponse{
		SessionInfo: toAuditSessionResponse(session),
		Report:      storeResponse,
	}, nil
}

func (a *auditSessionService) AddProductToAudit(req request.AddProductToAuditRequest, creatBy uint) (*response.AuditSessionResponse, error) {
	//session, err := a.auditSessionRepo.FindById(req.AuditSessionID)
	//if err != nil {
	//	return nil, errors.New("the session does not exist")
	//}
	//if session.Status != "OPEN" {
	//	return nil, errors.New("the session is not open")
	//}
	//var errs []string
	//var sessionProduct []*domain.

}

func (a *auditSessionService) RemoveProductFromAudit(auditSessionID uint, productID uint) error {
	//TODO implement me
	panic("implement me")
}

func (a *auditSessionService) CloseAuditSession(auditSessionID uint) error {
	//TODO implement me
	panic("implement me")
}
