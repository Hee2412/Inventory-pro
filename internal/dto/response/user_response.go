package response

import "time"

type UserResponse struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Role      string     `json:"role"`
	StoreName string     `json:"store_name"`
	StoreCode string     `json:"store_code"`
	IsActive  bool       `json:"is_active"`
	CreateAt  time.Time  `json:"create_at"`
	LastLogin *time.Time `json:"last_login"`
}
