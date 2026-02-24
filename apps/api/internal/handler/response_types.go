package handler

// ProductResponse is used for API documentation
// It mirrors db.Product but with plain Go types swag understands
type ProductResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Price       string `json:"price"`
	Stock       int32  `json:"stock"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}
