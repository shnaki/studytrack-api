package domain

import "time"

type Subject struct {
	ID        string
	UserID    string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewSubject(id, userID, name string) (*Subject, error) {
	if err := validateSubjectName(name); err != nil {
		return nil, err
	}
	if userID == "" {
		return nil, ErrValidation("user ID is required")
	}
	now := time.Now()
	return &Subject{
		ID:        id,
		UserID:    userID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func ReconstructSubject(id, userID, name string, createdAt, updatedAt time.Time) *Subject {
	return &Subject{
		ID:        id,
		UserID:    userID,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func (s *Subject) UpdateName(name string) error {
	if err := validateSubjectName(name); err != nil {
		return err
	}
	s.Name = name
	s.UpdatedAt = time.Now()
	return nil
}

func validateSubjectName(name string) error {
	if name == "" {
		return ErrValidation("subject name is required")
	}
	if len(name) > 200 {
		return ErrValidation("subject name must be 200 characters or less")
	}
	return nil
}
