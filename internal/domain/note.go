package domain

import "time"

// Note represents a note attached to a project.
type Note struct {
	ID        string
	ProjectID string
	UserID    string
	Title     string
	Content   string
	Tags      []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewNote creates a new Note entity.
func NewNote(id, projectID, userID, title, content string, tags []string) (*Note, error) {
	if projectID == "" {
		return nil, ErrValidation("project ID is required")
	}
	if userID == "" {
		return nil, ErrValidation("user ID is required")
	}
	if err := validateNoteTitle(title); err != nil {
		return nil, err
	}
	if err := validateNoteContent(content); err != nil {
		return nil, err
	}
	if err := validateNoteTags(tags); err != nil {
		return nil, err
	}
	now := time.Now()
	return &Note{
		ID:        id,
		ProjectID: projectID,
		UserID:    userID,
		Title:     title,
		Content:   content,
		Tags:      tags,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// ReconstructNote reconstructs a Note entity from existing data.
func ReconstructNote(id, projectID, userID, title, content string, tags []string, createdAt, updatedAt time.Time) *Note {
	return &Note{
		ID:        id,
		ProjectID: projectID,
		UserID:    userID,
		Title:     title,
		Content:   content,
		Tags:      tags,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// Update updates the title, content, and tags of the note.
func (n *Note) Update(title, content string, tags []string) error {
	if err := validateNoteTitle(title); err != nil {
		return err
	}
	if err := validateNoteContent(content); err != nil {
		return err
	}
	if err := validateNoteTags(tags); err != nil {
		return err
	}
	n.Title = title
	n.Content = content
	n.Tags = tags
	n.UpdatedAt = time.Now()
	return nil
}

func validateNoteTitle(title string) error {
	if title == "" {
		return ErrValidation("note title is required")
	}
	if len(title) > 200 {
		return ErrValidation("note title must be 200 characters or less")
	}
	return nil
}

func validateNoteContent(content string) error {
	if len(content) > 10000 {
		return ErrValidation("note content must be 10000 characters or less")
	}
	return nil
}

func validateNoteTags(tags []string) error {
	if len(tags) > 10 {
		return ErrValidation("note tags must be 10 or less")
	}
	for _, tag := range tags {
		if len(tag) > 50 {
			return ErrValidation("each tag must be 50 characters or less")
		}
	}
	return nil
}
