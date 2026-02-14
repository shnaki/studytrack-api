package domain_test

import (
	"strings"
	"testing"
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
)

func TestNewSubject_Valid(t *testing.T) {
	subject, err := domain.NewSubject("sub-1", "user-1", "Mathematics")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if subject.ID != "sub-1" {
		t.Errorf("expected ID 'sub-1', got '%s'", subject.ID)
	}
	if subject.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", subject.UserID)
	}
	if subject.Name != "Mathematics" {
		t.Errorf("expected Name 'Mathematics', got '%s'", subject.Name)
	}
	if subject.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if subject.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestNewSubject_EmptyName(t *testing.T) {
	_, err := domain.NewSubject("sub-1", "user-1", "")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestNewSubject_NameTooLong(t *testing.T) {
	longName := strings.Repeat("a", 201)
	_, err := domain.NewSubject("sub-1", "user-1", longName)
	if err == nil {
		t.Fatal("expected error for name too long")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestNewSubject_NameExactly200(t *testing.T) {
	name200 := strings.Repeat("a", 200)
	subject, err := domain.NewSubject("sub-1", "user-1", name200)
	if err != nil {
		t.Fatalf("unexpected error for 200-char name: %v", err)
	}
	if subject.Name != name200 {
		t.Errorf("expected name of length 200, got length %d", len(subject.Name))
	}
}

func TestNewSubject_EmptyUserID(t *testing.T) {
	_, err := domain.NewSubject("sub-1", "", "Mathematics")
	if err == nil {
		t.Fatal("expected error for empty user ID")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestReconstructSubject(t *testing.T) {
	createdAt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 6, 15, 14, 30, 0, 0, time.UTC)

	subject := domain.ReconstructSubject("sub-1", "user-1", "Physics", createdAt, updatedAt)

	if subject.ID != "sub-1" {
		t.Errorf("expected ID 'sub-1', got '%s'", subject.ID)
	}
	if subject.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", subject.UserID)
	}
	if subject.Name != "Physics" {
		t.Errorf("expected Name 'Physics', got '%s'", subject.Name)
	}
	if !subject.CreatedAt.Equal(createdAt) {
		t.Errorf("expected CreatedAt %v, got %v", createdAt, subject.CreatedAt)
	}
	if !subject.UpdatedAt.Equal(updatedAt) {
		t.Errorf("expected UpdatedAt %v, got %v", updatedAt, subject.UpdatedAt)
	}
}

func TestSubject_UpdateName_Valid(t *testing.T) {
	subject, err := domain.NewSubject("sub-1", "user-1", "Math")
	if err != nil {
		t.Fatalf("unexpected error creating subject: %v", err)
	}

	oldUpdatedAt := subject.UpdatedAt

	// Small delay to ensure UpdatedAt changes
	err = subject.UpdateName("Mathematics")
	if err != nil {
		t.Fatalf("unexpected error updating name: %v", err)
	}
	if subject.Name != "Mathematics" {
		t.Errorf("expected Name 'Mathematics', got '%s'", subject.Name)
	}
	if subject.UpdatedAt.Before(oldUpdatedAt) {
		t.Error("expected UpdatedAt to be updated")
	}
}

func TestSubject_UpdateName_EmptyName(t *testing.T) {
	subject, err := domain.NewSubject("sub-1", "user-1", "Math")
	if err != nil {
		t.Fatalf("unexpected error creating subject: %v", err)
	}

	err = subject.UpdateName("")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
	// Name should remain unchanged
	if subject.Name != "Math" {
		t.Errorf("expected name to remain 'Math', got '%s'", subject.Name)
	}
}

func TestSubject_UpdateName_TooLong(t *testing.T) {
	subject, err := domain.NewSubject("sub-1", "user-1", "Math")
	if err != nil {
		t.Fatalf("unexpected error creating subject: %v", err)
	}

	longName := strings.Repeat("x", 201)
	err = subject.UpdateName(longName)
	if err == nil {
		t.Fatal("expected error for name too long")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
	// Name should remain unchanged
	if subject.Name != "Math" {
		t.Errorf("expected name to remain 'Math', got '%s'", subject.Name)
	}
}
