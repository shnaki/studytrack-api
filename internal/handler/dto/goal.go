package dto

import (
	"time"

	"github.com/shnaki/studytrack-api/internal/domain"
)

type UpsertGoalRequest struct {
	TargetMinutesPerWeek int     `json:"targetMinutesPerWeek" minimum:"1" doc:"Target study minutes per week"`
	StartDate            string  `json:"startDate" doc:"Start date (YYYY-MM-DD)"`
	EndDate              *string `json:"endDate,omitempty" doc:"End date (YYYY-MM-DD), optional"`
}

type GoalResponse struct {
	ID                   string    `json:"id" doc:"Goal ID"`
	UserID               string    `json:"userId" doc:"User ID"`
	SubjectID            string    `json:"subjectId" doc:"Subject ID"`
	TargetMinutesPerWeek int       `json:"targetMinutesPerWeek" doc:"Target minutes per week"`
	StartDate            string    `json:"startDate" doc:"Start date"`
	EndDate              *string   `json:"endDate,omitempty" doc:"End date"`
	CreatedAt            time.Time `json:"createdAt" doc:"Creation timestamp"`
	UpdatedAt            time.Time `json:"updatedAt" doc:"Last update timestamp"`
}

func ToGoalResponse(g *domain.Goal) GoalResponse {
	resp := GoalResponse{
		ID:                   g.ID,
		UserID:               g.UserID,
		SubjectID:            g.SubjectID,
		TargetMinutesPerWeek: g.TargetMinutesPerWeek,
		StartDate:            g.StartDate.Format("2006-01-02"),
		CreatedAt:            g.CreatedAt,
		UpdatedAt:            g.UpdatedAt,
	}
	if g.EndDate != nil {
		s := g.EndDate.Format("2006-01-02")
		resp.EndDate = &s
	}
	return resp
}

func ToGoalResponseList(goals []*domain.Goal) []GoalResponse {
	result := make([]GoalResponse, len(goals))
	for i, g := range goals {
		result[i] = ToGoalResponse(g)
	}
	return result
}
