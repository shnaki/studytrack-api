package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase/port"
)

// ProjectUsecase provides methods for managing projects.
type ProjectUsecase struct {
	projectRepo port.ProjectRepository
	userRepo    port.UserRepository
}

// NewProjectUsecase creates a new ProjectUsecase.
func NewProjectUsecase(projectRepo port.ProjectRepository, userRepo port.UserRepository) *ProjectUsecase {
	return &ProjectUsecase{
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

// CreateProject creates a new project for a user.
func (u *ProjectUsecase) CreateProject(ctx context.Context, userID, name string) (*domain.Project, error) {
	if _, err := u.userRepo.FindByID(ctx, userID); err != nil {
		return nil, err
	}
	id := uuid.New().String()
	project, err := domain.NewProject(id, userID, name)
	if err != nil {
		return nil, err
	}
	if err := u.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

// ListProjects returns all projects for a user.
func (u *ProjectUsecase) ListProjects(ctx context.Context, userID string) ([]*domain.Project, error) {
	if _, err := u.userRepo.FindByID(ctx, userID); err != nil {
		return nil, err
	}
	return u.projectRepo.FindByUserID(ctx, userID)
}

// UpdateProject updates an existing project.
func (u *ProjectUsecase) UpdateProject(ctx context.Context, id, name string) (*domain.Project, error) {
	project, err := u.projectRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := project.UpdateName(name); err != nil {
		return nil, err
	}
	if err := u.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}
	return project, nil
}

// DeleteProject deletes a project by ID.
func (u *ProjectUsecase) DeleteProject(ctx context.Context, id string) error {
	if _, err := u.projectRepo.FindByID(ctx, id); err != nil {
		return err
	}
	return u.projectRepo.Delete(ctx, id)
}
