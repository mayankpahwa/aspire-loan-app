package ahttp

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
	UserID string `json:"user_id"`
}

// CreateUserResponse represents the response while creating a new user
type CreateUserResponse struct {
	UserID string `json:"user_id"`
}
