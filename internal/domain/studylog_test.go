package domain_test

import (
	"testing"
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
)

func TestNewStudyLog_Valid(t *testing.T) {
	now := time.Now()
	log, err := domain.NewStudyLog("log-1", "user-1", "subject-1", now, 60, "studied math")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if log.Minutes != 60 {
		t.Errorf("expected 60 minutes, got %d", log.Minutes)
	}
	if log.Note != "studied math" {
		t.Errorf("expected note 'studied math', got '%s'", log.Note)
	}
}

func TestNewStudyLog_ZeroMinutes(t *testing.T) {
	_, err := domain.NewStudyLog("log-1", "user-1", "subject-1", time.Now(), 0, "")
	if err == nil {
		t.Fatal("expected error for zero minutes")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestNewStudyLog_NegativeMinutes(t *testing.T) {
	_, err := domain.NewStudyLog("log-1", "user-1", "subject-1", time.Now(), -10, "")
	if err == nil {
		t.Fatal("expected error for negative minutes")
	}
}

func TestNewStudyLog_TooManyMinutes(t *testing.T) {
	_, err := domain.NewStudyLog("log-1", "user-1", "subject-1", time.Now(), 1441, "")
	if err == nil {
		t.Fatal("expected error for > 1440 minutes")
	}
}

func TestNewStudyLog_EmptyUserID(t *testing.T) {
	_, err := domain.NewStudyLog("log-1", "", "subject-1", time.Now(), 60, "")
	if err == nil {
		t.Fatal("expected error for empty user ID")
	}
}
