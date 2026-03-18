package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type ProductService interface {
	GetAllProducts() ([]*response.ProductResponse, error)
	FindActiveProducts() ([]*response.ProductResponse, error)
	GetProductById(productId uint) (*response.ProductResponse, error)
	CreateProduct(req request.CreateProductRequest) (*response.ProductResponse, error)
	UpdateProduct(productId uint, req request.UpdateProductRequest) error
	DeactivateProduct(productId uint) error
	ActivateProduct(productId uint) error
	DeleteProduct(productId uint) error
	HardDeleteProduct(productId uint) error
	GetAllProductsPaginated(params request.ProductSearchParams) ([]*response.ProductResponse, int64, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{productRepo: repo}
}

func (p *productService) FindActiveProducts() ([]*response.ProductResponse, error) {
	products, err := p.productRepo.FindActiveProducts()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}
	}
	var result []*response.ProductResponse
	for _, product := range products {
		result = append(result, toProductResponse(product))
	}
	return result, nil
}
func (p *productService) GetAllProducts() ([]*response.ProductResponse, error) {
	product, err := p.productRepo.FindAll()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("%w: failed to get product: %v", domain.ErrDatabase, err)
	}
	var products []*response.ProductResponse
	for _, product := range product {
		products = append(products, toProductResponse(product))
	}
	return products, nil
}

func (p *productService) GetProductById(productId uint) (*response.ProductResponse, error) {
	product, err := p.productRepo.FindById(productId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("%w: failed to get product: %v", domain.ErrDatabase, productId)
	}
	result := toProductResponse(product)
	return result, nil
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
		IsActive:    true,
	}
	err := p.productRepo.Create(newProduct)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to create product: %v", domain.ErrInternalServer, err)
	}
	result := toProductResponse(newProduct)
	return result, nil
}

func (p *productService) UpdateProduct(productId uint, req request.UpdateProductRequest) error {
	product, err := p.productRepo.FindById(productId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrProductNotFound
		}
		return fmt.Errorf("%w: failed to get product: %v", domain.ErrDatabase, productId)
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
	err = p.productRepo.Update(product)
	if err != nil {
		return fmt.Errorf("%w: failed to update product: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (p *productService) DeactivateProduct(productId uint) error {
	product, err := p.productRepo.FindById(productId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrProductNotFound
		}
		return fmt.Errorf("%w: failed to get product: %v", domain.ErrDatabase, productId)
	}
	if product.IsActive == false {
		return fmt.Errorf("%w: product is inactive", domain.ErrInvalidInput)
	}
	product.IsActive = false
	err = p.productRepo.Update(product)
	if err != nil {
		return fmt.Errorf("%w: failed to deactivate product: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (p *productService) ActivateProduct(productId uint) error {
	product, err := p.productRepo.FindById(productId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrProductNotFound
		}
		return fmt.Errorf("%w: failed to get product: %v", domain.ErrDatabase, productId)
	}
	if product.IsActive == true {
		return fmt.Errorf("%w: product is active", domain.ErrInvalidInput)
	}
	product.IsActive = true
	err = p.productRepo.Update(product)
	if err != nil {
		return fmt.Errorf("%w: failed to activate product: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (p *productService) DeleteProduct(productId uint) error {
	_, err := p.productRepo.FindById(productId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrProductNotFound
		}
		return fmt.Errorf("%w: failed to get product: %v", domain.ErrDatabase, productId)
	}
	err = p.productRepo.Delete(productId)
	if err != nil {
		return fmt.Errorf("%w: failed to delete product: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (p *productService) HardDeleteProduct(productId uint) error {
	_, err := p.productRepo.FindById(productId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrProductNotFound
		}
		return fmt.Errorf("%w: failed to get product: %v", domain.ErrDatabase, productId)
	}
	err = p.productRepo.HardDelete(productId)
	if err != nil {
		return fmt.Errorf("%w: failed to delete product: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (p *productService) GetAllProductsPaginated(params request.ProductSearchParams) ([]*response.ProductResponse, int64, error) {
	products, total, err := p.productRepo.FindAllPaginated(params)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*response.ProductResponse, 0, len(products))
	for _, product := range products {
		result = append(result, toProductResponse(product))
	}
	return result, total, nil
}
