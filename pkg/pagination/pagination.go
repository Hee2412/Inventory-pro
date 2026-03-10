package pagination

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
)

// Params - Request parameters
type Params struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

// SetDefaults - Set default values
func (p *Params) SetDefaults() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 20
	}
}

// GetOffset - Calculate offset
func (p *Params) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// Response - Pagination response
type Response struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}

// NewResponse - Create pagination response
func NewResponse(data interface{}, page, limit int, total int64) *Response {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &Response{
		Data:       data,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// Paginate - GORM scope for pagination
func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}

// ParseParams - parse pagination params from request
func ParseParams(c *gin.Context) *Params {
	var params Params
	if err := c.ShouldBindQuery(&params); err != nil {
		// return default if parse error
		return &Params{Page: 1, Limit: 20}
	}
	params.SetDefaults()
	return &params
}
