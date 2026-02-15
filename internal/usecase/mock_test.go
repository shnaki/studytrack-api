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

// --- Mock ProjectRepository (map-based, with duplicate check) ---

type mockProjectRepository struct {
	projects map[string]*domain.Project
}

func newMockProjectRepository() *mockProjectRepository {
	return &mockProjectRepository{projects: make(map[string]*domain.Project)}
}

func (m *mockProjectRepository) Create(_ context.Context, p *domain.Project) error {
	// Check for duplicate (userID, name) combination
	for _, existing := range m.projects {
		if existing.UserID == p.UserID && existing.Name == p.Name {
			return domain.ErrConflict("project with this name already exists for this user")
		}
	}
	m.projects[p.ID] = p
	return nil
}

func (m *mockProjectRepository) FindByID(_ context.Context, id string) (*domain.Project, error) {
	p, ok := m.projects[id]
	if !ok {
		return nil, domain.ErrNotFound("project")
	}
	return p, nil
}

func (m *mockProjectRepository) FindByUserID(_ context.Context, userID string) ([]*domain.Project, error) {
	var result []*domain.Project
	for _, p := range m.projects {
		if p.UserID == userID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (m *mockProjectRepository) Update(_ context.Context, project *domain.Project) error {
	if _, ok := m.projects[project.ID]; !ok {
		return domain.ErrNotFound("project")
	}
	m.projects[project.ID] = project
	return nil
}

func (m *mockProjectRepository) Delete(_ context.Context, id string) error {
	if _, ok := m.projects[id]; !ok {
		return domain.ErrNotFound("project")
	}
	delete(m.projects, id)
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
		if filter.ProjectID != nil && l.ProjectID != *filter.ProjectID {
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
	// Check if a goal with the same user_id + project_id exists; update if so
	for i, existing := range m.goals {
		if existing.UserID == g.UserID && existing.ProjectID == g.ProjectID {
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

// --- Mock NoteRepository (map-based) ---

type mockNoteRepository struct {
	notes map[string]*domain.Note
}

func newMockNoteRepository() *mockNoteRepository {
	return &mockNoteRepository{notes: make(map[string]*domain.Note)}
}

func (m *mockNoteRepository) Create(_ context.Context, note *domain.Note) error {
	m.notes[note.ID] = note
	return nil
}

func (m *mockNoteRepository) FindByID(_ context.Context, id string) (*domain.Note, error) {
	n, ok := m.notes[id]
	if !ok {
		return nil, domain.ErrNotFound("note")
	}
	return n, nil
}

func (m *mockNoteRepository) FindByProjectID(_ context.Context, projectID string) ([]*domain.Note, error) {
	var result []*domain.Note
	for _, n := range m.notes {
		if n.ProjectID == projectID {
			result = append(result, n)
		}
	}
	return result, nil
}

func (m *mockNoteRepository) Update(_ context.Context, note *domain.Note) error {
	if _, ok := m.notes[note.ID]; !ok {
		return domain.ErrNotFound("note")
	}
	m.notes[note.ID] = note
	return nil
}

func (m *mockNoteRepository) Delete(_ context.Context, id string) error {
	if _, ok := m.notes[id]; !ok {
		return domain.ErrNotFound("note")
	}
	delete(m.notes, id)
	return nil
}
