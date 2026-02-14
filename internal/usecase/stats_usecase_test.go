package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase"
)

func TestGetWeeklyStats(t *testing.T) {
	weekStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // Monday

	subjectRepo := newMockSubjectRepository()
	subjectRepo.subjects["s1"] = &domain.Subject{ID: "s1", UserID: "u1", Name: "Math"}
	subjectRepo.subjects["s2"] = &domain.Subject{ID: "s2", UserID: "u1", Name: "English"}

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

	// Find Math and English in the results (map iteration order is not guaranteed)
	var math, english *domain.SubjectWeeklyStats
	for i := range stats.Subjects {
		switch stats.Subjects[i].SubjectName {
		case "Math":
			math = &stats.Subjects[i]
		case "English":
			english = &stats.Subjects[i]
		}
	}

	if math == nil {
		t.Fatal("expected Math subject in stats")
	} else {
		if math.TotalMinutes != 150 {
			t.Errorf("expected 150 minutes for Math, got %d", math.TotalMinutes)
		}
		if math.TargetMinutesPerWeek != 200 {
			t.Errorf("expected target 200, got %d", math.TargetMinutesPerWeek)
		}
		if math.AchievementRate != 75.0 {
			t.Errorf("expected 75%% achievement, got %.1f%%", math.AchievementRate)
		}
	}

	if english == nil {
		t.Fatal("expected English subject in stats")
	} else {
		if english.TotalMinutes != 30 {
			t.Errorf("expected 30 minutes for English, got %d", english.TotalMinutes)
		}
		if english.TargetMinutesPerWeek != 0 {
			t.Errorf("expected target 0 (no goal), got %d", english.TargetMinutesPerWeek)
		}
	}
}
