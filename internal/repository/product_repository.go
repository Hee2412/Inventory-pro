package repository

import (
	"Inventory-pro/internal/domain"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *domain.Product) error
	FindById(id uint) (*domain.Product, error)
	FindByProductName(productName string) (*domain.Product, error)
	FindByProductCode(productCode string) (*domain.Product, error)
	FindAll() ([]domain.Product, error)
	Update(product *domain.Product) error
	Delete(product *domain.Product) error
	HardDelete(product *domain.Product) error

	FindActiveProduct() ([]*domain.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func (p *productRepository) Create(product *domain.Product) error {
	return p.db.Create(product).Error
}

func (p *productRepository) FindById(id uint) (*domain.Product, error) {
	var product domain.Product
	err := p.db.First(&product, id).Error
	return &product, err
}

func (p *productRepository) FindByProductName(productName string) (*domain.Product, error) {
	var product domain.Product
	err := p.db.Where("product_name = ?", productName).First(&product).Error
	return &product, err
}
func (p *productRepository) FindByProductCode(productCode string) (*domain.Product, error) {
	var product domain.Product
	err := p.db.Where("product_code = ?", productCode).First(&product).Error
	return &product, err
}

func (p *productRepository) FindAll() ([]domain.Product, error) {
	var product []domain.Product
	err := p.db.Find(&product).Error
	return product, err
}

func (p *productRepository) Update(product *domain.Product) error {
	return p.db.Updates(product).Error
}

func (p *productRepository) Delete(product *domain.Product) error {
	return p.db.Delete(product).Error
}

func (p *productRepository) HardDelete(product *domain.Product) error {
	return p.db.Unscoped().Delete(product).Error
}

func (p *productRepository) FindActiveProduct() ([]*domain.Product, error) {
	var product []*domain.Product
	err := p.db.Where("is_active = true", true).Find(&product).Error
	return product, err
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}
