package usecase_test

import (
	"context"

	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase/port"
)

// --- Mock UserRepository (map-based) ---

type mockUserRepository struct {
	users map[string]*domain.User
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{users: make(map[string]*domain.User)}
}

func (m *mockUserRepository) Create(_ context.Context, user *domain.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) FindByID(_ context.Context, id string) (*domain.User, error) {
	u, ok := m.users[id]
	if !ok {
		return nil, domain.ErrNotFound("user")
	}
	return u, nil
}

// --- Mock SubjectRepository (map-based, with duplicate check) ---

type mockSubjectRepository struct {
	subjects map[string]*domain.Subject
}

func newMockSubjectRepository() *mockSubjectRepository {
	return &mockSubjectRepository{subjects: make(map[string]*domain.Subject)}
}

func (m *mockSubjectRepository) Create(_ context.Context, s *domain.Subject) error {
	// Check for duplicate (userID, name) combination
	for _, existing := range m.subjects {
		if existing.UserID == s.UserID && existing.Name == s.Name {
			return domain.ErrConflict("subject with this name already exists for this user")
		}
	}
	m.subjects[s.ID] = s
	return nil
}

func (m *mockSubjectRepository) FindByID(_ context.Context, id string) (*domain.Subject, error) {
	s, ok := m.subjects[id]
	if !ok {
		return nil, domain.ErrNotFound("subject")
	}
	return s, nil
}

func (m *mockSubjectRepository) FindByUserID(_ context.Context, userID string) ([]*domain.Subject, error) {
	var result []*domain.Subject
	for _, s := range m.subjects {
		if s.UserID == userID {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *mockSubjectRepository) Update(_ context.Context, subject *domain.Subject) error {
	if _, ok := m.subjects[subject.ID]; !ok {
		return domain.ErrNotFound("subject")
	}
	m.subjects[subject.ID] = subject
	return nil
}

func (m *mockSubjectRepository) Delete(_ context.Context, id string) error {
	if _, ok := m.subjects[id]; !ok {
		return domain.ErrNotFound("subject")
	}
	delete(m.subjects, id)
	return nil
}

// --- Mock StudyLogRepository (slice-based) ---

type mockStudyLogRepository struct {
	logs []*domain.StudyLog
}

func newMockStudyLogRepository() *mockStudyLogRepository {
	return &mockStudyLogRepository{}
}

func (m *mockStudyLogRepository) Create(_ context.Context, l *domain.StudyLog) error {
	m.logs = append(m.logs, l)
	return nil
}

func (m *mockStudyLogRepository) FindByID(_ context.Context, id string) (*domain.StudyLog, error) {
	for _, l := range m.logs {
		if l.ID == id {
			return l, nil
		}
	}
	return nil, domain.ErrNotFound("study log")
}

func (m *mockStudyLogRepository) FindByUserID(_ context.Context, userID string, filter port.StudyLogFilter) ([]*domain.StudyLog, error) {
	var result []*domain.StudyLog
	for _, l := range m.logs {
		if l.UserID != userID {
			continue
		}
		if filter.From != nil && l.StudiedAt.Before(*filter.From) {
			continue
		}
		if filter.To != nil && !l.StudiedAt.Before(*filter.To) {
			continue
		}
		if filter.SubjectID != nil && l.SubjectID != *filter.SubjectID {
			continue
		}
		result = append(result, l)
	}
	return result, nil
}

func (m *mockStudyLogRepository) Delete(_ context.Context, id string) error {
	for i, l := range m.logs {
		if l.ID == id {
			m.logs = append(m.logs[:i], m.logs[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound("study log")
}

// --- Mock GoalRepository (slice-based, with upsert logic) ---

type mockGoalRepository struct {
	goals []*domain.Goal
}

func newMockGoalRepository() *mockGoalRepository {
	return &mockGoalRepository{}
}

func (m *mockGoalRepository) Upsert(_ context.Context, g *domain.Goal) error {
	// Check if a goal with the same user_id + subject_id exists; update if so
	for i, existing := range m.goals {
		if existing.UserID == g.UserID && existing.SubjectID == g.SubjectID {
			m.goals[i] = g
			return nil
		}
	}
	m.goals = append(m.goals, g)
	return nil
}

func (m *mockGoalRepository) FindByUserID(_ context.Context, userID string) ([]*domain.Goal, error) {
	var result []*domain.Goal
	for _, g := range m.goals {
		if g.UserID == userID {
			result = append(result, g)
		}
	}
	return result, nil
}
