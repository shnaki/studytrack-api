package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/shnaki/studytrack-api/internal/controller/dto"
	"github.com/shnaki/studytrack-api/internal/usecase"
)

type getWeeklyStatsInput struct {
	UserID    string `path:"userId" doc:"User ID"`
	WeekStart string `query:"weekStart" required:"true" doc:"Week start date (YYYY-MM-DD)" example:"2024-01-01"`
}

type getWeeklyStatsOutput struct {
	Body dto.WeeklyStatsResponse
}

// RegisterStatsRoutes registers statistics-related routes to the Huma API.
func RegisterStatsRoutes(api huma.API, uc *usecase.StatsUsecase) {
	huma.Register(api, huma.Operation{
		OperationID: "get-weekly-stats",
		Method:      http.MethodGet,
		Path:        "/users/{userId}/stats/weekly",
		Summary:     "Get weekly study statistics",
		Tags:        []string{"Stats"},
	}, func(ctx context.Context, input *getWeeklyStatsInput) (*getWeeklyStatsOutput, error) {
		weekStart, err := time.Parse("2006-01-02", input.WeekStart)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid weekStart format, expected YYYY-MM-DD")
		}

		stats, err := uc.GetWeeklyStats(ctx, input.UserID, weekStart)
		if err != nil {
			return nil, toHTTPError(err)
		}
		return &getWeeklyStatsOutput{Body: dto.ToWeeklyStatsResponse(stats)}, nil
	})
}
