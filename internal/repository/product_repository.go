package repository

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *domain.Product) error
	FindById(id uint) (*domain.Product, error)
	FindByProductName(productName string) (*domain.Product, error)
	FindByProductCode(productCode string) (*domain.Product, error)
	FindAll() ([]*domain.Product, error)
	Update(product *domain.Product) error
	Delete(id uint) error
	HardDelete(id uint) error
	FindActiveProducts() ([]*domain.Product, error)
	FindAllPaginated(params request.ProductSearchParams) ([]*domain.Product, int64, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (p *productRepository) Create(product *domain.Product) error {
	return p.db.Create(product).Error
}

func (p *productRepository) FindById(id uint) (*domain.Product, error) {
	var product domain.Product
	err := p.db.First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *productRepository) FindByProductName(productName string) (*domain.Product, error) {
	var product domain.Product
	err := p.db.Where("product_name = ?", productName).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}
func (p *productRepository) FindByProductCode(productCode string) (*domain.Product, error) {
	var product domain.Product
	err := p.db.Where("product_code = ?", productCode).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *productRepository) FindAll() ([]*domain.Product, error) {
	var products []*domain.Product
	err := p.db.Find(&products).Error
	return products, err
}

func (p *productRepository) Update(product *domain.Product) error {
	return p.db.Save(product).Error
}

func (p *productRepository) Delete(id uint) error {
	return p.db.Delete(&domain.Product{}, id).Error
}

func (p *productRepository) HardDelete(id uint) error {
	return p.db.Unscoped().Delete(&domain.Product{}, id).Error
}

func (p *productRepository) FindActiveProducts() ([]*domain.Product, error) {
	var products []*domain.Product
	err := p.db.Where("is_active = true").Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *productRepository) FindAllPaginated(params request.ProductSearchParams) ([]*domain.Product, int64, error) {
	var products []*domain.Product
	var total int64
	query := p.db.Model(&domain.Product{})
	if params.Search != "" {
		searchTerm := "%" + params.Search + "%"
		query = query.Where(
			"product_name LIKE ? OR product_code LIKE ?",
			searchTerm, searchTerm)
	}
	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}
	if params.MinPrice != nil {
		query = query.Where("min_price >= ?", *params.MinPrice)
	}
	if params.MaxPrice != nil {
		query = query.Where("max_price <= ?", *params.MaxPrice)
	}
	query.Count(&total)
	offset := (params.Page - 1) * params.Limit
	err := query.Offset(offset).Limit(params.Limit).
		Order("created_at DESC").
		Find(&products).Error
	return products, total, err
}
