package domain

import "time"

type WeeklyStats struct {
	WeekStart    time.Time
	Subjects     []SubjectWeeklyStats
	TotalMinutes int
}

type SubjectWeeklyStats struct {
	SubjectID            string
	SubjectName          string
	TotalMinutes         int
	TargetMinutesPerWeek int
	AchievementRate      float64
}
