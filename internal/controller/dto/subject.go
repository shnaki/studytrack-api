package dto

import (
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
)

// CreateSubjectRequest represents the request body for creating a subject.
type CreateSubjectRequest struct {
	Name string `json:"name" minLength:"1" maxLength:"200" doc:"Subject name"`
}

// UpdateSubjectRequest represents the request body for updating a subject.
type UpdateSubjectRequest struct {
	Name string `json:"name" minLength:"1" maxLength:"200" doc:"Subject name"`
}

// SubjectResponse represents the response body for a subject.
type SubjectResponse struct {
	ID        string    `json:"id" doc:"Subject ID"`
	UserID    string    `json:"userId" doc:"Owner user ID"`
	Name      string    `json:"name" doc:"Subject name"`
	CreatedAt time.Time `json:"createdAt" doc:"Creation timestamp"`
	UpdatedAt time.Time `json:"updatedAt" doc:"Last update timestamp"`
}

// ToSubjectResponse converts a domain.Subject to a SubjectResponse.
func ToSubjectResponse(s *domain.Subject) SubjectResponse {
	return SubjectResponse{
		ID:        s.ID,
		UserID:    s.UserID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

// ToSubjectResponseList converts a list of domain.Subject to a list of SubjectResponse.
func ToSubjectResponseList(subjects []*domain.Subject) []SubjectResponse {
	result := make([]SubjectResponse, len(subjects))
	for i, s := range subjects {
		result[i] = ToSubjectResponse(s)
	}
	return result
}
