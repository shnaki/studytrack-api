package domain_test

import (
	"testing"
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
)

func TestNewGoal_Valid(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)
	goal, err := domain.NewGoal("goal-1", "user-1", "project-1", 300, start, &end)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if goal.TargetMinutesPerWeek != 300 {
		t.Errorf("expected 300, got %d", goal.TargetMinutesPerWeek)
	}
}

func TestNewGoal_NoEndDate(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	goal, err := domain.NewGoal("goal-1", "user-1", "project-1", 300, start, nil)
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
	_, err := domain.NewGoal("goal-1", "user-1", "project-1", 300, start, &end)
	if err == nil {
		t.Fatal("expected error for end before start")
	}
}

func TestNewGoal_ZeroTarget(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := domain.NewGoal("goal-1", "user-1", "project-1", 0, start, nil)
	if err == nil {
		t.Fatal("expected error for zero target")
	}
}

func TestReconstructGoal(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	createdAt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 3, 15, 14, 0, 0, 0, time.UTC)

	goal := domain.ReconstructGoal("goal-1", "user-1", "project-1", 300, startDate, &endDate, createdAt, updatedAt)

	if goal.ID != "goal-1" {
		t.Errorf("expected ID 'goal-1', got '%s'", goal.ID)
	}
	if goal.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", goal.UserID)
	}
	if goal.ProjectID != "project-1" {
		t.Errorf("expected ProjectID 'project-1', got '%s'", goal.ProjectID)
	}
	if goal.TargetMinutesPerWeek != 300 {
		t.Errorf("expected TargetMinutesPerWeek 300, got %d", goal.TargetMinutesPerWeek)
	}
	if !goal.StartDate.Equal(startDate) {
		t.Errorf("expected StartDate %v, got %v", startDate, goal.StartDate)
	}
	if goal.EndDate == nil {
		t.Fatal("expected EndDate to be set")
	}
	if !goal.EndDate.Equal(endDate) {
		t.Errorf("expected EndDate %v, got %v", endDate, *goal.EndDate)
	}
	if !goal.CreatedAt.Equal(createdAt) {
		t.Errorf("expected CreatedAt %v, got %v", createdAt, goal.CreatedAt)
	}
	if !goal.UpdatedAt.Equal(updatedAt) {
		t.Errorf("expected UpdatedAt %v, got %v", updatedAt, goal.UpdatedAt)
	}
}

func TestReconstructGoal_NilEndDate(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	createdAt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)

	goal := domain.ReconstructGoal("goal-1", "user-1", "project-1", 300, startDate, nil, createdAt, updatedAt)

	if goal.EndDate != nil {
		t.Error("expected EndDate to be nil")
	}
}
