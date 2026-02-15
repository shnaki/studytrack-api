package domain_test

import (
	"strings"
	"testing"
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
)

func TestNewNote_Valid(t *testing.T) {
	note, err := domain.NewNote("note-1", "proj-1", "user-1", "My Note", "some content", []string{"go", "api"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.ID != "note-1" {
		t.Errorf("expected ID 'note-1', got '%s'", note.ID)
	}
	if note.ProjectID != "proj-1" {
		t.Errorf("expected ProjectID 'proj-1', got '%s'", note.ProjectID)
	}
	if note.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", note.UserID)
	}
	if note.Title != "My Note" {
		t.Errorf("expected Title 'My Note', got '%s'", note.Title)
	}
	if note.Content != "some content" {
		t.Errorf("expected Content 'some content', got '%s'", note.Content)
	}
	if len(note.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(note.Tags))
	}
	if note.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if note.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestNewNote_EmptyTitle(t *testing.T) {
	_, err := domain.NewNote("note-1", "proj-1", "user-1", "", "content", nil)
	if err == nil {
		t.Fatal("expected error for empty title")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestNewNote_TitleTooLong(t *testing.T) {
	longTitle := strings.Repeat("a", 201)
	_, err := domain.NewNote("note-1", "proj-1", "user-1", longTitle, "", nil)
	if err == nil {
		t.Fatal("expected error for title too long")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestNewNote_ContentTooLong(t *testing.T) {
	longContent := strings.Repeat("a", 10001)
	_, err := domain.NewNote("note-1", "proj-1", "user-1", "Title", longContent, nil)
	if err == nil {
		t.Fatal("expected error for content too long")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestNewNote_TooManyTags(t *testing.T) {
	tags := make([]string, 11)
	for i := range tags {
		tags[i] = "tag"
	}
	_, err := domain.NewNote("note-1", "proj-1", "user-1", "Title", "", tags)
	if err == nil {
		t.Fatal("expected error for too many tags")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestNewNote_TagTooLong(t *testing.T) {
	longTag := strings.Repeat("a", 51)
	_, err := domain.NewNote("note-1", "proj-1", "user-1", "Title", "", []string{longTag})
	if err == nil {
		t.Fatal("expected error for tag too long")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestNewNote_EmptyProjectID(t *testing.T) {
	_, err := domain.NewNote("note-1", "", "user-1", "Title", "", nil)
	if err == nil {
		t.Fatal("expected error for empty project ID")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestNewNote_EmptyUserID(t *testing.T) {
	_, err := domain.NewNote("note-1", "proj-1", "", "Title", "", nil)
	if err == nil {
		t.Fatal("expected error for empty user ID")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestReconstructNote(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)

	note := domain.ReconstructNote("note-1", "proj-1", "user-1", "Title", "Content", []string{"tag1"}, createdAt, updatedAt)

	if note.ID != "note-1" {
		t.Errorf("expected ID 'note-1', got '%s'", note.ID)
	}
	if note.ProjectID != "proj-1" {
		t.Errorf("expected ProjectID 'proj-1', got '%s'", note.ProjectID)
	}
	if note.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", note.UserID)
	}
	if note.Title != "Title" {
		t.Errorf("expected Title 'Title', got '%s'", note.Title)
	}
	if !note.CreatedAt.Equal(createdAt) {
		t.Errorf("expected CreatedAt %v, got %v", createdAt, note.CreatedAt)
	}
	if !note.UpdatedAt.Equal(updatedAt) {
		t.Errorf("expected UpdatedAt %v, got %v", updatedAt, note.UpdatedAt)
	}
}

func TestNote_Update_Valid(t *testing.T) {
	note, _ := domain.NewNote("note-1", "proj-1", "user-1", "Old Title", "old content", []string{"old"})
	oldUpdatedAt := note.UpdatedAt

	err := note.Update("New Title", "new content", []string{"new", "updated"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.Title != "New Title" {
		t.Errorf("expected Title 'New Title', got '%s'", note.Title)
	}
	if note.Content != "new content" {
		t.Errorf("expected Content 'new content', got '%s'", note.Content)
	}
	if len(note.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(note.Tags))
	}
	if note.UpdatedAt.Before(oldUpdatedAt) {
		t.Error("expected UpdatedAt to be updated")
	}
}

func TestNote_Update_EmptyTitle(t *testing.T) {
	note, _ := domain.NewNote("note-1", "proj-1", "user-1", "Title", "content", nil)

	err := note.Update("", "content", nil)
	if err == nil {
		t.Fatal("expected error for empty title")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
	if note.Title != "Title" {
		t.Errorf("expected title to remain 'Title', got '%s'", note.Title)
	}
}

func TestNewNote_EmptyContent(t *testing.T) {
	note, err := domain.NewNote("note-1", "proj-1", "user-1", "Title", "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.Content != "" {
		t.Errorf("expected empty content, got '%s'", note.Content)
	}
}

func TestNewNote_NilTags(t *testing.T) {
	note, err := domain.NewNote("note-1", "proj-1", "user-1", "Title", "", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if note.Tags != nil {
		t.Errorf("expected nil tags, got %v", note.Tags)
	}
}
