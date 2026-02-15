package domain

import "time"

// StudyLog represents a record of study time for a specific project.
type StudyLog struct {
	ID        string
	UserID    string
	ProjectID string
	StudiedAt time.Time
	Minutes   int
	Note      string
	CreatedAt time.Time
}

// NewStudyLog creates a new StudyLog entity.
func NewStudyLog(id, userID, projectID string, studiedAt time.Time, minutes int, note string) (*StudyLog, error) {
	if userID == "" {
		return nil, ErrValidation("user ID is required")
	}
	if projectID == "" {
		return nil, ErrValidation("project ID is required")
	}
	if err := validateMinutes(minutes); err != nil {
		return nil, err
	}
	return &StudyLog{
		ID:        id,
		UserID:    userID,
		ProjectID: projectID,
		StudiedAt: studiedAt,
		Minutes:   minutes,
		Note:      note,
		CreatedAt: time.Now(),
	}, nil
}

// ReconstructStudyLog reconstructs a StudyLog entity from existing data.
func ReconstructStudyLog(id, userID, projectID string, studiedAt time.Time, minutes int, note string, createdAt time.Time) *StudyLog {
	return &StudyLog{
		ID:        id,
		UserID:    userID,
		ProjectID: projectID,
		StudiedAt: studiedAt,
		Minutes:   minutes,
		Note:      note,
		CreatedAt: createdAt,
	}
}

func validateMinutes(minutes int) error {
	if minutes <= 0 {
		return ErrValidation("minutes must be greater than 0")
	}
	if minutes > 1440 {
		return ErrValidation("minutes must be 1440 or less")
	}
	return nil
}
