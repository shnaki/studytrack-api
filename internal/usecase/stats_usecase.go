package usecase

import (
	"context"
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase/port"
)

// StatsUsecase provides methods for calculating study statistics.
type StatsUsecase struct {
	studyLogRepo port.StudyLogRepository
	goalRepo     port.GoalRepository
	subjectRepo  port.SubjectRepository
}

// NewStatsUsecase creates a new StatsUsecase.
func NewStatsUsecase(
	studyLogRepo port.StudyLogRepository,
	goalRepo port.GoalRepository,
	subjectRepo port.SubjectRepository,
) *StatsUsecase {
	return &StatsUsecase{
		studyLogRepo: studyLogRepo,
		goalRepo:     goalRepo,
		subjectRepo:  subjectRepo,
	}
}

// GetWeeklyStats calculates study statistics for a specific week.
func (u *StatsUsecase) GetWeeklyStats(ctx context.Context, userID string, weekStart time.Time) (*domain.WeeklyStats, error) {
	weekEnd := weekStart.AddDate(0, 0, 7)

	subjects, err := u.subjectRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	filter := port.StudyLogFilter{
		From: &weekStart,
		To:   &weekEnd,
	}
	logs, err := u.studyLogRepo.FindByUserID(ctx, userID, filter)
	if err != nil {
		return nil, err
	}

	goals, err := u.goalRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	minutesBySubject := make(map[string]int)
	for _, log := range logs {
		minutesBySubject[log.SubjectID] += log.Minutes
	}

	goalBySubject := make(map[string]*domain.Goal)
	for _, goal := range goals {
		goalBySubject[goal.SubjectID] = goal
	}

	stats := &domain.WeeklyStats{
		WeekStart: weekStart,
	}
	var totalMinutes int

	for _, subject := range subjects {
		minutes := minutesBySubject[subject.ID]
		totalMinutes += minutes

		ss := domain.SubjectWeeklyStats{
			SubjectID:    subject.ID,
			SubjectName:  subject.Name,
			TotalMinutes: minutes,
		}

		if goal, ok := goalBySubject[subject.ID]; ok {
			ss.TargetMinutesPerWeek = goal.TargetMinutesPerWeek
			if goal.TargetMinutesPerWeek > 0 {
				ss.AchievementRate = float64(minutes) / float64(goal.TargetMinutesPerWeek) * 100
			}
		}

		stats.Subjects = append(stats.Subjects, ss)
	}

	stats.TotalMinutes = totalMinutes
	return stats, nil
}
