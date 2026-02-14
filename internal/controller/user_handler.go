package controller

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/shnaki/studytrack-api/internal/controller/dto"
	"github.com/shnaki/studytrack-api/internal/usecase"
)

type createUserInput struct {
	Body dto.CreateUserRequest
}

type createUserOutput struct {
	Body dto.UserResponse
}

type getUserInput struct {
	ID string `path:"id" doc:"User ID"`
}

type getUserOutput struct {
	Body dto.UserResponse
}

func RegisterUserRoutes(api huma.API, uc *usecase.UserUsecase) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-user",
		Method:        http.MethodPost,
		Path:          "/users",
		Summary:       "Create a new user",
		Tags:          []string{"Users"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *createUserInput) (*createUserOutput, error) {
		user, err := uc.CreateUser(ctx, input.Body.Name)
		if err != nil {
			return nil, toHTTPError(err)
		}
		return &createUserOutput{Body: dto.ToUserResponse(user)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/users/{id}",
		Summary:     "Get a user by ID",
		Tags:        []string{"Users"},
	}, func(ctx context.Context, input *getUserInput) (*getUserOutput, error) {
		user, err := uc.GetUser(ctx, input.ID)
		if err != nil {
			return nil, toHTTPError(err)
		}
		return &getUserOutput{Body: dto.ToUserResponse(user)}, nil
	})
}
