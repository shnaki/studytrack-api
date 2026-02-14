package domain_test

import (
	"testing"

	"github.com/shnaki/studytrack-api/internal/domain"
)

func TestNewUser_Valid(t *testing.T) {
	user, err := domain.NewUser("test-id", "Alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "test-id" {
		t.Errorf("expected ID 'test-id', got '%s'", user.ID)
	}
	if user.Name != "Alice" {
		t.Errorf("expected Name 'Alice', got '%s'", user.Name)
	}
	if user.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestNewUser_EmptyName(t *testing.T) {
	_, err := domain.NewUser("test-id", "")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestNewUser_NameTooLong(t *testing.T) {
	longName := make([]byte, 101)
	for i := range longName {
		longName[i] = 'a'
	}
	_, err := domain.NewUser("test-id", string(longName))
	if err == nil {
		t.Fatal("expected error for long name")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}
