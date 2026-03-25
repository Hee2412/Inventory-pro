package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type OrderSessionService interface {
	CreateSession(req request.CreateOrderSessionRequest, createBy uint) (*response.OrderSessionResponse, error)
	GetAllSessions(params request.SessionSearchParams) ([]*response.OrderSessionResponse, int64, error)
	GetSessionById(sessionId uint) (*response.OrderSessionDetailResponse, error)
	AddProductToSession(req request.AddProductToSessionRequest) (*response.AddProductResponse, error)
	RemoveProductFromSession(sessionId uint, productId uint) error
	CloseSession(sessionId uint) error
}

type orderSessionService struct {
	orderSessionRepo        repository.OrderSessionRepository
	productRepo             repository.ProductRepository
	orderSessionProductRepo repository.OrderSessionProductRepository
}

func NewOrderSessionService(
	repo repository.OrderSessionRepository,
	productRepo repository.ProductRepository,
	sessionProductRepo repository.OrderSessionProductRepository,
) OrderSessionService {
	return &orderSessionService{
		orderSessionRepo:        repo,
		productRepo:             productRepo,
		orderSessionProductRepo: sessionProductRepo}
}

func (o *orderSessionService) CreateSession(req request.CreateOrderSessionRequest, createdBy uint) (*response.OrderSessionResponse, error) {
	if req.Deadline.After(req.DeliveryDate) {
		return nil, fmt.Errorf("%w: invalid delivery date", domain.ErrInvalidInput)
	}
	title := req.Title
	if title == "" {
		cycleUpper := strings.ToUpper(req.OrderCycle)
		datePart := req.Deadline.Format("0201")
		title = fmt.Sprintf("%s %s", cycleUpper, datePart)
	}
	newOrderSession := &domain.OrderSession{
		Title:        title,
		OrderCycle:   req.OrderCycle,
		Status:       "OPEN",
		Deadline:     req.Deadline,
		DeliveryDate: req.DeliveryDate,
		CreatedBy:    createdBy,
	}
	err := o.orderSessionRepo.Create(newOrderSession)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create order session: %v", domain.ErrDatabase, err)
	}
	result := toOrderSessionResponse(newOrderSession)
	return result, nil
}

func (o *orderSessionService) GetAllSessions(params request.SessionSearchParams) ([]*response.OrderSessionResponse, int64, error) {
	sessions, total, err := o.orderSessionRepo.FindAllPaginated(params)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: failed to find sessions: %v", domain.ErrDatabase, err)
	}
	result := make([]*response.OrderSessionResponse, 0, len(sessions))
	for _, session := range sessions {
		result = append(result, toOrderSessionResponse(session))
	}
	return result, total, nil
}

func (o *orderSessionService) GetSessionById(sessionId uint) (*response.OrderSessionDetailResponse, error) {
	session, err := o.orderSessionRepo.FindById(sessionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%w: failed to get session: %v", domain.ErrDatabase, sessionId)
	}
	sessionProduct, err := o.orderSessionProductRepo.FindBySessionId(sessionId)
	if err != nil {
		return &response.OrderSessionDetailResponse{
			Session:  toOrderSessionResponse(session),
			Products: []*response.ProductResponse{},
		}, nil
	}
	var products []*response.ProductResponse
	for _, sp := range sessionProduct {
		product, err := o.productRepo.FindById(sp.ProductID)
		if err != nil {
			continue
		}
		products = append(products, toProductResponse(product))
	}
	return &response.OrderSessionDetailResponse{
		Session:  toOrderSessionResponse(session),
		Products: products,
	}, nil
}

func (o *orderSessionService) AddProductToSession(req request.AddProductToSessionRequest) (*response.AddProductResponse, error) {
	session, err := o.orderSessionRepo.FindById(req.SessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%w: failed to get session: %v", domain.ErrDatabase, session.ID)
	}
	if session.Status != "OPEN" {
		return nil, fmt.Errorf("%w: cannot add product from a closed session", domain.ErrSessionClosed)
	}

	var errs []string
	var sessionProduct []*domain.OrderSessionProduct

	for _, productID := range req.ProductID {
		_, err := o.productRepo.FindById(productID)
		if err != nil {
			errs = append(errs, fmt.Sprintf("product id %d not found", productID))
			continue
		}
		existing, _ := o.orderSessionProductRepo.FindBySessionAndProduct(session.ID, productID)
		if existing != nil {
			errs = append(errs, fmt.Sprintf("product id %d already exists", productID))
			continue
		}
		sessionProduct = append(sessionProduct, &domain.OrderSessionProduct{
			SessionID: session.ID,
			ProductID: productID,
		})
	}
	if len(sessionProduct) > 0 {
		err := o.orderSessionProductRepo.Create(sessionProduct...)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to add product to session: %v", domain.ErrDatabase, err)
		}
	}
	return &response.AddProductResponse{
		Added:   len(sessionProduct),
		Skipped: len(errs),
		Errs:    errs,
	}, nil
}

func (o *orderSessionService) RemoveProductFromSession(sessionId uint, productId uint) error {
	session, err := o.orderSessionRepo.FindById(sessionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrSessionNotFound
		}
		return fmt.Errorf("%w: failed to fetch session: %v", domain.ErrDatabase, sessionId)
	}
	if session.Status != "OPEN" {
		return fmt.Errorf("%w: cannot remove product from a closed session", domain.ErrSessionClosed)
	}
	sessionProduct, err := o.orderSessionProductRepo.FindBySessionAndProduct(sessionId, productId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: product is not assigned to this session", domain.ErrNotFound)
		}
		return fmt.Errorf("%w: failed to verify session product: %v", domain.ErrDatabase, err)
	}
	err = o.orderSessionProductRepo.Delete(sessionProduct.ID)
	if err != nil {
		return fmt.Errorf("%w: failed to remove product from session: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (o *orderSessionService) CloseSession(sessionId uint) error {
	session, err := o.orderSessionRepo.FindById(sessionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrSessionNotFound
		}
		return fmt.Errorf("%w: failed to get session: %v", domain.ErrDatabase, sessionId)
	}
	if session.Status != "OPEN" {
		return fmt.Errorf("%w: session closed", domain.ErrSessionClosed)
	}
	session.Status = "CLOSED"
	err = o.orderSessionRepo.Update(session)
	if err != nil {
		return fmt.Errorf("%w failed to update session", domain.ErrInternalServer)
	}
	return nil
}
