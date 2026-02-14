package domain

import "time"

type Goal struct {
	ID                   string
	UserID               string
	SubjectID            string
	TargetMinutesPerWeek int
	StartDate            time.Time
	EndDate              *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func NewGoal(id, userID, subjectID string, targetMinutesPerWeek int, startDate time.Time, endDate *time.Time) (*Goal, error) {
	if userID == "" {
		return nil, ErrValidation("user ID is required")
	}
	if subjectID == "" {
		return nil, ErrValidation("subject ID is required")
	}
	if err := validateTargetMinutes(targetMinutesPerWeek); err != nil {
		return nil, err
	}
	if endDate != nil && endDate.Before(startDate) {
		return nil, ErrValidation("end date must be after start date")
	}
	now := time.Now()
	return &Goal{
		ID:                   id,
		UserID:               userID,
		SubjectID:            subjectID,
		TargetMinutesPerWeek: targetMinutesPerWeek,
		StartDate:            startDate,
		EndDate:              endDate,
		CreatedAt:            now,
		UpdatedAt:            now,
	}, nil
}

func ReconstructGoal(id, userID, subjectID string, targetMinutesPerWeek int, startDate time.Time, endDate *time.Time, createdAt, updatedAt time.Time) *Goal {
	return &Goal{
		ID:                   id,
		UserID:               userID,
		SubjectID:            subjectID,
		TargetMinutesPerWeek: targetMinutesPerWeek,
		StartDate:            startDate,
		EndDate:              endDate,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
	}
}

func validateTargetMinutes(minutes int) error {
	if minutes <= 0 {
		return ErrValidation("target minutes per week must be greater than 0")
	}
	return nil
}
