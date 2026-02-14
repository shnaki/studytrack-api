package domain

import "time"

// WeeklyStats represents study statistics for a specific week.
type WeeklyStats struct {
	WeekStart    time.Time
	Subjects     []SubjectWeeklyStats
	TotalMinutes int
}

// SubjectWeeklyStats represents study statistics for a specific subject in a week.
type SubjectWeeklyStats struct {
	SubjectID            string
	SubjectName          string
	TotalMinutes         int
	TargetMinutesPerWeek int
	AchievementRate      float64
}
