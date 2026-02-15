package controller

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/shnaki/studytrack-api/internal/controller/dto"
	"github.com/shnaki/studytrack-api/internal/usecase"
)

type createProjectInput struct {
	UserID string `path:"userId" doc:"User ID"`
	Body   dto.CreateProjectRequest
}

type createProjectOutput struct {
	Body dto.ProjectResponse
}

type listProjectsInput struct {
	UserID string `path:"userId" doc:"User ID"`
}

type listProjectsOutput struct {
	Body []dto.ProjectResponse
}

type updateProjectInput struct {
	ID   string `path:"id" doc:"Project ID"`
	Body dto.UpdateProjectRequest
}

type updateProjectOutput struct {
	Body dto.ProjectResponse
}

type deleteProjectInput struct {
	ID string `path:"id" doc:"Project ID"`
}

// RegisterProjectRoutes registers project-related routes to the Huma API.
func RegisterProjectRoutes(api huma.API, uc *usecase.ProjectUsecase) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-project",
		Method:        http.MethodPost,
		Path:          "/users/{userId}/projects",
		Summary:       "Create a new project",
		Tags:          []string{"Projects"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *createProjectInput) (*createProjectOutput, error) {
		project, err := uc.CreateProject(ctx, input.UserID, input.Body.Name)
		if err != nil {
			return nil, toHTTPError(err)
		}
		return &createProjectOutput{Body: dto.ToProjectResponse(project)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-projects",
		Method:      http.MethodGet,
		Path:        "/users/{userId}/projects",
		Summary:     "List projects for a user",
		Tags:        []string{"Projects"},
	}, func(ctx context.Context, input *listProjectsInput) (*listProjectsOutput, error) {
		projects, err := uc.ListProjects(ctx, input.UserID)
		if err != nil {
			return nil, toHTTPError(err)
		}
		return &listProjectsOutput{Body: dto.ToProjectResponseList(projects)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-project",
		Method:      http.MethodPut,
		Path:        "/projects/{id}",
		Summary:     "Update a project",
		Tags:        []string{"Projects"},
	}, func(ctx context.Context, input *updateProjectInput) (*updateProjectOutput, error) {
		project, err := uc.UpdateProject(ctx, input.ID, input.Body.Name)
		if err != nil {
			return nil, toHTTPError(err)
		}
		return &updateProjectOutput{Body: dto.ToProjectResponse(project)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-project",
		Method:        http.MethodDelete,
		Path:          "/projects/{id}",
		Summary:       "Delete a project",
		Tags:          []string{"Projects"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *deleteProjectInput) (*struct{}, error) {
		if err := uc.DeleteProject(ctx, input.ID); err != nil {
			return nil, toHTTPError(err)
		}
		return nil, nil
	})
}
