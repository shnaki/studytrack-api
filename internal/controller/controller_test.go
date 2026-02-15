package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shnaki/studytrack-api/internal/controller"
	"github.com/shnaki/studytrack-api/internal/domain"
	"github.com/shnaki/studytrack-api/internal/usecase"
	"github.com/shnaki/studytrack-api/internal/usecase/port"
)

// --- Mock Repositories ---

type mockUserRepository struct {
	users map[string]*domain.User
}

func newMockUserRepo() *mockUserRepository {
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

type mockProjectRepository struct {
	projects map[string]*domain.Project
}

func newMockProjectRepo() *mockProjectRepository {
	return &mockProjectRepository{projects: make(map[string]*domain.Project)}
}

func (m *mockProjectRepository) Create(_ context.Context, p *domain.Project) error {
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

func (m *mockProjectRepository) Update(_ context.Context, p *domain.Project) error {
	m.projects[p.ID] = p
	return nil
}

func (m *mockProjectRepository) Delete(_ context.Context, id string) error {
	delete(m.projects, id)
	return nil
}

type mockStudyLogRepository struct {
	logs map[string]*domain.StudyLog
}

func newMockStudyLogRepo() *mockStudyLogRepository {
	return &mockStudyLogRepository{logs: make(map[string]*domain.StudyLog)}
}

func (m *mockStudyLogRepository) Create(_ context.Context, l *domain.StudyLog) error {
	m.logs[l.ID] = l
	return nil
}

func (m *mockStudyLogRepository) FindByID(_ context.Context, id string) (*domain.StudyLog, error) {
	l, ok := m.logs[id]
	if !ok {
		return nil, domain.ErrNotFound("study log")
	}
	return l, nil
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
	delete(m.logs, id)
	return nil
}

type mockGoalRepository struct {
	goals map[string]*domain.Goal
}

func newMockGoalRepo() *mockGoalRepository {
	return &mockGoalRepository{goals: make(map[string]*domain.Goal)}
}

func (m *mockGoalRepository) Upsert(_ context.Context, g *domain.Goal) error {
	m.goals[g.UserID+":"+g.ProjectID] = g
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

type mockNoteRepository struct {
	notes map[string]*domain.Note
}

func newMockNoteRepo() *mockNoteRepository {
	return &mockNoteRepository{notes: make(map[string]*domain.Note)}
}

func (m *mockNoteRepository) Create(_ context.Context, n *domain.Note) error {
	m.notes[n.ID] = n
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

func (m *mockNoteRepository) Update(_ context.Context, n *domain.Note) error {
	if _, ok := m.notes[n.ID]; !ok {
		return domain.ErrNotFound("note")
	}
	m.notes[n.ID] = n
	return nil
}

func (m *mockNoteRepository) Delete(_ context.Context, id string) error {
	if _, ok := m.notes[id]; !ok {
		return domain.ErrNotFound("note")
	}
	delete(m.notes, id)
	return nil
}

// --- Helpers ---

func setupRouter(t *testing.T) (http.Handler, *mockUserRepository, *mockProjectRepository, *mockStudyLogRepository, *mockGoalRepository) {
	t.Helper()
	userRepo := newMockUserRepo()
	projectRepo := newMockProjectRepo()
	studyLogRepo := newMockStudyLogRepo()
	goalRepo := newMockGoalRepo()
	noteRepo := newMockNoteRepo()

	usecases := &controller.Usecases{
		User:     usecase.NewUserUsecase(userRepo),
		Project:  usecase.NewProjectUsecase(projectRepo, userRepo),
		StudyLog: usecase.NewStudyLogUsecase(studyLogRepo, userRepo, projectRepo),
		Goal:     usecase.NewGoalUsecase(goalRepo, userRepo, projectRepo),
		Stats:    usecase.NewStatsUsecase(studyLogRepo, goalRepo, projectRepo),
		Note:     usecase.NewNoteUsecase(noteRepo, projectRepo, userRepo),
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := controller.NewRouter(usecases, []string{"*"}, logger)
	return router, userRepo, projectRepo, studyLogRepo, goalRepo
}

func jsonRequest(method, path string, body any) *http.Request {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			panic(err) // This is for test setup, should not happen with valid input
		}
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	return req
}

func doRequest(handler http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

func parseJSON(t *testing.T, rr *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.NewDecoder(rr.Body).Decode(v); err != nil {
		t.Fatalf("failed to parse response body: %v\nbody: %s", err, rr.Body.String())
	}
}

// --- User Tests ---

func TestCreateUser_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	req := jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"})
	rr := doRequest(handler, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	var resp map[string]any
	parseJSON(t, rr, &resp)

	if resp["name"] != "Alice" {
		t.Errorf("expected name 'Alice', got '%v'", resp["name"])
	}
	if resp["id"] == nil || resp["id"] == "" {
		t.Error("expected id to be set")
	}
	if resp["createdAt"] == nil {
		t.Error("expected createdAt to be set")
	}
	if resp["updatedAt"] == nil {
		t.Error("expected updatedAt to be set")
	}
}

func TestCreateUser_ValidationError_EmptyName(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Huma validates minLength:"1" on the struct tag, so empty name yields 422
	req := jsonRequest("POST", "/v1/users", map[string]string{"name": ""})
	rr := doRequest(handler, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusUnprocessableEntity, rr.Code, rr.Body.String())
	}
}

func TestGetUser_Found(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// First create a user
	createReq := jsonRequest("POST", "/v1/users", map[string]string{"name": "Bob"})
	createRR := doRequest(handler, createReq)
	if createRR.Code != http.StatusCreated {
		t.Fatalf("setup: create user failed with status %d", createRR.Code)
	}

	var created map[string]any
	parseJSON(t, createRR, &created)
	userID := created["id"].(string)

	// Now get the user
	getReq := jsonRequest("GET", "/v1/users/"+userID, nil)
	getRR := doRequest(handler, getReq)

	if getRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusOK, getRR.Code, getRR.Body.String())
	}

	var resp map[string]any
	parseJSON(t, getRR, &resp)
	if resp["name"] != "Bob" {
		t.Errorf("expected name 'Bob', got '%v'", resp["name"])
	}
	if resp["id"] != userID {
		t.Errorf("expected id '%s', got '%v'", userID, resp["id"])
	}
}

func TestGetUser_NotFound(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	req := jsonRequest("GET", "/v1/users/nonexistent-id", nil)
	rr := doRequest(handler, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}

// --- Project Tests ---

func TestCreateProject_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user first
	createUserReq := jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"})
	createUserRR := doRequest(handler, createUserReq)
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	// Create project
	createProjReq := jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "Mathematics"})
	createProjRR := doRequest(handler, createProjReq)

	if createProjRR.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusCreated, createProjRR.Code, createProjRR.Body.String())
	}

	var proj map[string]any
	parseJSON(t, createProjRR, &proj)
	if proj["name"] != "Mathematics" {
		t.Errorf("expected name 'Mathematics', got '%v'", proj["name"])
	}
	if proj["userId"] != userID {
		t.Errorf("expected userId '%s', got '%v'", userID, proj["userId"])
	}
	if proj["id"] == nil || proj["id"] == "" {
		t.Error("expected project id to be set")
	}
}

func TestListProjects_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user
	createUserReq := jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"})
	createUserRR := doRequest(handler, createUserReq)
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	// Create two projects
	doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "Math"}))
	doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "English"}))

	// List projects
	listReq := jsonRequest("GET", "/v1/users/"+userID+"/projects", nil)
	listRR := doRequest(handler, listReq)

	if listRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusOK, listRR.Code, listRR.Body.String())
	}

	var projects []map[string]any
	parseJSON(t, listRR, &projects)
	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}
}

func TestUpdateProject_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user and project
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createProjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "Math"}))
	var proj map[string]any
	parseJSON(t, createProjRR, &proj)
	projectID := proj["id"].(string)

	// Update project
	updateReq := jsonRequest("PUT", "/v1/projects/"+projectID, map[string]string{"name": "Advanced Math"})
	updateRR := doRequest(handler, updateReq)

	if updateRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusOK, updateRR.Code, updateRR.Body.String())
	}

	var updated map[string]any
	parseJSON(t, updateRR, &updated)
	if updated["name"] != "Advanced Math" {
		t.Errorf("expected name 'Advanced Math', got '%v'", updated["name"])
	}
}

func TestDeleteProject_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user and project
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createProjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "Math"}))
	var proj map[string]any
	parseJSON(t, createProjRR, &proj)
	projectID := proj["id"].(string)

	// Delete project
	deleteReq := jsonRequest("DELETE", "/v1/projects/"+projectID, nil)
	deleteRR := doRequest(handler, deleteReq)

	if deleteRR.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNoContent, deleteRR.Code, deleteRR.Body.String())
	}
}

func TestListProjects_UserNotFound(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	req := jsonRequest("GET", "/v1/users/nonexistent-user/projects", nil)
	rr := doRequest(handler, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}

// --- StudyLog Tests ---

func TestCreateStudyLog_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	// Create project
	createProjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "Math"}))
	var proj map[string]any
	parseJSON(t, createProjRR, &proj)
	projectID := proj["id"].(string)

	// Create study log
	studiedAt := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	body := map[string]any{
		"projectId": projectID,
		"studiedAt": studiedAt.Format(time.RFC3339),
		"minutes":   60,
		"note":      "Chapter 1",
	}
	createLogReq := jsonRequest("POST", "/v1/users/"+userID+"/study-logs", body)
	createLogRR := doRequest(handler, createLogReq)

	if createLogRR.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusCreated, createLogRR.Code, createLogRR.Body.String())
	}

	var logResp map[string]any
	parseJSON(t, createLogRR, &logResp)
	if logResp["userId"] != userID {
		t.Errorf("expected userId '%s', got '%v'", userID, logResp["userId"])
	}
	if logResp["projectId"] != projectID {
		t.Errorf("expected projectId '%s', got '%v'", projectID, logResp["projectId"])
	}
	if int(logResp["minutes"].(float64)) != 60 {
		t.Errorf("expected minutes 60, got %v", logResp["minutes"])
	}
	if logResp["note"] != "Chapter 1" {
		t.Errorf("expected note 'Chapter 1', got '%v'", logResp["note"])
	}
}

func TestListStudyLogs_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user and project
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createProjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "Math"}))
	var proj map[string]any
	parseJSON(t, createProjRR, &proj)
	projectID := proj["id"].(string)

	// Create two study logs
	studiedAt1 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	studiedAt2 := time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC)
	doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/study-logs", map[string]any{
		"projectId": projectID, "studiedAt": studiedAt1.Format(time.RFC3339), "minutes": 60, "note": "session 1",
	}))
	doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/study-logs", map[string]any{
		"projectId": projectID, "studiedAt": studiedAt2.Format(time.RFC3339), "minutes": 45, "note": "session 2",
	}))

	// List study logs
	listReq := jsonRequest("GET", "/v1/users/"+userID+"/study-logs", nil)
	listRR := doRequest(handler, listReq)

	if listRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusOK, listRR.Code, listRR.Body.String())
	}

	var logs []map[string]any
	parseJSON(t, listRR, &logs)
	if len(logs) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(logs))
	}
}

func TestListStudyLogs_WithDateFilter(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user and project
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createProjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "Math"}))
	var proj map[string]any
	parseJSON(t, createProjRR, &proj)
	projectID := proj["id"].(string)

	// Create logs on different dates
	jan15 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	feb10 := time.Date(2024, 2, 10, 10, 0, 0, 0, time.UTC)
	doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/study-logs", map[string]any{
		"projectId": projectID, "studiedAt": jan15.Format(time.RFC3339), "minutes": 60, "note": "january",
	}))
	doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/study-logs", map[string]any{
		"projectId": projectID, "studiedAt": feb10.Format(time.RFC3339), "minutes": 45, "note": "february",
	}))

	// Filter to January only
	listReq := jsonRequest("GET", "/v1/users/"+userID+"/study-logs?from=2024-01-01&to=2024-01-31", nil)
	listRR := doRequest(handler, listReq)

	if listRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusOK, listRR.Code, listRR.Body.String())
	}

	var logs []map[string]any
	parseJSON(t, listRR, &logs)
	if len(logs) != 1 {
		t.Fatalf("expected 1 log in January, got %d", len(logs))
	}
	if logs[0]["note"] != "january" {
		t.Errorf("expected note 'january', got '%v'", logs[0]["note"])
	}
}

func TestDeleteStudyLog_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user, project, and study log
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createProjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "Math"}))
	var proj map[string]any
	parseJSON(t, createProjRR, &proj)
	projectID := proj["id"].(string)

	studiedAt := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	createLogRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/study-logs", map[string]any{
		"projectId": projectID, "studiedAt": studiedAt.Format(time.RFC3339), "minutes": 60, "note": "to delete",
	}))
	var logResp map[string]any
	parseJSON(t, createLogRR, &logResp)
	logID := logResp["id"].(string)

	// Delete study log
	deleteReq := jsonRequest("DELETE", "/v1/study-logs/"+logID, nil)
	deleteRR := doRequest(handler, deleteReq)

	if deleteRR.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNoContent, deleteRR.Code, deleteRR.Body.String())
	}
}

// --- Goal Tests ---

func TestUpsertGoal_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user and project
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createProjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "Math"}))
	var proj map[string]any
	parseJSON(t, createProjRR, &proj)
	projectID := proj["id"].(string)

	// Upsert goal
	body := map[string]any{
		"targetMinutesPerWeek": 300,
		"startDate":            "2024-01-01",
	}
	upsertReq := jsonRequest("PUT", "/v1/users/"+userID+"/goals/"+projectID, body)
	upsertRR := doRequest(handler, upsertReq)

	if upsertRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusOK, upsertRR.Code, upsertRR.Body.String())
	}

	var goalResp map[string]any
	parseJSON(t, upsertRR, &goalResp)
	if goalResp["userId"] != userID {
		t.Errorf("expected userId '%s', got '%v'", userID, goalResp["userId"])
	}
	if goalResp["projectId"] != projectID {
		t.Errorf("expected projectId '%s', got '%v'", projectID, goalResp["projectId"])
	}
	if int(goalResp["targetMinutesPerWeek"].(float64)) != 300 {
		t.Errorf("expected targetMinutesPerWeek 300, got %v", goalResp["targetMinutesPerWeek"])
	}
	if goalResp["startDate"] != "2024-01-01" {
		t.Errorf("expected startDate '2024-01-01', got '%v'", goalResp["startDate"])
	}
	if goalResp["endDate"] != nil {
		t.Errorf("expected endDate nil, got '%v'", goalResp["endDate"])
	}
}

func TestUpsertGoal_WithEndDate(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user and project
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createProjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "Math"}))
	var proj map[string]any
	parseJSON(t, createProjRR, &proj)
	projectID := proj["id"].(string)

	endDate := "2024-06-30"
	body := map[string]any{
		"targetMinutesPerWeek": 200,
		"startDate":            "2024-01-01",
		"endDate":              endDate,
	}
	upsertReq := jsonRequest("PUT", "/v1/users/"+userID+"/goals/"+projectID, body)
	upsertRR := doRequest(handler, upsertReq)

	if upsertRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusOK, upsertRR.Code, upsertRR.Body.String())
	}

	var goalResp map[string]any
	parseJSON(t, upsertRR, &goalResp)
	if goalResp["endDate"] != "2024-06-30" {
		t.Errorf("expected endDate '2024-06-30', got '%v'", goalResp["endDate"])
	}
}

func TestListGoals_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user and two projects
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createProj1RR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "Math"}))
	var proj1 map[string]any
	parseJSON(t, createProj1RR, &proj1)
	projectID1 := proj1["id"].(string)

	createProj2RR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "English"}))
	var proj2 map[string]any
	parseJSON(t, createProj2RR, &proj2)
	projectID2 := proj2["id"].(string)

	// Upsert goals for both projects
	doRequest(handler, jsonRequest("PUT", "/v1/users/"+userID+"/goals/"+projectID1, map[string]any{
		"targetMinutesPerWeek": 300, "startDate": "2024-01-01",
	}))
	doRequest(handler, jsonRequest("PUT", "/v1/users/"+userID+"/goals/"+projectID2, map[string]any{
		"targetMinutesPerWeek": 150, "startDate": "2024-01-01",
	}))

	// List goals
	listReq := jsonRequest("GET", "/v1/users/"+userID+"/goals", nil)
	listRR := doRequest(handler, listReq)

	if listRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusOK, listRR.Code, listRR.Body.String())
	}

	var goals []map[string]any
	parseJSON(t, listRR, &goals)
	if len(goals) != 2 {
		t.Fatalf("expected 2 goals, got %d", len(goals))
	}
}

// --- Stats Tests ---

func TestGetWeeklyStats_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user and project
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createProjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "Math"}))
	var proj map[string]any
	parseJSON(t, createProjRR, &proj)
	projectID := proj["id"].(string)

	// Create a study log within the week
	studiedAt := time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC) // Tuesday in the week of 2024-01-01
	doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/study-logs", map[string]any{
		"projectId": projectID, "studiedAt": studiedAt.Format(time.RFC3339), "minutes": 90, "note": "study session",
	}))

	// Upsert a goal
	doRequest(handler, jsonRequest("PUT", "/v1/users/"+userID+"/goals/"+projectID, map[string]any{
		"targetMinutesPerWeek": 200, "startDate": "2024-01-01",
	}))

	// Get weekly stats
	statsReq := jsonRequest("GET", "/v1/users/"+userID+"/stats/weekly?weekStart=2024-01-01", nil)
	statsRR := doRequest(handler, statsReq)

	if statsRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusOK, statsRR.Code, statsRR.Body.String())
	}

	var statsResp map[string]any
	parseJSON(t, statsRR, &statsResp)

	if statsResp["weekStart"] != "2024-01-01" {
		t.Errorf("expected weekStart '2024-01-01', got '%v'", statsResp["weekStart"])
	}
	if int(statsResp["totalMinutes"].(float64)) != 90 {
		t.Errorf("expected totalMinutes 90, got %v", statsResp["totalMinutes"])
	}

	projects, ok := statsResp["projects"].([]any)
	if !ok {
		t.Fatalf("expected projects to be a list")
	}
	if len(projects) != 1 {
		t.Fatalf("expected 1 project, got %d", len(projects))
	}

	mathStats := projects[0].(map[string]any)
	if mathStats["projectName"] != "Math" {
		t.Errorf("expected projectName 'Math', got '%v'", mathStats["projectName"])
	}
	if int(mathStats["totalMinutes"].(float64)) != 90 {
		t.Errorf("expected totalMinutes 90 for Math, got %v", mathStats["totalMinutes"])
	}
	if int(mathStats["targetMinutesPerWeek"].(float64)) != 200 {
		t.Errorf("expected targetMinutesPerWeek 200, got %v", mathStats["targetMinutesPerWeek"])
	}
	// Achievement rate: 90/200 * 100 = 45.0
	if mathStats["achievementRate"].(float64) != 45.0 {
		t.Errorf("expected achievementRate 45.0, got %v", mathStats["achievementRate"])
	}
}

func TestGetWeeklyStats_MissingWeekStart(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create a user first so we don't get 404 for user
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	// weekStart is required:"true" in the handler, Huma will return 422 for missing required query param
	statsReq := jsonRequest("GET", "/v1/users/"+userID+"/stats/weekly", nil)
	statsRR := doRequest(handler, statsReq)

	// Huma returns 422 for missing required query parameters
	if statsRR.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusUnprocessableEntity, statsRR.Code, statsRR.Body.String())
	}
}

func TestGetWeeklyStats_InvalidWeekStartFormat(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	// Invalid date format
	statsReq := jsonRequest("GET", "/v1/users/"+userID+"/stats/weekly?weekStart=not-a-date", nil)
	statsRR := doRequest(handler, statsReq)

	if statsRR.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusBadRequest, statsRR.Code, statsRR.Body.String())
	}
}

// --- Additional edge case tests ---

func TestCreateProject_UserNotFound(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	req := jsonRequest("POST", "/v1/users/nonexistent-user/projects", map[string]string{"name": "Math"})
	rr := doRequest(handler, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}

func TestUpdateProject_NotFound(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	req := jsonRequest("PUT", "/v1/projects/nonexistent-id", map[string]string{"name": "Updated"})
	rr := doRequest(handler, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}

func TestDeleteProject_NotFound(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	req := jsonRequest("DELETE", "/v1/projects/nonexistent-id", nil)
	rr := doRequest(handler, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}

func TestDeleteStudyLog_NotFound(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	req := jsonRequest("DELETE", "/v1/study-logs/nonexistent-id", nil)
	rr := doRequest(handler, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}

func TestUpsertGoal_UserNotFound(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	body := map[string]any{
		"targetMinutesPerWeek": 300,
		"startDate":            "2024-01-01",
	}
	req := jsonRequest("PUT", "/v1/users/nonexistent/goals/some-project", body)
	rr := doRequest(handler, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}

func TestListGoals_UserNotFound(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	req := jsonRequest("GET", "/v1/users/nonexistent/goals", nil)
	rr := doRequest(handler, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}

// --- Note Tests ---

func TestCreateNote_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user and project
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createProjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "Math"}))
	var proj map[string]any
	parseJSON(t, createProjRR, &proj)
	projectID := proj["id"].(string)

	// Create note
	body := map[string]any{
		"title":   "My Note",
		"content": "some content",
		"tags":    []string{"go", "api"},
	}
	createNoteReq := jsonRequest("POST", "/v1/users/"+userID+"/projects/"+projectID+"/notes", body)
	createNoteRR := doRequest(handler, createNoteReq)

	if createNoteRR.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusCreated, createNoteRR.Code, createNoteRR.Body.String())
	}

	var noteResp map[string]any
	parseJSON(t, createNoteRR, &noteResp)
	if noteResp["title"] != "My Note" {
		t.Errorf("expected title 'My Note', got '%v'", noteResp["title"])
	}
	if noteResp["content"] != "some content" {
		t.Errorf("expected content 'some content', got '%v'", noteResp["content"])
	}
	if noteResp["userId"] != userID {
		t.Errorf("expected userId '%s', got '%v'", userID, noteResp["userId"])
	}
	if noteResp["projectId"] != projectID {
		t.Errorf("expected projectId '%s', got '%v'", projectID, noteResp["projectId"])
	}
	tags, ok := noteResp["tags"].([]any)
	if !ok || len(tags) != 2 {
		t.Errorf("expected 2 tags, got %v", noteResp["tags"])
	}
}

func TestListNotes_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user and project
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createProjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "Math"}))
	var proj map[string]any
	parseJSON(t, createProjRR, &proj)
	projectID := proj["id"].(string)

	// Create two notes
	doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects/"+projectID+"/notes", map[string]any{
		"title": "Note 1", "content": "content 1",
	}))
	doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects/"+projectID+"/notes", map[string]any{
		"title": "Note 2", "content": "content 2",
	}))

	// List notes
	listReq := jsonRequest("GET", "/v1/users/"+userID+"/projects/"+projectID+"/notes", nil)
	listRR := doRequest(handler, listReq)

	if listRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusOK, listRR.Code, listRR.Body.String())
	}

	var notes []map[string]any
	parseJSON(t, listRR, &notes)
	if len(notes) != 2 {
		t.Fatalf("expected 2 notes, got %d", len(notes))
	}
}

func TestGetNote_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user, project, and note
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createProjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "Math"}))
	var proj map[string]any
	parseJSON(t, createProjRR, &proj)
	projectID := proj["id"].(string)

	createNoteRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects/"+projectID+"/notes", map[string]any{
		"title": "My Note", "content": "some content",
	}))
	var created map[string]any
	parseJSON(t, createNoteRR, &created)
	noteID := created["id"].(string)

	// Get note
	getReq := jsonRequest("GET", "/v1/notes/"+noteID, nil)
	getRR := doRequest(handler, getReq)

	if getRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusOK, getRR.Code, getRR.Body.String())
	}

	var noteResp map[string]any
	parseJSON(t, getRR, &noteResp)
	if noteResp["title"] != "My Note" {
		t.Errorf("expected title 'My Note', got '%v'", noteResp["title"])
	}
}

func TestGetNote_NotFound(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	req := jsonRequest("GET", "/v1/notes/nonexistent-id", nil)
	rr := doRequest(handler, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}

func TestUpdateNote_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user, project, and note
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createProjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "Math"}))
	var proj map[string]any
	parseJSON(t, createProjRR, &proj)
	projectID := proj["id"].(string)

	createNoteRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects/"+projectID+"/notes", map[string]any{
		"title": "Old Title", "content": "old content",
	}))
	var created map[string]any
	parseJSON(t, createNoteRR, &created)
	noteID := created["id"].(string)

	// Update note
	updateReq := jsonRequest("PUT", "/v1/notes/"+noteID, map[string]any{
		"title": "New Title", "content": "new content", "tags": []string{"updated"},
	})
	updateRR := doRequest(handler, updateReq)

	if updateRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusOK, updateRR.Code, updateRR.Body.String())
	}

	var updated map[string]any
	parseJSON(t, updateRR, &updated)
	if updated["title"] != "New Title" {
		t.Errorf("expected title 'New Title', got '%v'", updated["title"])
	}
	if updated["content"] != "new content" {
		t.Errorf("expected content 'new content', got '%v'", updated["content"])
	}
}

func TestUpdateNote_NotFound(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	req := jsonRequest("PUT", "/v1/notes/nonexistent-id", map[string]any{
		"title": "Title", "content": "content",
	})
	rr := doRequest(handler, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}

func TestDeleteNote_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user, project, and note
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createProjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects", map[string]string{"name": "Math"}))
	var proj map[string]any
	parseJSON(t, createProjRR, &proj)
	projectID := proj["id"].(string)

	createNoteRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/projects/"+projectID+"/notes", map[string]any{
		"title": "To Delete", "content": "content",
	}))
	var created map[string]any
	parseJSON(t, createNoteRR, &created)
	noteID := created["id"].(string)

	// Delete note
	deleteReq := jsonRequest("DELETE", "/v1/notes/"+noteID, nil)
	deleteRR := doRequest(handler, deleteReq)

	if deleteRR.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNoContent, deleteRR.Code, deleteRR.Body.String())
	}
}

func TestDeleteNote_NotFound(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	req := jsonRequest("DELETE", "/v1/notes/nonexistent-id", nil)
	rr := doRequest(handler, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}

func TestCreateNote_UserNotFound(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	body := map[string]any{"title": "Note", "content": "content"}
	req := jsonRequest("POST", "/v1/users/nonexistent/projects/some-project/notes", body)
	rr := doRequest(handler, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}
