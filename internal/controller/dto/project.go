package dto

import (
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
)

// CreateProjectRequest represents the request body for creating a project.
type CreateProjectRequest struct {
	Name string `json:"name" minLength:"1" maxLength:"200" doc:"Project name"`
}

// UpdateProjectRequest represents the request body for updating a project.
type UpdateProjectRequest struct {
	Name string `json:"name" minLength:"1" maxLength:"200" doc:"Project name"`
}

// ProjectResponse represents the response body for a project.
type ProjectResponse struct {
	ID        string    `json:"id" doc:"Project ID"`
	UserID    string    `json:"userId" doc:"Owner user ID"`
	Name      string    `json:"name" doc:"Project name"`
	CreatedAt time.Time `json:"createdAt" doc:"Creation timestamp"`
	UpdatedAt time.Time `json:"updatedAt" doc:"Last update timestamp"`
}

// ToProjectResponse converts a domain.Project to a ProjectResponse.
func ToProjectResponse(p *domain.Project) ProjectResponse {
	return ProjectResponse{
		ID:        p.ID,
		UserID:    p.UserID,
		Name:      p.Name,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

// ToProjectResponseList converts a list of domain.Project to a list of ProjectResponse.
func ToProjectResponseList(projects []*domain.Project) []ProjectResponse {
	result := make([]ProjectResponse, len(projects))
	for i, p := range projects {
		result[i] = ToProjectResponse(p)
	}
	return result
}
