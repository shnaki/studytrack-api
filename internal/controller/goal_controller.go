package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/shnaki/studytrack-api/internal/controller/dto"
	"github.com/shnaki/studytrack-api/internal/usecase"
)

type upsertGoalInput struct {
	UserID    string `path:"userId" doc:"User ID"`
	SubjectID string `path:"subjectId" doc:"Subject ID"`
	Body      dto.UpsertGoalRequest
}

type upsertGoalOutput struct {
	Body dto.GoalResponse
}

type listGoalsInput struct {
	UserID string `path:"userId" doc:"User ID"`
}

type listGoalsOutput struct {
	Body []dto.GoalResponse
}

// RegisterGoalRoutes registers goal-related routes to the Huma API.
func RegisterGoalRoutes(api huma.API, uc *usecase.GoalUsecase) {
	huma.Register(api, huma.Operation{
		OperationID: "upsert-goal",
		Method:      http.MethodPut,
		Path:        "/users/{userId}/goals/{subjectId}",
		Summary:     "Create or update a goal for a subject",
		Tags:        []string{"Goals"},
	}, func(ctx context.Context, input *upsertGoalInput) (*upsertGoalOutput, error) {
		startDate, err := time.Parse("2006-01-02", input.Body.StartDate)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid startDate format, expected YYYY-MM-DD")
		}

		var endDate *time.Time
		if input.Body.EndDate != nil {
			t, err := time.Parse("2006-01-02", *input.Body.EndDate)
			if err != nil {
				return nil, huma.Error400BadRequest("invalid endDate format, expected YYYY-MM-DD")
			}
			endDate = &t
		}

		goal, err := uc.UpsertGoal(ctx, input.UserID, input.SubjectID, input.Body.TargetMinutesPerWeek, startDate, endDate)
		if err != nil {
			return nil, toHTTPError(err)
		}
		return &upsertGoalOutput{Body: dto.ToGoalResponse(goal)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-goals",
		Method:      http.MethodGet,
		Path:        "/users/{userId}/goals",
		Summary:     "List goals for a user",
		Tags:        []string{"Goals"},
	}, func(ctx context.Context, input *listGoalsInput) (*listGoalsOutput, error) {
		goals, err := uc.ListGoals(ctx, input.UserID)
		if err != nil {
			return nil, toHTTPError(err)
		}
		return &listGoalsOutput{Body: dto.ToGoalResponseList(goals)}, nil
	})
}
