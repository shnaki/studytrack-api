package domain

import "time"

// WeeklyStats represents study statistics for a specific week.
type WeeklyStats struct {
	WeekStart    time.Time
	Projects     []ProjectWeeklyStats
	TotalMinutes int
}

// ProjectWeeklyStats represents study statistics for a specific project in a week.
type ProjectWeeklyStats struct {
	ProjectID            string
	ProjectName          string
	TotalMinutes         int
	TargetMinutesPerWeek int
	AchievementRate      float64
}
