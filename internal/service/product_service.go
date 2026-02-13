package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
)

type ProductService interface {
	GetAllProducts() ([]response.ProductResponse, error)
	FindActiveProducts() ([]response.ProductResponse, error)
	GetProductById(productId uint) (*response.ProductResponse, error)
	CreateProduct(req request.CreateProductRequest) (*response.ProductResponse, error)
	UpdateProduct(productId uint, req request.UpdateProductRequest) error
	DeactivateProduct(productId uint) error
	ActivateProduct(productId uint) error
	DeleteProduct(productId uint) error
	HardDeleteProduct(productId uint) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func toProductResponse(product *domain.Product) response.ProductResponse {
	return response.ProductResponse{
		ID:          product.ID,
		ProductName: product.ProductName,
		ProductCode: product.ProductCode,
		Unit:        product.Unit,
		MOQ:         product.MOQ,
		OM:          product.OM,
		Type:        product.Type,
		OrderCycle:  product.OrderCycle,
		AuditCycle:  product.AuditCycle,
		IsActive:    product.IsActive,
	}
}

func (p *productService) FindActiveProducts() ([]response.ProductResponse, error) {
	products, err := p.repo.FindActiveProducts()
	if err != nil {
		return nil, err
	}
	var result []response.ProductResponse
	for _, product := range products {
		result = append(result, toProductResponse(product))
	}
	return result, nil
}
func (p *productService) GetAllProducts() ([]response.ProductResponse, error) {
	product, err := p.repo.FindAll()
	if err != nil {
		return nil, err
	}
	var products []response.ProductResponse
	for _, product := range product {
		products = append(products, toProductResponse(product))
	}
	return products, nil
}

func (p *productService) GetProductById(productId uint) (*response.ProductResponse, error) {
	product, err := p.repo.FindById(productId)
	if err != nil {
		return nil, err
	}
	result := toProductResponse(product)
	return &result, nil
}

func (p *productService) CreateProduct(req request.CreateProductRequest) (*response.ProductResponse, error) {
	newProduct := &domain.Product{
		ProductName: req.ProductName,
		Unit:        req.Unit,
		MOQ:         req.MOQ,
		OM:          req.OM,
		Type:        req.Type,
		OrderCycle:  req.OrderCycle,
		AuditCycle:  req.AuditCycle,
		IsActive:    req.IsActive,
	}
	err := p.repo.Create(newProduct)
	if err != nil {
		return nil, err
	}
	result := toProductResponse(newProduct)
	return &result, nil
}

func (p *productService) UpdateProduct(productId uint, req request.UpdateProductRequest) error {
	product, err := p.repo.FindById(productId)
	if err != nil {
		return errors.New("product not found")
	}
	if req.ProductName != nil {
		product.ProductName = *req.ProductName
	}
	if req.Unit != nil {
		product.Unit = *req.Unit
	}
	if req.MOQ != nil {
		product.MOQ = *req.MOQ
	}
	if req.OM != nil {
		product.OM = *req.OM
	}
	if req.Type != nil {
		product.Type = *req.Type
	}
	if req.OrderCycle != nil {
		product.OrderCycle = *req.OrderCycle
	}
	if req.AuditCycle != nil {
		product.AuditCycle = *req.AuditCycle
	}
	if req.IsActive != nil {
		product.IsActive = *req.IsActive
	}
	return p.repo.Update(product)
}

func (p *productService) DeactivateProduct(productId uint) error {
	product, err := p.repo.FindById(productId)
	if err != nil {
		return errors.New("product not found")
	}
	if product.IsActive == false {
		return errors.New("product is not active")
	}
	product.IsActive = false
	return p.repo.Update(product)
}

func (p *productService) ActivateProduct(productId uint) error {
	product, err := p.repo.FindById(productId)
	if err != nil {
		return errors.New("product not found")
	}
	if product.IsActive == true {
		return errors.New("product is active")
	}
	product.IsActive = true
	return p.repo.Update(product)
}

func (p *productService) DeleteProduct(productId uint) error {
	_, err := p.repo.FindById(productId)
	if err != nil {
		return errors.New("product not found")
	}
	return p.repo.Delete(productId)
}

func (p *productService) HardDeleteProduct(productId uint) error {
	_, err := p.repo.FindById(productId)
	if err != nil {
		return errors.New("product not found")
	}
	return p.repo.HardDelete(productId)
}
