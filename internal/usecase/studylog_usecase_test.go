package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase"
	"github.com/shnaki/studytrack-api/internal/usecase/port"
)

func setupStudyLogTest() (*usecase.StudyLogUsecase, *mockUserRepository, *mockProjectRepository, *mockStudyLogRepository) {
	userRepo := newMockUserRepository()
	projectRepo := newMockProjectRepository()
	studyLogRepo := newMockStudyLogRepository()
	uc := usecase.NewStudyLogUsecase(studyLogRepo, userRepo, projectRepo)
	return uc, userRepo, projectRepo, studyLogRepo
}

func TestCreateStudyLog_Success(t *testing.T) {
	uc, userRepo, projectRepo, _ := setupStudyLogTest()
	createTestUser(userRepo, "user-1", "Alice")
	projectRepo.projects["proj-1"] = &domain.Project{ID: "proj-1", UserID: "user-1", Name: "Math"}

	now := time.Now()
	log, err := uc.CreateStudyLog(context.Background(), "user-1", "proj-1", now, 60, "chapter 3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if log.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", log.UserID)
	}
	if log.ProjectID != "proj-1" {
		t.Errorf("expected ProjectID 'proj-1', got '%s'", log.ProjectID)
	}
	if log.Minutes != 60 {
		t.Errorf("expected 60 minutes, got %d", log.Minutes)
	}
	if log.Note != "chapter 3" {
		t.Errorf("expected note 'chapter 3', got '%s'", log.Note)
	}
	if log.ID == "" {
		t.Error("expected ID to be generated")
	}
}

func TestCreateStudyLog_UserNotFound(t *testing.T) {
	uc, _, projectRepo, _ := setupStudyLogTest()
	projectRepo.projects["proj-1"] = &domain.Project{ID: "proj-1", UserID: "user-1", Name: "Math"}

	_, err := uc.CreateStudyLog(context.Background(), "nonexistent", "proj-1", time.Now(), 60, "")
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestCreateStudyLog_ProjectNotFound(t *testing.T) {
	uc, userRepo, _, _ := setupStudyLogTest()
	createTestUser(userRepo, "user-1", "Alice")

	_, err := uc.CreateStudyLog(context.Background(), "user-1", "nonexistent", time.Now(), 60, "")
	if err == nil {
		t.Fatal("expected error for nonexistent project")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestCreateStudyLog_ProjectBelongsToDifferentUser(t *testing.T) {
	uc, userRepo, projectRepo, _ := setupStudyLogTest()
	createTestUser(userRepo, "user-1", "Alice")
	createTestUser(userRepo, "user-2", "Bob")
	projectRepo.projects["proj-1"] = &domain.Project{ID: "proj-1", UserID: "user-2", Name: "Math"}

	_, err := uc.CreateStudyLog(context.Background(), "user-1", "proj-1", time.Now(), 60, "")
	if err == nil {
		t.Fatal("expected error when project belongs to different user")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}

func TestCreateStudyLog_InvalidMinutes(t *testing.T) {
	uc, userRepo, projectRepo, _ := setupStudyLogTest()
	createTestUser(userRepo, "user-1", "Alice")
	projectRepo.projects["proj-1"] = &domain.Project{ID: "proj-1", UserID: "user-1", Name: "Math"}

	// Zero minutes
	_, err := uc.CreateStudyLog(context.Background(), "user-1", "proj-1", time.Now(), 0, "")
	if err == nil {
		t.Fatal("expected error for zero minutes")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error for zero minutes, got: %v", err)
	}

	// Negative minutes
	_, err = uc.CreateStudyLog(context.Background(), "user-1", "proj-1", time.Now(), -5, "")
	if err == nil {
		t.Fatal("expected error for negative minutes")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error for negative minutes, got: %v", err)
	}

	// Over 1440 minutes
	_, err = uc.CreateStudyLog(context.Background(), "user-1", "proj-1", time.Now(), 1441, "")
	if err == nil {
		t.Fatal("expected error for > 1440 minutes")
	}
	if !domain.IsValidation(err) {
		t.Errorf("expected validation error for > 1440 minutes, got: %v", err)
	}
}

func TestListStudyLogs_Success(t *testing.T) {
	uc, userRepo, projectRepo, _ := setupStudyLogTest()
	createTestUser(userRepo, "user-1", "Alice")
	projectRepo.projects["proj-1"] = &domain.Project{ID: "proj-1", UserID: "user-1", Name: "Math"}

	_, err := uc.CreateStudyLog(context.Background(), "user-1", "proj-1", time.Now(), 60, "session 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = uc.CreateStudyLog(context.Background(), "user-1", "proj-1", time.Now(), 30, "session 2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	logs, err := uc.ListStudyLogs(context.Background(), "user-1", port.StudyLogFilter{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(logs) != 2 {
		t.Errorf("expected 2 logs, got %d", len(logs))
	}
}

func TestListStudyLogs_WithFilter(t *testing.T) {
	uc, userRepo, projectRepo, _ := setupStudyLogTest()
	createTestUser(userRepo, "user-1", "Alice")
	projectRepo.projects["proj-1"] = &domain.Project{ID: "proj-1", UserID: "user-1", Name: "Math"}
	projectRepo.projects["proj-2"] = &domain.Project{ID: "proj-2", UserID: "user-1", Name: "English"}

	day1 := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	day2 := time.Date(2024, 1, 5, 10, 0, 0, 0, time.UTC)
	day3 := time.Date(2024, 1, 10, 10, 0, 0, 0, time.UTC)

	_, _ = uc.CreateStudyLog(context.Background(), "user-1", "proj-1", day1, 60, "day1 math")
	_, _ = uc.CreateStudyLog(context.Background(), "user-1", "proj-2", day2, 45, "day2 english")
	_, _ = uc.CreateStudyLog(context.Background(), "user-1", "proj-1", day3, 30, "day3 math")

	// Filter by date range: from day2 to day3 (exclusive)
	from := time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)
	logs, err := uc.ListStudyLogs(context.Background(), "user-1", port.StudyLogFilter{
		From: &from,
		To:   &to,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(logs) != 1 {
		t.Errorf("expected 1 log with date filter, got %d", len(logs))
	}

	// Filter by project
	projectID := "proj-1"
	logs, err = uc.ListStudyLogs(context.Background(), "user-1", port.StudyLogFilter{
		ProjectID: &projectID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(logs) != 2 {
		t.Errorf("expected 2 logs for proj-1, got %d", len(logs))
	}
}

func TestDeleteStudyLog_Success(t *testing.T) {
	uc, userRepo, projectRepo, studyLogRepo := setupStudyLogTest()
	createTestUser(userRepo, "user-1", "Alice")
	projectRepo.projects["proj-1"] = &domain.Project{ID: "proj-1", UserID: "user-1", Name: "Math"}

	created, err := uc.CreateStudyLog(context.Background(), "user-1", "proj-1", time.Now(), 60, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = uc.DeleteStudyLog(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify deleted
	if len(studyLogRepo.logs) != 0 {
		t.Errorf("expected 0 logs after deletion, got %d", len(studyLogRepo.logs))
	}
}

func TestDeleteStudyLog_NotFound(t *testing.T) {
	uc, _, _, _ := setupStudyLogTest()

	err := uc.DeleteStudyLog(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent study log")
	}
	if !domain.IsNotFound(err) {
		t.Errorf("expected not found error, got: %v", err)
	}
}
