package domain

import "time"

type StudyLog struct {
	ID        string
	UserID    string
	SubjectID string
	StudiedAt time.Time
	Minutes   int
	Note      string
	CreatedAt time.Time
}

func NewStudyLog(id, userID, subjectID string, studiedAt time.Time, minutes int, note string) (*StudyLog, error) {
	if userID == "" {
		return nil, ErrValidation("user ID is required")
	}
	if subjectID == "" {
		return nil, ErrValidation("subject ID is required")
	}
	if err := validateMinutes(minutes); err != nil {
		return nil, err
	}
	return &StudyLog{
		ID:        id,
		UserID:    userID,
		SubjectID: subjectID,
		StudiedAt: studiedAt,
		Minutes:   minutes,
		Note:      note,
		CreatedAt: time.Now(),
	}, nil
}

func ReconstructStudyLog(id, userID, subjectID string, studiedAt time.Time, minutes int, note string, createdAt time.Time) *StudyLog {
	return &StudyLog{
		ID:        id,
		UserID:    userID,
		SubjectID: subjectID,
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
