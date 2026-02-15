package controller

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/shnaki/studytrack-api/internal/controller/dto"
	"github.com/shnaki/studytrack-api/internal/usecase"
)

type createNoteInput struct {
	UserID    string `path:"userId" doc:"User ID"`
	ProjectID string `path:"projectId" doc:"Project ID"`
	Body      dto.CreateNoteRequest
}

type createNoteOutput struct {
	Body dto.NoteResponse
}

type listNotesInput struct {
	UserID    string `path:"userId" doc:"User ID"`
	ProjectID string `path:"projectId" doc:"Project ID"`
}

type listNotesOutput struct {
	Body []dto.NoteResponse
}

type getNoteInput struct {
	ID string `path:"id" doc:"Note ID"`
}

type getNoteOutput struct {
	Body dto.NoteResponse
}

type updateNoteInput struct {
	ID   string `path:"id" doc:"Note ID"`
	Body dto.UpdateNoteRequest
}

type updateNoteOutput struct {
	Body dto.NoteResponse
}

type deleteNoteInput struct {
	ID string `path:"id" doc:"Note ID"`
}

// RegisterNoteRoutes registers note-related routes to the Huma API.
func RegisterNoteRoutes(api huma.API, uc *usecase.NoteUsecase) {
	huma.Register(api, huma.Operation{
		OperationID:   "create-note",
		Method:        http.MethodPost,
		Path:          "/users/{userId}/projects/{projectId}/notes",
		Summary:       "Create a new note",
		Tags:          []string{"Notes"},
		DefaultStatus: http.StatusCreated,
	}, func(ctx context.Context, input *createNoteInput) (*createNoteOutput, error) {
		note, err := uc.CreateNote(ctx, input.UserID, input.ProjectID, input.Body.Title, input.Body.Content, input.Body.Tags)
		if err != nil {
			return nil, toHTTPError(err)
		}
		return &createNoteOutput{Body: dto.ToNoteResponse(note)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list-notes",
		Method:      http.MethodGet,
		Path:        "/users/{userId}/projects/{projectId}/notes",
		Summary:     "List notes for a project",
		Tags:        []string{"Notes"},
	}, func(ctx context.Context, input *listNotesInput) (*listNotesOutput, error) {
		notes, err := uc.ListNotes(ctx, input.UserID, input.ProjectID)
		if err != nil {
			return nil, toHTTPError(err)
		}
		return &listNotesOutput{Body: dto.ToNoteResponseList(notes)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get-note",
		Method:      http.MethodGet,
		Path:        "/notes/{id}",
		Summary:     "Get a note by ID",
		Tags:        []string{"Notes"},
	}, func(ctx context.Context, input *getNoteInput) (*getNoteOutput, error) {
		note, err := uc.GetNote(ctx, input.ID)
		if err != nil {
			return nil, toHTTPError(err)
		}
		return &getNoteOutput{Body: dto.ToNoteResponse(note)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-note",
		Method:      http.MethodPut,
		Path:        "/notes/{id}",
		Summary:     "Update a note",
		Tags:        []string{"Notes"},
	}, func(ctx context.Context, input *updateNoteInput) (*updateNoteOutput, error) {
		note, err := uc.UpdateNote(ctx, input.ID, input.Body.Title, input.Body.Content, input.Body.Tags)
		if err != nil {
			return nil, toHTTPError(err)
		}
		return &updateNoteOutput{Body: dto.ToNoteResponse(note)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID:   "delete-note",
		Method:        http.MethodDelete,
		Path:          "/notes/{id}",
		Summary:       "Delete a note",
		Tags:          []string{"Notes"},
		DefaultStatus: http.StatusNoContent,
	}, func(ctx context.Context, input *deleteNoteInput) (*struct{}, error) {
		if err := uc.DeleteNote(ctx, input.ID); err != nil {
			return nil, toHTTPError(err)
		}
		return nil, nil
	})
}
