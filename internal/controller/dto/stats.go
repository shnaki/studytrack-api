package dto

import "github.com/shnaki/studytrack-api/internal/domain"

// ProjectWeeklyStatsResponse represents weekly statistics for a specific project.
type ProjectWeeklyStatsResponse struct {
	ProjectID            string  `json:"projectId" doc:"Project ID"`
	ProjectName          string  `json:"projectName" doc:"Project name"`
	TotalMinutes         int     `json:"totalMinutes" doc:"Total minutes studied this week"`
	TargetMinutesPerWeek int     `json:"targetMinutesPerWeek" doc:"Weekly goal target (0 if no goal)"`
	AchievementRate      float64 `json:"achievementRate" doc:"Achievement rate percentage (0 if no goal)"`
}

// WeeklyStatsResponse represents weekly statistics for all projects.
type WeeklyStatsResponse struct {
	WeekStart    string                       `json:"weekStart" doc:"Week start date"`
	Projects     []ProjectWeeklyStatsResponse `json:"projects" doc:"Per-project stats"`
	TotalMinutes int                          `json:"totalMinutes" doc:"Total minutes across all projects"`
}

// ToWeeklyStatsResponse converts domain.WeeklyStats to WeeklyStatsResponse.
func ToWeeklyStatsResponse(s *domain.WeeklyStats) WeeklyStatsResponse {
	projects := make([]ProjectWeeklyStatsResponse, len(s.Projects))
	for i, proj := range s.Projects {
		projects[i] = ProjectWeeklyStatsResponse{
			ProjectID:            proj.ProjectID,
			ProjectName:          proj.ProjectName,
			TotalMinutes:         proj.TotalMinutes,
			TargetMinutesPerWeek: proj.TargetMinutesPerWeek,
			AchievementRate:      proj.AchievementRate,
		}
	}
	return WeeklyStatsResponse{
		WeekStart:    s.WeekStart.Format("2006-01-02"),
		Projects:     projects,
		TotalMinutes: s.TotalMinutes,
	}
}
