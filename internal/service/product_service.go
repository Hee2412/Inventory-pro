package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
	"github.com/jinzhu/copier"
	"log"
)

type ProductService interface {
	GetAllProducts() ([]response.ProductResponse, error)
	GetProductById(productId uint) (*response.ProductResponse, error)
	CreateProduct(req request.CreateProductRequest) (*response.ProductResponse, error)
	UpdateProduct(productId uint, req request.UpdateProductRequest) error
	DeactivateProduct(productId uint) error
	ActivateProducts(productId uint) error
	DeleteProducts(productId uint) error
	HardDeleteProducts(productId uint) error
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
	err = copier.Copy(product, req)
	if err != nil {
		log.Printf("copy error: %v", err)
		return errors.New("error 500")
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

func (p *productService) ActivateProducts(productId uint) error {
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

func (p *productService) DeleteProducts(productId uint) error {
	_, err := p.repo.FindById(productId)
	if err != nil {
		return errors.New("product not found")
	}
	return p.repo.Delete(productId)
}

func (p *productService) HardDeleteProducts(productId uint) error {
	_, err := p.repo.FindById(productId)
	if err != nil {
		return errors.New("product not found")
	}
	return p.repo.HardDelete(productId)
}
