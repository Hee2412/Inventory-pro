package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProductHandler struct {
	productHandler service.ProductService
}

func NewProductHandler(productHandler service.ProductService) *ProductHandler {
	return &ProductHandler{productHandler: productHandler}
}

// GetAllProducts GET /api/products
func (p *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := p.productHandler.FindActiveProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": products})
}

// GetAllProductsForAdmin GET /api/admin/products
func (p *ProductHandler) GetAllProductsForAdmin(c *gin.Context) {
	products, err := p.productHandler.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": products})
}

// GetProductById GET /api/products/:id
func (p *ProductHandler) GetProductById(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Products ID"})
		return
	}
	product, err := p.productHandler.GetProductById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product": product})
}

// CreateProduct POST /api/admin/products
func (p *ProductHandler) CreateProduct(c *gin.Context) {
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	data, err := p.productHandler.CreateProduct(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fails to create products",
			"error":   err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": "Products created successfully",
		"product": data})
}

// UpdateProduct PUT /api/admin/products/:id
func (p *ProductHandler) UpdateProduct(c *gin.Context) {
	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Input"})
		return
	}
	id, err := getIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Products ID"})
		return
	}
	if err = p.productHandler.UpdateProduct(id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Fails to update product",
			"error":   err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "Products updated successfully"})
}

// DeactivateProduct PATCH /api/admin/products/:id/deactivate
func (p *ProductHandler) DeactivateProduct(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Products ID"})
		return
	}
	if err = p.productHandler.DeactivateProduct(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Fails to deactivate product",
			"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "Products deactivated successfully"})
}

// ActivateProduct PATCH /api/admin/products/:id/activate
func (p *ProductHandler) ActivateProduct(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid Products ID"})
		return
	}
	if err = p.productHandler.ActivateProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fails to activate products",
			"error":   err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "Products activated successfully"})
}

// DeleteProduct DELETE /api/admin/products/:id
func (p *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Products ID"})
		return
	}
	if err = p.productHandler.DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fails to delete products",
			"error":   err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "Products deleted successfully"})
}

// HardDeleteProduct DELETE /api/superadmin/products/:id/hard
func (p *ProductHandler) HardDeleteProduct(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invalid Products ID"})
		return
	}
	if err = p.productHandler.HardDeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fails to delete products",
			"error":   err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "Products deleted successfully"})
}
