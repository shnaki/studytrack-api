package dto

import (
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
)

// CreateNoteRequest represents the request body for creating a note.
type CreateNoteRequest struct {
	Title   string   `json:"title" minLength:"1" maxLength:"200" doc:"Note title"`
	Content string   `json:"content,omitempty" maxLength:"10000" doc:"Note content"`
	Tags    []string `json:"tags,omitempty" maxItems:"10" doc:"Note tags"`
}

// UpdateNoteRequest represents the request body for updating a note.
type UpdateNoteRequest struct {
	Title   string   `json:"title" minLength:"1" maxLength:"200" doc:"Note title"`
	Content string   `json:"content,omitempty" maxLength:"10000" doc:"Note content"`
	Tags    []string `json:"tags,omitempty" maxItems:"10" doc:"Note tags"`
}

// NoteResponse represents the response body for a note.
type NoteResponse struct {
	ID        string    `json:"id" doc:"Note ID"`
	ProjectID string    `json:"projectId" doc:"Project ID"`
	UserID    string    `json:"userId" doc:"Owner user ID"`
	Title     string    `json:"title" doc:"Note title"`
	Content   string    `json:"content" doc:"Note content"`
	Tags      []string  `json:"tags" doc:"Note tags"`
	CreatedAt time.Time `json:"createdAt" doc:"Creation timestamp"`
	UpdatedAt time.Time `json:"updatedAt" doc:"Last update timestamp"`
}

// ToNoteResponse converts a domain.Note to a NoteResponse.
func ToNoteResponse(n *domain.Note) NoteResponse {
	tags := n.Tags
	if tags == nil {
		tags = []string{}
	}
	return NoteResponse{
		ID:        n.ID,
		ProjectID: n.ProjectID,
		UserID:    n.UserID,
		Title:     n.Title,
		Content:   n.Content,
		Tags:      tags,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

// ToNoteResponseList converts a list of domain.Note to a list of NoteResponse.
func ToNoteResponseList(notes []*domain.Note) []NoteResponse {
	result := make([]NoteResponse, len(notes))
	for i, n := range notes {
		result[i] = ToNoteResponse(n)
	}
	return result
}
