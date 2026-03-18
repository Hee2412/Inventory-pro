package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/service"
	"Inventory-pro/pkg/pagination"
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
		response.HandleError(c, err)
		return
	}
	response.Success(c, products)
}

// GetAllProductsForAdmin GET /api/admin/products
func (p *ProductHandler) GetAllProductsForAdmin(c *gin.Context) {
	//check request
	var params request.ProductSearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.HandleError(c, err)
		return
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 20
	}
	products, total, err := p.productHandler.GetAllProductsPaginated(params)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	paginationResponse := pagination.NewResponse(products, params.Page, params.Limit, total)
	response.Success(c, paginationResponse)
}

// GetProductById GET /api/products/:id
func (p *ProductHandler) GetProductById(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	product, err := p.productHandler.GetProductById(id)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, product)
}

// CreateProduct POST /api/admin/products
func (p *ProductHandler) CreateProduct(c *gin.Context) {
	var req request.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	data, err := p.productHandler.CreateProduct(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, data)
}

// UpdateProduct PUT /api/admin/products/:id
func (p *ProductHandler) UpdateProduct(c *gin.Context) {
	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	id, err := getIDParam(c, "id")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	if err = p.productHandler.UpdateProduct(id, req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Product updated")
}

// DeactivateProduct PATCH /api/admin/products/:id/deactivate
func (p *ProductHandler) DeactivateProduct(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	if err = p.productHandler.DeactivateProduct(id); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Product deactivated")
}

// ActivateProduct PATCH /api/admin/products/:id/activate
func (p *ProductHandler) ActivateProduct(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	if err = p.productHandler.ActivateProduct(id); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Product activated")
}

// DeleteProduct DELETE /api/admin/products/:id
func (p *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	if err = p.productHandler.DeleteProduct(id); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Product deleted")
}

// HardDeleteProduct DELETE /api/superadmin/products/:id/hard
func (p *ProductHandler) HardDeleteProduct(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	if err = p.productHandler.HardDeleteProduct(id); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Product deleted")
}
