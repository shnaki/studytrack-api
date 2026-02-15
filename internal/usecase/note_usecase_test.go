package usecase_test

import (
	"context"
	"testing"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase"
)

func setupNoteTest() (*usecase.NoteUsecase, *mockUserRepository, *mockProjectRepository, *mockNoteRepository) {
	userRepo := newMockUserRepository()
	projectRepo := newMockProjectRepository()
	noteRepo := newMockNoteRepository()
	uc := usecase.NewNoteUsecase(noteRepo, projectRepo, userRepo)
	return uc, userRepo, projectRepo, noteRepo
}

func createTestProject(repo *mockProjectRepository, id, userID, name string) {
	repo.projects[id] = &domain.Project{ID: id, UserID: userID, Name: name}
}

func TestCreateNote_Success(t *testing.T) {
	uc, userRepo, projectRepo, _ := setupNoteTest()
	createTestUser(userRepo, "user-1", "Alice")
	createTestProject(projectRepo, "proj-1", "user-1", "Math")

	note, err := uc.CreateNote(context.Background(), "user-1", "proj-1", "My Note", "some content", []string{"go", "api"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.Title != "My Note" {
		t.Errorf("expected title 'My Note', got '%s'", note.Title)
	}
	if note.Content != "some content" {
		t.Errorf("expected content 'some content', got '%s'", note.Content)
	}
	if note.ProjectID != "proj-1" {
		t.Errorf("expected ProjectID 'proj-1', got '%s'", note.ProjectID)
	}
	if note.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", note.UserID)
	}
	if note.ID == "" {
		t.Error("expected ID to be generated")
	}
	if len(note.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(note.Tags))
	}
}

func TestCreateNote_UserNotFound(t *testing.T) {
	uc, _, projectRepo, _ := setupNoteTest()
	createTestProject(projectRepo, "proj-1", "user-1", "Math")

	_, err := uc.CreateNote(context.Background(), "nonexistent", "proj-1", "Title", "", nil)
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestCreateNote_ProjectNotFound(t *testing.T) {
	uc, userRepo, _, _ := setupNoteTest()
	createTestUser(userRepo, "user-1", "Alice")

	_, err := uc.CreateNote(context.Background(), "user-1", "nonexistent", "Title", "", nil)
	if err == nil {
		t.Fatal("expected error for nonexistent project")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestCreateNote_ProjectNotOwnedByUser(t *testing.T) {
	uc, userRepo, projectRepo, _ := setupNoteTest()
	createTestUser(userRepo, "user-1", "Alice")
	createTestUser(userRepo, "user-2", "Bob")
	createTestProject(projectRepo, "proj-1", "user-2", "Math")

	_, err := uc.CreateNote(context.Background(), "user-1", "proj-1", "Title", "", nil)
	if err == nil {
		t.Fatal("expected error for project not owned by user")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestCreateNote_EmptyTitle(t *testing.T) {
	uc, userRepo, projectRepo, _ := setupNoteTest()
	createTestUser(userRepo, "user-1", "Alice")
	createTestProject(projectRepo, "proj-1", "user-1", "Math")

	_, err := uc.CreateNote(context.Background(), "user-1", "proj-1", "", "", nil)
	if err == nil {
		t.Fatal("expected error for empty title")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestGetNote_Success(t *testing.T) {
	uc, userRepo, projectRepo, _ := setupNoteTest()
	createTestUser(userRepo, "user-1", "Alice")
	createTestProject(projectRepo, "proj-1", "user-1", "Math")

	created, err := uc.CreateNote(context.Background(), "user-1", "proj-1", "My Note", "content", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	note, err := uc.GetNote(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.Title != "My Note" {
		t.Errorf("expected title 'My Note', got '%s'", note.Title)
	}
}

func TestGetNote_NotFound(t *testing.T) {
	uc, _, _, _ := setupNoteTest()

	_, err := uc.GetNote(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent note")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestListNotes_Success(t *testing.T) {
	uc, userRepo, projectRepo, _ := setupNoteTest()
	createTestUser(userRepo, "user-1", "Alice")
	createTestProject(projectRepo, "proj-1", "user-1", "Math")

	_, err := uc.CreateNote(context.Background(), "user-1", "proj-1", "Note 1", "content 1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = uc.CreateNote(context.Background(), "user-1", "proj-1", "Note 2", "content 2", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	notes, err := uc.ListNotes(context.Background(), "user-1", "proj-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(notes) != 2 {
		t.Errorf("expected 2 notes, got %d", len(notes))
	}
}

func TestListNotes_UserNotFound(t *testing.T) {
	uc, _, _, _ := setupNoteTest()

	_, err := uc.ListNotes(context.Background(), "nonexistent", "proj-1")
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestListNotes_ProjectNotOwnedByUser(t *testing.T) {
	uc, userRepo, projectRepo, _ := setupNoteTest()
	createTestUser(userRepo, "user-1", "Alice")
	createTestUser(userRepo, "user-2", "Bob")
	createTestProject(projectRepo, "proj-1", "user-2", "Math")

	_, err := uc.ListNotes(context.Background(), "user-1", "proj-1")
	if err == nil {
		t.Fatal("expected error for project not owned by user")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestUpdateNote_Success(t *testing.T) {
	uc, userRepo, projectRepo, _ := setupNoteTest()
	createTestUser(userRepo, "user-1", "Alice")
	createTestProject(projectRepo, "proj-1", "user-1", "Math")

	created, err := uc.CreateNote(context.Background(), "user-1", "proj-1", "Old Title", "old content", []string{"old"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updated, err := uc.UpdateNote(context.Background(), created.ID, "New Title", "new content", []string{"new", "updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Title != "New Title" {
		t.Errorf("expected title 'New Title', got '%s'", updated.Title)
	}
	if updated.Content != "new content" {
		t.Errorf("expected content 'new content', got '%s'", updated.Content)
	}
	if len(updated.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(updated.Tags))
	}
}

func TestUpdateNote_NotFound(t *testing.T) {
	uc, _, _, _ := setupNoteTest()

	_, err := uc.UpdateNote(context.Background(), "nonexistent", "Title", "content", nil)
	if err == nil {
		t.Fatal("expected error for nonexistent note")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestUpdateNote_EmptyTitle(t *testing.T) {
	uc, userRepo, projectRepo, _ := setupNoteTest()
	createTestUser(userRepo, "user-1", "Alice")
	createTestProject(projectRepo, "proj-1", "user-1", "Math")

	created, err := uc.CreateNote(context.Background(), "user-1", "proj-1", "Title", "content", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = uc.UpdateNote(context.Background(), created.ID, "", "content", nil)
	if err == nil {
		t.Fatal("expected error for empty title")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestDeleteNote_Success(t *testing.T) {
	uc, userRepo, projectRepo, noteRepo := setupNoteTest()
	createTestUser(userRepo, "user-1", "Alice")
	createTestProject(projectRepo, "proj-1", "user-1", "Math")

	created, err := uc.CreateNote(context.Background(), "user-1", "proj-1", "Title", "content", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = uc.DeleteNote(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := noteRepo.notes[created.ID]; ok {
		t.Error("expected note to be deleted from repository")
	}
}

func TestDeleteNote_NotFound(t *testing.T) {
	uc, _, _, _ := setupNoteTest()

	err := uc.DeleteNote(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent note")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}
