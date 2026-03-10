package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/service"
	"github.com/gin-gonic/gin"
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
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, products)
}

// GetAllProductsForAdmin GET /api/admin/products
func (p *ProductHandler) GetAllProductsForAdmin(c *gin.Context) {
	products, err := p.productHandler.GetAllProducts()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, products)
}

// GetProductById GET /api/products/:id
func (p *ProductHandler) GetProductById(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	product, err := p.productHandler.GetProductById(id)
	if err != nil {
		response.NotFound(c, "Product not found")
		return
	}
	response.Success(c, product)
}

// CreateProduct POST /api/admin/products
func (p *ProductHandler) CreateProduct(c *gin.Context) {
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	data, err := p.productHandler.CreateProduct(req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

// UpdateProduct PUT /api/admin/products/:id
func (p *ProductHandler) UpdateProduct(c *gin.Context) {
	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	id, err := getIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err = p.productHandler.UpdateProduct(id, req); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Message(c, "Product updated")
}

// DeactivateProduct PATCH /api/admin/products/:id/deactivate
func (p *ProductHandler) DeactivateProduct(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err = p.productHandler.DeactivateProduct(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Message(c, "Product deactivated")
}

// ActivateProduct PATCH /api/admin/products/:id/activate
func (p *ProductHandler) ActivateProduct(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err = p.productHandler.ActivateProduct(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Message(c, "Product activated")
}

// DeleteProduct DELETE /api/admin/products/:id
func (p *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err = p.productHandler.DeleteProduct(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Message(c, "Product deleted")
}

// HardDeleteProduct DELETE /api/superadmin/products/:id/hard
func (p *ProductHandler) HardDeleteProduct(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err = p.productHandler.HardDeleteProduct(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Message(c, "Product deleted")
}
