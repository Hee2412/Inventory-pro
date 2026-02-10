package request

type UpdateUserRequest struct {
	StoreName *string `json:"store_name"`
	Username  *string `json:"username"`
	Password  *string `json:"password"`
	Role      *string `json:"role"`
	IsActive  *bool   `json:"is_active"`
}
