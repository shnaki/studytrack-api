package dto_test

import (
	"testing"
	"time"

	"github.com/shnaki/studytrack-api/internal/controller/dto"
	"github.com/shnaki/studytrack-api/internal/domain"
)

func TestToUserResponse(t *testing.T) {
	now := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	user := &domain.User{
		ID:        "user-123",
		Name:      "Alice",
		CreatedAt: now,
		UpdatedAt: now.Add(1 * time.Hour),
	}

	resp := dto.ToUserResponse(user)

	if resp.ID != "user-123" {
		t.Errorf("expected ID 'user-123', got '%s'", resp.ID)
	}
	if resp.Name != "Alice" {
		t.Errorf("expected Name 'Alice', got '%s'", resp.Name)
	}
	if !resp.CreatedAt.Equal(now) {
		t.Errorf("expected CreatedAt %v, got %v", now, resp.CreatedAt)
	}
	if !resp.UpdatedAt.Equal(now.Add(1 * time.Hour)) {
		t.Errorf("expected UpdatedAt %v, got %v", now.Add(1*time.Hour), resp.UpdatedAt)
	}
}

func TestToProjectResponse(t *testing.T) {
	now := time.Date(2024, 3, 10, 8, 0, 0, 0, time.UTC)
	project := &domain.Project{
		ID:        "proj-1",
		UserID:    "user-1",
		Name:      "Mathematics",
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := dto.ToProjectResponse(project)

	if resp.ID != "proj-1" {
		t.Errorf("expected ID 'proj-1', got '%s'", resp.ID)
	}
	if resp.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", resp.UserID)
	}
	if resp.Name != "Mathematics" {
		t.Errorf("expected Name 'Mathematics', got '%s'", resp.Name)
	}
	if !resp.CreatedAt.Equal(now) {
		t.Errorf("expected CreatedAt %v, got %v", now, resp.CreatedAt)
	}
	if !resp.UpdatedAt.Equal(now) {
		t.Errorf("expected UpdatedAt %v, got %v", now, resp.UpdatedAt)
	}
}

func TestToProjectResponseList(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	projects := []*domain.Project{
		{ID: "p1", UserID: "u1", Name: "Math", CreatedAt: now, UpdatedAt: now},
		{ID: "p2", UserID: "u1", Name: "English", CreatedAt: now, UpdatedAt: now},
		{ID: "p3", UserID: "u1", Name: "Science", CreatedAt: now, UpdatedAt: now},
	}

	result := dto.ToProjectResponseList(projects)

	if len(result) != 3 {
		t.Fatalf("expected 3 responses, got %d", len(result))
	}
	if result[0].Name != "Math" {
		t.Errorf("expected first project 'Math', got '%s'", result[0].Name)
	}
	if result[1].Name != "English" {
		t.Errorf("expected second project 'English', got '%s'", result[1].Name)
	}
	if result[2].Name != "Science" {
		t.Errorf("expected third project 'Science', got '%s'", result[2].Name)
	}
}

func TestToProjectResponseList_Empty(t *testing.T) {
	result := dto.ToProjectResponseList([]*domain.Project{})

	if len(result) != 0 {
		t.Errorf("expected empty list, got %d items", len(result))
	}
}

func TestToStudyLogResponse(t *testing.T) {
	studiedAt := time.Date(2024, 5, 20, 14, 0, 0, 0, time.UTC)
	createdAt := time.Date(2024, 5, 20, 14, 30, 0, 0, time.UTC)
	log := &domain.StudyLog{
		ID:        "log-1",
		UserID:    "user-1",
		ProjectID: "proj-1",
		StudiedAt: studiedAt,
		Minutes:   90,
		Note:      "Chapter 5 exercises",
		CreatedAt: createdAt,
	}

	resp := dto.ToStudyLogResponse(log)

	if resp.ID != "log-1" {
		t.Errorf("expected ID 'log-1', got '%s'", resp.ID)
	}
	if resp.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", resp.UserID)
	}
	if resp.ProjectID != "proj-1" {
		t.Errorf("expected ProjectID 'proj-1', got '%s'", resp.ProjectID)
	}
	if !resp.StudiedAt.Equal(studiedAt) {
		t.Errorf("expected StudiedAt %v, got %v", studiedAt, resp.StudiedAt)
	}
	if resp.Minutes != 90 {
		t.Errorf("expected Minutes 90, got %d", resp.Minutes)
	}
	if resp.Note != "Chapter 5 exercises" {
		t.Errorf("expected Note 'Chapter 5 exercises', got '%s'", resp.Note)
	}
	if !resp.CreatedAt.Equal(createdAt) {
		t.Errorf("expected CreatedAt %v, got %v", createdAt, resp.CreatedAt)
	}
}

func TestToStudyLogResponseList(t *testing.T) {
	now := time.Date(2024, 1, 10, 12, 0, 0, 0, time.UTC)
	logs := []*domain.StudyLog{
		{ID: "l1", UserID: "u1", ProjectID: "p1", StudiedAt: now, Minutes: 60, Note: "note1", CreatedAt: now},
		{ID: "l2", UserID: "u1", ProjectID: "p2", StudiedAt: now, Minutes: 30, Note: "note2", CreatedAt: now},
	}

	result := dto.ToStudyLogResponseList(logs)

	if len(result) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(result))
	}
	if result[0].ID != "l1" {
		t.Errorf("expected first ID 'l1', got '%s'", result[0].ID)
	}
	if result[0].Minutes != 60 {
		t.Errorf("expected first Minutes 60, got %d", result[0].Minutes)
	}
	if result[1].ID != "l2" {
		t.Errorf("expected second ID 'l2', got '%s'", result[1].ID)
	}
	if result[1].Minutes != 30 {
		t.Errorf("expected second Minutes 30, got %d", result[1].Minutes)
	}
}

func TestToStudyLogResponseList_Empty(t *testing.T) {
	result := dto.ToStudyLogResponseList([]*domain.StudyLog{})

	if len(result) != 0 {
		t.Errorf("expected empty list, got %d items", len(result))
	}
}

func TestToGoalResponse_WithoutEndDate(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	createdAt := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC)
	goal := &domain.Goal{
		ID:                   "goal-1",
		UserID:               "user-1",
		ProjectID:            "proj-1",
		TargetMinutesPerWeek: 300,
		StartDate:            startDate,
		EndDate:              nil,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
	}

	resp := dto.ToGoalResponse(goal)

	if resp.ID != "goal-1" {
		t.Errorf("expected ID 'goal-1', got '%s'", resp.ID)
	}
	if resp.UserID != "user-1" {
		t.Errorf("expected UserID 'user-1', got '%s'", resp.UserID)
	}
	if resp.ProjectID != "proj-1" {
		t.Errorf("expected ProjectID 'proj-1', got '%s'", resp.ProjectID)
	}
	if resp.TargetMinutesPerWeek != 300 {
		t.Errorf("expected TargetMinutesPerWeek 300, got %d", resp.TargetMinutesPerWeek)
	}
	if resp.StartDate != "2024-01-01" {
		t.Errorf("expected StartDate '2024-01-01', got '%s'", resp.StartDate)
	}
	if resp.EndDate != nil {
		t.Errorf("expected EndDate nil, got '%v'", resp.EndDate)
	}
	if !resp.CreatedAt.Equal(createdAt) {
		t.Errorf("expected CreatedAt %v, got %v", createdAt, resp.CreatedAt)
	}
	if !resp.UpdatedAt.Equal(updatedAt) {
		t.Errorf("expected UpdatedAt %v, got %v", updatedAt, resp.UpdatedAt)
	}
}

func TestToGoalResponse_WithEndDate(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	now := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	goal := &domain.Goal{
		ID:                   "goal-2",
		UserID:               "user-1",
		ProjectID:            "proj-1",
		TargetMinutesPerWeek: 120,
		StartDate:            startDate,
		EndDate:              &endDate,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	resp := dto.ToGoalResponse(goal)

	if resp.StartDate != "2024-01-01" {
		t.Errorf("expected StartDate '2024-01-01', got '%s'", resp.StartDate)
	}
	if resp.EndDate == nil {
		t.Fatal("expected EndDate to be non-nil")
	}
	if *resp.EndDate != "2024-06-30" {
		t.Errorf("expected EndDate '2024-06-30', got '%s'", *resp.EndDate)
	}
}

func TestToGoalResponseList(t *testing.T) {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	goals := []*domain.Goal{
		{ID: "g1", UserID: "u1", ProjectID: "p1", TargetMinutesPerWeek: 100, StartDate: now, CreatedAt: now, UpdatedAt: now},
		{ID: "g2", UserID: "u1", ProjectID: "p2", TargetMinutesPerWeek: 200, StartDate: now, CreatedAt: now, UpdatedAt: now},
	}

	result := dto.ToGoalResponseList(goals)

	if len(result) != 2 {
		t.Fatalf("expected 2 responses, got %d", len(result))
	}
	if result[0].ID != "g1" {
		t.Errorf("expected first ID 'g1', got '%s'", result[0].ID)
	}
	if result[0].TargetMinutesPerWeek != 100 {
		t.Errorf("expected first target 100, got %d", result[0].TargetMinutesPerWeek)
	}
	if result[1].ID != "g2" {
		t.Errorf("expected second ID 'g2', got '%s'", result[1].ID)
	}
	if result[1].TargetMinutesPerWeek != 200 {
		t.Errorf("expected second target 200, got %d", result[1].TargetMinutesPerWeek)
	}
}

func TestToGoalResponseList_Empty(t *testing.T) {
	result := dto.ToGoalResponseList([]*domain.Goal{})

	if len(result) != 0 {
		t.Errorf("expected empty list, got %d items", len(result))
	}
}

func TestToWeeklyStatsResponse(t *testing.T) {
	weekStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	stats := &domain.WeeklyStats{
		WeekStart: weekStart,
		Projects: []domain.ProjectWeeklyStats{
			{
				ProjectID:            "p1",
				ProjectName:          "Math",
				TotalMinutes:         150,
				TargetMinutesPerWeek: 200,
				AchievementRate:      75.0,
			},
			{
				ProjectID:            "p2",
				ProjectName:          "English",
				TotalMinutes:         30,
				TargetMinutesPerWeek: 0,
				AchievementRate:      0,
			},
		},
		TotalMinutes: 180,
	}

	resp := dto.ToWeeklyStatsResponse(stats)

	if resp.WeekStart != "2024-01-01" {
		t.Errorf("expected WeekStart '2024-01-01', got '%s'", resp.WeekStart)
	}
	if resp.TotalMinutes != 180 {
		t.Errorf("expected TotalMinutes 180, got %d", resp.TotalMinutes)
	}
	if len(resp.Projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(resp.Projects))
	}

	math := resp.Projects[0]
	if math.ProjectID != "p1" {
		t.Errorf("expected ProjectID 'p1', got '%s'", math.ProjectID)
	}
	if math.ProjectName != "Math" {
		t.Errorf("expected ProjectName 'Math', got '%s'", math.ProjectName)
	}
	if math.TotalMinutes != 150 {
		t.Errorf("expected TotalMinutes 150, got %d", math.TotalMinutes)
	}
	if math.TargetMinutesPerWeek != 200 {
		t.Errorf("expected TargetMinutesPerWeek 200, got %d", math.TargetMinutesPerWeek)
	}
	if math.AchievementRate != 75.0 {
		t.Errorf("expected AchievementRate 75.0, got %f", math.AchievementRate)
	}

	english := resp.Projects[1]
	if english.ProjectID != "p2" {
		t.Errorf("expected ProjectID 'p2', got '%s'", english.ProjectID)
	}
	if english.ProjectName != "English" {
		t.Errorf("expected ProjectName 'English', got '%s'", english.ProjectName)
	}
	if english.TotalMinutes != 30 {
		t.Errorf("expected TotalMinutes 30, got %d", english.TotalMinutes)
	}
	if english.TargetMinutesPerWeek != 0 {
		t.Errorf("expected TargetMinutesPerWeek 0, got %d", english.TargetMinutesPerWeek)
	}
	if english.AchievementRate != 0 {
		t.Errorf("expected AchievementRate 0, got %f", english.AchievementRate)
	}
}

func TestToWeeklyStatsResponse_NoProjects(t *testing.T) {
	weekStart := time.Date(2024, 2, 5, 0, 0, 0, 0, time.UTC)
	stats := &domain.WeeklyStats{
		WeekStart:    weekStart,
		Projects:     []domain.ProjectWeeklyStats{},
		TotalMinutes: 0,
	}

	resp := dto.ToWeeklyStatsResponse(stats)

	if resp.WeekStart != "2024-02-05" {
		t.Errorf("expected WeekStart '2024-02-05', got '%s'", resp.WeekStart)
	}
	if resp.TotalMinutes != 0 {
		t.Errorf("expected TotalMinutes 0, got %d", resp.TotalMinutes)
	}
	if len(resp.Projects) != 0 {
		t.Errorf("expected 0 projects, got %d", len(resp.Projects))
	}
}
