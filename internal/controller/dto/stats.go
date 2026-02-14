package dto

import "github.com/shnaki/studytrack-api/internal/domain"

type SubjectWeeklyStatsResponse struct {
	SubjectID            string  `json:"subjectId" doc:"Subject ID"`
	SubjectName          string  `json:"subjectName" doc:"Subject name"`
	TotalMinutes         int     `json:"totalMinutes" doc:"Total minutes studied this week"`
	TargetMinutesPerWeek int     `json:"targetMinutesPerWeek" doc:"Weekly goal target (0 if no goal)"`
	AchievementRate      float64 `json:"achievementRate" doc:"Achievement rate percentage (0 if no goal)"`
}

type WeeklyStatsResponse struct {
	WeekStart    string                       `json:"weekStart" doc:"Week start date"`
	Subjects     []SubjectWeeklyStatsResponse `json:"subjects" doc:"Per-subject stats"`
	TotalMinutes int                          `json:"totalMinutes" doc:"Total minutes across all subjects"`
}

func ToWeeklyStatsResponse(s *domain.WeeklyStats) WeeklyStatsResponse {
	subjects := make([]SubjectWeeklyStatsResponse, len(s.Subjects))
	for i, sub := range s.Subjects {
		subjects[i] = SubjectWeeklyStatsResponse{
			SubjectID:            sub.SubjectID,
			SubjectName:          sub.SubjectName,
			TotalMinutes:         sub.TotalMinutes,
			TargetMinutesPerWeek: sub.TargetMinutesPerWeek,
			AchievementRate:      sub.AchievementRate,
		}
	}
	return WeeklyStatsResponse{
		WeekStart:    s.WeekStart.Format("2006-01-02"),
		Subjects:     subjects,
		TotalMinutes: s.TotalMinutes,
	}
}
