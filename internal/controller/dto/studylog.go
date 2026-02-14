package dto

import (
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
)

// CreateStudyLogRequest represents the request body for creating a study log.
type CreateStudyLogRequest struct {
	SubjectID string    `json:"subjectId" doc:"Subject ID"`
	StudiedAt time.Time `json:"studiedAt" doc:"When the study session occurred"`
	Minutes   int       `json:"minutes" minimum:"1" maximum:"1440" doc:"Duration in minutes"`
	Note      string    `json:"note,omitempty" maxLength:"1000" doc:"Optional note"`
}

// StudyLogResponse represents the response body for a study log.
type StudyLogResponse struct {
	ID        string    `json:"id" doc:"Study log ID"`
	UserID    string    `json:"userId" doc:"User ID"`
	SubjectID string    `json:"subjectId" doc:"Subject ID"`
	StudiedAt time.Time `json:"studiedAt" doc:"When the study session occurred"`
	Minutes   int       `json:"minutes" doc:"Duration in minutes"`
	Note      string    `json:"note" doc:"Note"`
	CreatedAt time.Time `json:"createdAt" doc:"Creation timestamp"`
}

// ToStudyLogResponse converts a domain.StudyLog to a StudyLogResponse.
func ToStudyLogResponse(l *domain.StudyLog) StudyLogResponse {
	return StudyLogResponse{
		ID:        l.ID,
		UserID:    l.UserID,
		SubjectID: l.SubjectID,
		StudiedAt: l.StudiedAt,
		Minutes:   l.Minutes,
		Note:      l.Note,
		CreatedAt: l.CreatedAt,
	}
}

// ToStudyLogResponseList converts a list of domain.StudyLog to a list of StudyLogResponse.
func ToStudyLogResponseList(logs []*domain.StudyLog) []StudyLogResponse {
	result := make([]StudyLogResponse, len(logs))
	for i, l := range logs {
		result[i] = ToStudyLogResponse(l)
	}
	return result
}
