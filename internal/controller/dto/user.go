package dto

import (
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
)

// CreateUserRequest represents the request body for creating a user.
type CreateUserRequest struct {
	Name string `json:"name" minLength:"1" maxLength:"100" doc:"User name"`
}

// UserResponse represents the response body for a user.
type UserResponse struct {
	ID        string    `json:"id" doc:"User ID"`
	Name      string    `json:"name" doc:"User name"`
	CreatedAt time.Time `json:"createdAt" doc:"Creation timestamp"`
	UpdatedAt time.Time `json:"updatedAt" doc:"Last update timestamp"`
}

// ToUserResponse converts a domain.User to a UserResponse.
func ToUserResponse(u *domain.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
