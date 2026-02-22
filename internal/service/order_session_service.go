package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
	"time"
)

type OrderSessionService interface {
	CreateSession(req request.CreateOrderSessionRequest, createBy uint) (*response.OrderSessionResponse, error)
	GetAllSessions() ([]response.OrderSessionResponse, error)
	GetSessionById(sessionId uint) (*response.StoreOrderDetailResponse, error)
	AddProductToSession(req request.AddProductToSessionRequest) error
	RemoveProductFromSession(sessionId uint, productId uint) error
	CloseSession(sessionId uint) error
}

type orderSessionService struct {
	repo               repository.OrderSessionRepository
	productRepo        repository.ProductRepository
	sessionProductRepo repository.OrderSessionProductRepository
}

func NewOrderSessionService(
	repo repository.OrderSessionRepository,
	productRepo repository.ProductRepository,
	sessionProductRepo repository.OrderSessionProductRepository,
) OrderSessionService {
	return &orderSessionService{
		repo:               repo,
		productRepo:        productRepo,
		sessionProductRepo: sessionProductRepo}
}

func toOrderSessionResponse(order *domain.OrderSession) response.OrderSessionResponse {
	return response.OrderSessionResponse{
		ID:         order.ID,
		Title:      order.Title,
		OrderCycle: order.OrderCycle,
		Status:     order.Status,
		Deadline:   order.Deadline,
		DeliveryAt: order.DeliveryDate,
		CreatedAt:  order.CreatedAt,
	}
}
func (o orderSessionService) CreateSession(req request.CreateOrderSessionRequest, createdBy uint) (*response.OrderSessionResponse, error) {
	if req.Deadline.After(req.DeliveryDate) {
		return nil, errors.New("delivery date is in the past")
	}
	newOrderSession := &domain.OrderSession{
		Title:        req.Title,
		OrderCycle:   req.OrderCycle,
		Status:       "OPEN",
		Deadline:     req.Deadline,
		DeliveryDate: req.DeliveryDate,
		CreatedBy:    createdBy,
	}
	err := o.repo.Create(newOrderSession)
	if err != nil {
		return nil, err
	}
	result := toOrderSessionResponse(newOrderSession)
	return &result, nil
}

func (o orderSessionService) GetAllSessions() ([]response.OrderSessionResponse, error) {
	sessions, err := o.repo.FindAll()
	if err != nil {
		return nil, err
	}
	var result []response.OrderSessionResponse
	for _, session := range sessions {
		result = append(result, toOrderSessionResponse(session))
	}
	return result, nil
}

func (o orderSessionService) GetSessionById(sessionId uint) (*response.StoreOrderDetailResponse, error) {
	session, err := o.repo.FindById(sessionId)
	if err != nil {
		return nil, err
	}
	result := response.StoreOrderDetailResponse{
		Order: response.StoreOrderResponse{
			ID:          0,
			SessionID:   session.ID,
			StoreID:     0,
			Status:      "",
			SubmittedAt: nil,
			ApprovedAt:  nil,
			CreatedAt:   time.Time{},
		},
		Items: []response.OrderItemResponse{},
	}
	return &result, nil
}

func (o orderSessionService) AddProductToSession(req request.AddProductToSessionRequest) error {
	session, err := o.repo.FindById(req.SessionID)
	if err != nil {
		return errors.New("session not found")
	}
	if session.Status != "OPEN" {
		return errors.New("session is closed")
	}
	_, errs := o.productRepo.FindById(req.ProductID)
	if errs != nil {
		return errors.New("product not found")
	}
	existing, _ := o.sessionProductRepo.FindBySessionAndProduct(session.ID, req.ProductID)
	if existing != nil {
		return errors.New("product already exists")
	}
	sessionProduct := &domain.OrderSessionProducts{
		SessionID: session.ID,
		ProductID: req.ProductID,
	}
	return o.sessionProductRepo.Create(sessionProduct)
}

func (o orderSessionService) RemoveProductFromSession(sessionId uint, productId uint) error {
	session, err := o.repo.FindById(sessionId)
	if err != nil {
		return errors.New("session not found")
	}
	if session.Status != "OPEN" {
		return errors.New("session is closed")
	}
	_, err = o.productRepo.FindById(productId)
	if err != nil {
		return errors.New("product not found")
	}
	sessionProduct, err := o.sessionProductRepo.FindBySessionAndProduct(sessionId, productId)
	if err != nil {
		return errors.New("product not found in session")
	}
	return o.repo.Delete(sessionProduct.ID)
}

func (o orderSessionService) CloseSession(sessionId uint) error {
	session, err := o.repo.FindById(sessionId)
	if err != nil {
		return errors.New("session not found")
	}
	if session.Status != "OPEN" {
		return errors.New("session is closed")
	}
	session.Status = "CLOSED"
	return o.repo.Update(session)
}
