package domain

import "time"

// Project represents a learning project.
type Project struct {
	ID        string
	UserID    string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewProject creates a new Project entity.
func NewProject(id, userID, name string) (*Project, error) {
	if err := validateProjectName(name); err != nil {
		return nil, err
	}
	if userID == "" {
		return nil, ErrValidation("user ID is required")
	}
	now := time.Now()
	return &Project{
		ID:        id,
		UserID:    userID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// ReconstructProject reconstructs a Project entity from existing data.
func ReconstructProject(id, userID, name string, createdAt, updatedAt time.Time) *Project {
	return &Project{
		ID:        id,
		UserID:    userID,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// UpdateName updates the name of the project.
func (p *Project) UpdateName(name string) error {
	if err := validateProjectName(name); err != nil {
		return err
	}
	p.Name = name
	p.UpdatedAt = time.Now()
	return nil
}

func validateProjectName(name string) error {
	if name == "" {
		return ErrValidation("project name is required")
	}
	if len(name) > 200 {
		return ErrValidation("project name must be 200 characters or less")
	}
	return nil
}
