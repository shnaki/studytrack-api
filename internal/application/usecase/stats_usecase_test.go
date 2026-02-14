package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/shnaki/studytrack-api/internal/application/port"
	"github.com/shnaki/studytrack-api/internal/application/usecase"
	"github.com/shnaki/studytrack-api/internal/domain"
)

// --- Mock Repositories for Stats ---

type mockSubjectRepository struct {
	subjects []*domain.Subject
}

func (m *mockSubjectRepository) Create(_ context.Context, s *domain.Subject) error {
	m.subjects = append(m.subjects, s)
	return nil
}
func (m *mockSubjectRepository) FindByID(_ context.Context, id string) (*domain.Subject, error) {
	for _, s := range m.subjects {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, domain.ErrNotFound("subject")
}
func (m *mockSubjectRepository) FindByUserID(_ context.Context, userID string) ([]*domain.Subject, error) {
	var result []*domain.Subject
	for _, s := range m.subjects {
		if s.UserID == userID {
			result = append(result, s)
		}
	}
	return result, nil
}
func (m *mockSubjectRepository) Update(_ context.Context, _ *domain.Subject) error { return nil }
func (m *mockSubjectRepository) Delete(_ context.Context, _ string) error          { return nil }

type mockStudyLogRepository struct {
	logs []*domain.StudyLog
}

func (m *mockStudyLogRepository) Create(_ context.Context, l *domain.StudyLog) error {
	m.logs = append(m.logs, l)
	return nil
}
func (m *mockStudyLogRepository) FindByID(_ context.Context, id string) (*domain.StudyLog, error) {
	for _, l := range m.logs {
		if l.ID == id {
			return l, nil
		}
	}
	return nil, domain.ErrNotFound("study log")
}
func (m *mockStudyLogRepository) FindByUserID(_ context.Context, userID string, filter port.StudyLogFilter) ([]*domain.StudyLog, error) {
	var result []*domain.StudyLog
	for _, l := range m.logs {
		if l.UserID != userID {
			continue
		}
		if filter.From != nil && l.StudiedAt.Before(*filter.From) {
			continue
		}
		if filter.To != nil && !l.StudiedAt.Before(*filter.To) {
			continue
		}
		if filter.SubjectID != nil && l.SubjectID != *filter.SubjectID {
			continue
		}
		result = append(result, l)
	}
	return result, nil
}
func (m *mockStudyLogRepository) Delete(_ context.Context, _ string) error { return nil }

type mockGoalRepository struct {
	goals []*domain.Goal
}

func (m *mockGoalRepository) Upsert(_ context.Context, g *domain.Goal) error {
	m.goals = append(m.goals, g)
	return nil
}
func (m *mockGoalRepository) FindByUserID(_ context.Context, userID string) ([]*domain.Goal, error) {
	var result []*domain.Goal
	for _, g := range m.goals {
		if g.UserID == userID {
			result = append(result, g)
		}
	}
	return result, nil
}

// --- Tests ---

func TestGetWeeklyStats(t *testing.T) {
	weekStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // Monday

	subjectRepo := &mockSubjectRepository{
		subjects: []*domain.Subject{
			{ID: "s1", UserID: "u1", Name: "Math"},
			{ID: "s2", UserID: "u1", Name: "English"},
		},
	}

	studyLogRepo := &mockStudyLogRepository{
		logs: []*domain.StudyLog{
			{ID: "l1", UserID: "u1", SubjectID: "s1", StudiedAt: weekStart.Add(1 * time.Hour), Minutes: 60},
			{ID: "l2", UserID: "u1", SubjectID: "s1", StudiedAt: weekStart.Add(25 * time.Hour), Minutes: 90},
			{ID: "l3", UserID: "u1", SubjectID: "s2", StudiedAt: weekStart.Add(2 * time.Hour), Minutes: 30},
		},
	}

	goalRepo := &mockGoalRepository{
		goals: []*domain.Goal{
			{ID: "g1", UserID: "u1", SubjectID: "s1", TargetMinutesPerWeek: 200},
		},
	}

	uc := usecase.NewStatsUsecase(studyLogRepo, goalRepo, subjectRepo)
	stats, err := uc.GetWeeklyStats(context.Background(), "u1", weekStart)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.TotalMinutes != 180 {
		t.Errorf("expected total 180, got %d", stats.TotalMinutes)
	}
	if len(stats.Subjects) != 2 {
		t.Fatalf("expected 2 subjects, got %d", len(stats.Subjects))
	}

	// Math: 150 min, target 200, rate 75%
	math := stats.Subjects[0]
	if math.SubjectName != "Math" {
		t.Errorf("expected 'Math', got '%s'", math.SubjectName)
	}
	if math.TotalMinutes != 150 {
		t.Errorf("expected 150 minutes for Math, got %d", math.TotalMinutes)
	}
	if math.TargetMinutesPerWeek != 200 {
		t.Errorf("expected target 200, got %d", math.TargetMinutesPerWeek)
	}
	if math.AchievementRate != 75.0 {
		t.Errorf("expected 75%% achievement, got %.1f%%", math.AchievementRate)
	}

	// English: 30 min, no goal
	english := stats.Subjects[1]
	if english.TotalMinutes != 30 {
		t.Errorf("expected 30 minutes for English, got %d", english.TotalMinutes)
	}
	if english.TargetMinutesPerWeek != 0 {
		t.Errorf("expected target 0 (no goal), got %d", english.TargetMinutesPerWeek)
	}
}
