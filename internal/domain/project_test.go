package domain_test

import (
	"strings"
	"testing"
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
)

func TestNewProject_Valid(t *testing.T) {
	project, err := domain.NewProject("proj-1", "user-1", "Mathematics")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if project.ID != "proj-1" {
		t.Errorf("expected ID 'proj-1', got '%s'", project.ID)
	}
	if project.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", project.UserID)
	}
	if project.Name != "Mathematics" {
		t.Errorf("expected Name 'Mathematics', got '%s'", project.Name)
	}
	if project.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if project.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestNewProject_EmptyName(t *testing.T) {
	_, err := domain.NewProject("proj-1", "user-1", "")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestNewProject_NameTooLong(t *testing.T) {
	longName := strings.Repeat("a", 201)
	_, err := domain.NewProject("proj-1", "user-1", longName)
	if err == nil {
		t.Fatal("expected error for name too long")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestNewProject_NameExactly200(t *testing.T) {
	name200 := strings.Repeat("a", 200)
	project, err := domain.NewProject("proj-1", "user-1", name200)
	if err != nil {
		t.Fatalf("unexpected error for 200-char name: %v", err)
	}
	if project.Name != name200 {
		t.Errorf("expected name of length 200, got length %d", len(project.Name))
	}
}

func TestNewProject_EmptyUserID(t *testing.T) {
	_, err := domain.NewProject("proj-1", "", "Mathematics")
	if err == nil {
		t.Fatal("expected error for empty user ID")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestReconstructProject(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)

	project := domain.ReconstructProject("proj-1", "user-1", "Physics", createdAt, updatedAt)

	if project.ID != "proj-1" {
		t.Errorf("expected ID 'proj-1', got '%s'", project.ID)
	}
	if project.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", project.UserID)
	}
	if project.Name != "Physics" {
		t.Errorf("expected Name 'Physics', got '%s'", project.Name)
	}
	if !project.CreatedAt.Equal(createdAt) {
		t.Errorf("expected CreatedAt %v, got %v", createdAt, project.CreatedAt)
	}
	if !project.UpdatedAt.Equal(updatedAt) {
		t.Errorf("expected UpdatedAt %v, got %v", updatedAt, project.UpdatedAt)
	}
}

func TestProject_UpdateName_Valid(t *testing.T) {
	project, err := domain.NewProject("proj-1", "user-1", "Math")
	if err != nil {
		t.Fatalf("unexpected error creating project: %v", err)
	}

	oldUpdatedAt := project.UpdatedAt

	err = project.UpdateName("Mathematics")
	if err != nil {
		t.Fatalf("unexpected error updating name: %v", err)
	}
	if project.Name != "Mathematics" {
		t.Errorf("expected Name 'Mathematics', got '%s'", project.Name)
	}
	if project.UpdatedAt.Before(oldUpdatedAt) {
		t.Error("expected UpdatedAt to be updated")
	}
}

func TestProject_UpdateName_EmptyName(t *testing.T) {
	project, err := domain.NewProject("proj-1", "user-1", "Math")
	if err != nil {
		t.Fatalf("unexpected error creating project: %v", err)
	}

	err = project.UpdateName("")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
	if project.Name != "Math" {
		t.Errorf("expected name to remain 'Math', got '%s'", project.Name)
	}
}

func TestProject_UpdateName_TooLong(t *testing.T) {
	project, err := domain.NewProject("proj-1", "user-1", "Math")
	if err != nil {
		t.Fatalf("unexpected error creating project: %v", err)
	}

	longName := strings.Repeat("x", 201)
	err = project.UpdateName(longName)
	if err == nil {
		t.Fatal("expected error for name too long")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
	if project.Name != "Math" {
		t.Errorf("expected name to remain 'Math', got '%s'", project.Name)
	}
}
