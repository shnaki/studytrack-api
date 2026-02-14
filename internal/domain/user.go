package domain

import "time"

// User represents a system user.
type User struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new User entity.
func NewUser(id, name string) (*User, error) {
	if err := validateUserName(name); err != nil {
		return nil, err
	}
	now := time.Now()
	return &User{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// ReconstructUser reconstructs a User entity from existing data.
func ReconstructUser(id, name string, createdAt, updatedAt time.Time) *User {
	return &User{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func validateUserName(name string) error {
	if name == "" {
		return ErrValidation("user name is required")
	}
	if len(name) > 100 {
		return ErrValidation("user name must be 100 characters or less")
	}
	return nil
}
