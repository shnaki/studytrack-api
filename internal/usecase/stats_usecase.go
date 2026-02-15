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
	projectRepo  port.ProjectRepository
}

// NewStatsUsecase creates a new StatsUsecase.
func NewStatsUsecase(
	studyLogRepo port.StudyLogRepository,
	goalRepo port.GoalRepository,
	projectRepo port.ProjectRepository,
) *StatsUsecase {
	return &StatsUsecase{
		studyLogRepo: studyLogRepo,
		goalRepo:     goalRepo,
		projectRepo:  projectRepo,
	}
}

// GetWeeklyStats calculates study statistics for a specific week.
func (u *StatsUsecase) GetWeeklyStats(ctx context.Context, userID string, weekStart time.Time) (*domain.WeeklyStats, error) {
	weekEnd := weekStart.AddDate(0, 0, 7)

	projects, err := u.projectRepo.FindByUserID(ctx, userID)
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

	minutesByProject := make(map[string]int)
	for _, log := range logs {
		minutesByProject[log.ProjectID] += log.Minutes
	}

	goalByProject := make(map[string]*domain.Goal)
	for _, goal := range goals {
		goalByProject[goal.ProjectID] = goal
	}

	stats := &domain.WeeklyStats{
		WeekStart: weekStart,
	}
	var totalMinutes int

	for _, project := range projects {
		minutes := minutesByProject[project.ID]
		totalMinutes += minutes

		ps := domain.ProjectWeeklyStats{
			ProjectID:    project.ID,
			ProjectName:  project.Name,
			TotalMinutes: minutes,
		}

		if goal, ok := goalByProject[project.ID]; ok {
			ps.TargetMinutesPerWeek = goal.TargetMinutesPerWeek
			if goal.TargetMinutesPerWeek > 0 {
				ps.AchievementRate = float64(minutes) / float64(goal.TargetMinutesPerWeek) * 100
			}
		}

		stats.Projects = append(stats.Projects, ps)
	}

	stats.TotalMinutes = totalMinutes
	return stats, nil
}
