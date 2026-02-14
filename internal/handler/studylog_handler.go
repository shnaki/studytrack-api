package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/shnaki/studytrack-api/internal/application/port"
	"github.com/shnaki/studytrack-api/internal/application/usecase"
	"github.com/shnaki/studytrack-api/internal/handler/dto"
)

type createStudyLogInput struct {
	UserID string `path:"userId" doc:"User ID"`
	Body   dto.CreateStudyLogRequest
}

type createStudyLogOutput struct {
	Body dto.StudyLogResponse
}

type listStudyLogsInput struct {
	UserID    string  `path:"userId" doc:"User ID"`
	From      *string `query:"from" doc:"Start date (YYYY-MM-DD)" example:"2024-01-01"`
	To        *string `query:"to" doc:"End date (YYYY-MM-DD)" example:"2024-01-31"`
	SubjectID *string `query:"subjectId" doc:"Filter by subject ID"`
}

type listStudyLogsOutput struct {
	Body []dto.StudyLogResponse
}

type deleteStudyLogInput struct {
	ID string `path:"id" doc:"Study log ID"`
}

func RegisterStudyLogRoutes(api huma.API, uc *usecase.StudyLogUsecase) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-study-log",
		Method:        http.MethodPost,
		Path:          "/users/{userId}/study-logs",
		Summary:       "Create a study log",
		Tags:          []string{"StudyLogs"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *createStudyLogInput) (*createStudyLogOutput, error) {
		log, err := uc.CreateStudyLog(ctx, input.UserID, input.Body.SubjectID, input.Body.StudiedAt, input.Body.Minutes, input.Body.Note)
		if err != nil {
			return nil, toHTTPError(err)
		}
		return &createStudyLogOutput{Body: dto.ToStudyLogResponse(log)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-study-logs",
		Method:      http.MethodGet,
		Path:        "/users/{userId}/study-logs",
		Summary:     "List study logs for a user",
		Tags:        []string{"StudyLogs"},
	}, func(ctx context.Context, input *listStudyLogsInput) (*listStudyLogsOutput, error) {
		filter, err := parseStudyLogFilter(input)
		if err != nil {
			return nil, err
		}
		logs, err := uc.ListStudyLogs(ctx, input.UserID, filter)
		if err != nil {
			return nil, toHTTPError(err)
		}
		return &listStudyLogsOutput{Body: dto.ToStudyLogResponseList(logs)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-study-log",
		Method:        http.MethodDelete,
		Path:          "/study-logs/{id}",
		Summary:       "Delete a study log",
		Tags:          []string{"StudyLogs"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *deleteStudyLogInput) (*struct{}, error) {
		if err := uc.DeleteStudyLog(ctx, input.ID); err != nil {
			return nil, toHTTPError(err)
		}
		return nil, nil
	})
}

func parseStudyLogFilter(input *listStudyLogsInput) (port.StudyLogFilter, error) {
	var filter port.StudyLogFilter

	if input.From != nil {
		t, err := time.Parse("2006-01-02", *input.From)
		if err != nil {
			return filter, huma.Error400BadRequest("invalid 'from' date format, expected YYYY-MM-DD")
		}
		filter.From = &t
	}
	if input.To != nil {
		t, err := time.Parse("2006-01-02", *input.To)
		if err != nil {
			return filter, huma.Error400BadRequest("invalid 'to' date format, expected YYYY-MM-DD")
		}
		endOfDay := t.AddDate(0, 0, 1)
		filter.To = &endOfDay
	}
	filter.SubjectID = input.SubjectID
	return filter, nil
}
