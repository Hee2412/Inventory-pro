package request

type UpdateUserRequest struct {
	StoreName *string `json:"store_name"`
	Username  *string `json:"username"`
	Password  *string `json:"password"`
	Role      *string `json:"role"`
	IsActive  *bool   `json:"is_active"`
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	Role      string `json:"role" binding:"required,oneof=admin store"`
	StoreName string `json:"store_name"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserSearchParams struct {
	Search   string `form:"search"`
	Role     string `form:"role"`
	IsActive *bool  `form:"is_active"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

func (p *UserSearchParams) SetDefaults() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 20
	}
}
