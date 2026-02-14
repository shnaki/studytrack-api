package dto

import (
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
)

type CreateSubjectRequest struct {
	Name string `json:"name" minLength:"1" maxLength:"200" doc:"Subject name"`
}

type UpdateSubjectRequest struct {
	Name string `json:"name" minLength:"1" maxLength:"200" doc:"Subject name"`
}

type SubjectResponse struct {
	ID        string    `json:"id" doc:"Subject ID"`
	UserID    string    `json:"userId" doc:"Owner user ID"`
	Name      string    `json:"name" doc:"Subject name"`
	CreatedAt time.Time `json:"createdAt" doc:"Creation timestamp"`
	UpdatedAt time.Time `json:"updatedAt" doc:"Last update timestamp"`
}

func ToSubjectResponse(s *domain.Subject) SubjectResponse {
	return SubjectResponse{
		ID:        s.ID,
		UserID:    s.UserID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func ToSubjectResponseList(subjects []*domain.Subject) []SubjectResponse {
	result := make([]SubjectResponse, len(subjects))
	for i, s := range subjects {
		result[i] = ToSubjectResponse(s)
	}
	return result
}
