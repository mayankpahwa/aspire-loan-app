package ahttp

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	ID       string `json:"id" validate:"min=3,max=10"`
	Password string `json:"password" validate:"min=6,max=15"`
}

// CreateUserResponse represents the response while creating a new user
type CreateUserResponse struct {
	ID string `json:"id"`
}
