package controller

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/shnaki/studytrack-api/internal/controller/dto"
	"github.com/shnaki/studytrack-api/internal/usecase"
)

type createSubjectInput struct {
	UserID string `path:"userId" doc:"User ID"`
	Body   dto.CreateSubjectRequest
}

type createSubjectOutput struct {
	Body dto.SubjectResponse
}

type listSubjectsInput struct {
	UserID string `path:"userId" doc:"User ID"`
}

type listSubjectsOutput struct {
	Body []dto.SubjectResponse
}

type updateSubjectInput struct {
	ID   string `path:"id" doc:"Subject ID"`
	Body dto.UpdateSubjectRequest
}

type updateSubjectOutput struct {
	Body dto.SubjectResponse
}

type deleteSubjectInput struct {
	ID string `path:"id" doc:"Subject ID"`
}

func RegisterSubjectRoutes(api huma.API, uc *usecase.SubjectUsecase) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-subject",
		Method:        http.MethodPost,
		Path:          "/users/{userId}/subjects",
		Summary:       "Create a new subject",
		Tags:          []string{"Subjects"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *createSubjectInput) (*createSubjectOutput, error) {
		subject, err := uc.CreateSubject(ctx, input.UserID, input.Body.Name)
		if err != nil {
			return nil, toHTTPError(err)
		}
		return &createSubjectOutput{Body: dto.ToSubjectResponse(subject)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-subjects",
		Method:      http.MethodGet,
		Path:        "/users/{userId}/subjects",
		Summary:     "List subjects for a user",
		Tags:        []string{"Subjects"},
	}, func(ctx context.Context, input *listSubjectsInput) (*listSubjectsOutput, error) {
		subjects, err := uc.ListSubjects(ctx, input.UserID)
		if err != nil {
			return nil, toHTTPError(err)
		}
		return &listSubjectsOutput{Body: dto.ToSubjectResponseList(subjects)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-subject",
		Method:      http.MethodPut,
		Path:        "/subjects/{id}",
		Summary:     "Update a subject",
		Tags:        []string{"Subjects"},
	}, func(ctx context.Context, input *updateSubjectInput) (*updateSubjectOutput, error) {
		subject, err := uc.UpdateSubject(ctx, input.ID, input.Body.Name)
		if err != nil {
			return nil, toHTTPError(err)
		}
		return &updateSubjectOutput{Body: dto.ToSubjectResponse(subject)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-subject",
		Method:        http.MethodDelete,
		Path:          "/subjects/{id}",
		Summary:       "Delete a subject",
		Tags:          []string{"Subjects"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *deleteSubjectInput) (*struct{}, error) {
		if err := uc.DeleteSubject(ctx, input.ID); err != nil {
			return nil, toHTTPError(err)
		}
		return nil, nil
	})
}
