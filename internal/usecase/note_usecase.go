package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase/port"
)

// NoteUsecase provides methods for managing notes.
type NoteUsecase struct {
	noteRepo    port.NoteRepository
	projectRepo port.ProjectRepository
	userRepo    port.UserRepository
}

// NewNoteUsecase creates a new NoteUsecase.
func NewNoteUsecase(noteRepo port.NoteRepository, projectRepo port.ProjectRepository, userRepo port.UserRepository) *NoteUsecase {
	return &NoteUsecase{
		noteRepo:    noteRepo,
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

// CreateNote creates a new note for a project.
func (u *NoteUsecase) CreateNote(ctx context.Context, userID, projectID, title, content string, tags []string) (*domain.Note, error) {
	if _, err := u.userRepo.FindByID(ctx, userID); err != nil {
		return nil, err
	}
	project, err := u.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, domain.ErrNotFound("project")
	}

	id := uuid.New().String()
	note, err := domain.NewNote(id, projectID, userID, title, content, tags)
	if err != nil {
		return nil, err
	}
	if err := u.noteRepo.Create(ctx, note); err != nil {
		return nil, err
	}
	return note, nil
}

// GetNote returns a note by ID.
func (u *NoteUsecase) GetNote(ctx context.Context, id string) (*domain.Note, error) {
	return u.noteRepo.FindByID(ctx, id)
}

// ListNotes returns all notes for a project.
func (u *NoteUsecase) ListNotes(ctx context.Context, userID, projectID string) ([]*domain.Note, error) {
	if _, err := u.userRepo.FindByID(ctx, userID); err != nil {
		return nil, err
	}
	project, err := u.projectRepo.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.UserID != userID {
		return nil, domain.ErrNotFound("project")
	}
	return u.noteRepo.FindByProjectID(ctx, projectID)
}

// UpdateNote updates an existing note.
func (u *NoteUsecase) UpdateNote(ctx context.Context, id, title, content string, tags []string) (*domain.Note, error) {
	note, err := u.noteRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := note.Update(title, content, tags); err != nil {
		return nil, err
	}
	if err := u.noteRepo.Update(ctx, note); err != nil {
		return nil, err
	}
	return note, nil
}

// DeleteNote deletes a note by ID.
func (u *NoteUsecase) DeleteNote(ctx context.Context, id string) error {
	if _, err := u.noteRepo.FindByID(ctx, id); err != nil {
		return err
	}
	return u.noteRepo.Delete(ctx, id)
}
