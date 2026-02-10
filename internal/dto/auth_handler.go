package dto

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

type LoginResponse struct {
	Token     string `json:"token"`
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	StoreName string `json:"store_name"`
	StoreCode string `json:"store_code"`
}
