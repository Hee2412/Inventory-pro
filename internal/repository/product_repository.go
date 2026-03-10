package repository

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/pkg/pagination"
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
	FindAllPaginated(page, limit int) ([]*domain.Product, int64, error)
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

func (p *productRepository) FindAllPaginated(page, limit int) ([]*domain.Product, int64, error) {
	var products []*domain.Product
	var total int64
	if err := p.db.Model(&domain.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := p.db.Scopes(pagination.Paginate(page, limit)).Find(&products).Error
	return products, total, err
}
