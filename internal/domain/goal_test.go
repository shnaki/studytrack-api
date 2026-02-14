package domain_test

import (
	"testing"
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
)

func TestNewGoal_Valid(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)
	goal, err := domain.NewGoal("goal-1", "user-1", "subject-1", 300, start, &end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if goal.TargetMinutesPerWeek != 300 {
		t.Errorf("expected 300, got %d", goal.TargetMinutesPerWeek)
	}
}

func TestNewGoal_NoEndDate(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	goal, err := domain.NewGoal("goal-1", "user-1", "subject-1", 300, start, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if goal.EndDate != nil {
		t.Error("expected nil EndDate")
	}
}

func TestNewGoal_EndBeforeStart(t *testing.T) {
	start := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := domain.NewGoal("goal-1", "user-1", "subject-1", 300, start, &end)
	if err == nil {
		t.Fatal("expected error for end before start")
	}
}

func TestNewGoal_ZeroTarget(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := domain.NewGoal("goal-1", "user-1", "subject-1", 0, start, nil)
	if err == nil {
		t.Fatal("expected error for zero target")
	}
}
