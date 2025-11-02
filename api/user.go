package api

import "time"

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=admin employee"`
	FullName string `json:"full_name" binding:"required,min=3"`
}

// Create a user response that will provide to the response
type userResponse struct {
	Username         string    `json:"username"`
	FullName         string    `json:"full_name"`
	Email            string    `json:"email"`
	Role             string    `json:"role"`
	PasswordChangeAt time.Time `json:"password_change_at"`
	CreatedAt        time.Time `json:"created_at"`
}
