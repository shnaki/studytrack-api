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

type mockSubjectRepository struct {
	subjects map[string]*domain.Subject
}

func newMockSubjectRepo() *mockSubjectRepository {
	return &mockSubjectRepository{subjects: make(map[string]*domain.Subject)}
}

func (m *mockSubjectRepository) Create(_ context.Context, s *domain.Subject) error {
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

func (m *mockSubjectRepository) Update(_ context.Context, s *domain.Subject) error {
	m.subjects[s.ID] = s
	return nil
}

func (m *mockSubjectRepository) Delete(_ context.Context, id string) error {
	delete(m.subjects, id)
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
		if filter.SubjectID != nil && l.SubjectID != *filter.SubjectID {
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
	m.goals[g.UserID+":"+g.SubjectID] = g
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

// --- Helpers ---

func setupRouter(t *testing.T) (http.Handler, *mockUserRepository, *mockSubjectRepository, *mockStudyLogRepository, *mockGoalRepository) {
	t.Helper()
	userRepo := newMockUserRepo()
	subjectRepo := newMockSubjectRepo()
	studyLogRepo := newMockStudyLogRepo()
	goalRepo := newMockGoalRepo()

	usecases := &controller.Usecases{
		User:     usecase.NewUserUsecase(userRepo),
		Subject:  usecase.NewSubjectUsecase(subjectRepo, userRepo),
		StudyLog: usecase.NewStudyLogUsecase(studyLogRepo, userRepo, subjectRepo),
		Goal:     usecase.NewGoalUsecase(goalRepo, userRepo, subjectRepo),
		Stats:    usecase.NewStatsUsecase(studyLogRepo, goalRepo, subjectRepo),
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	router := controller.NewRouter(usecases, []string{"*"}, logger)
	return router, userRepo, subjectRepo, studyLogRepo, goalRepo
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

// --- Subject Tests ---

func TestCreateSubject_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user first
	createUserReq := jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"})
	createUserRR := doRequest(handler, createUserReq)
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	// Create subject
	createSubjReq := jsonRequest("POST", "/v1/users/"+userID+"/subjects", map[string]string{"name": "Mathematics"})
	createSubjRR := doRequest(handler, createSubjReq)

	if createSubjRR.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusCreated, createSubjRR.Code, createSubjRR.Body.String())
	}

	var subj map[string]any
	parseJSON(t, createSubjRR, &subj)
	if subj["name"] != "Mathematics" {
		t.Errorf("expected name 'Mathematics', got '%v'", subj["name"])
	}
	if subj["userId"] != userID {
		t.Errorf("expected userId '%s', got '%v'", userID, subj["userId"])
	}
	if subj["id"] == nil || subj["id"] == "" {
		t.Error("expected subject id to be set")
	}
}

func TestListSubjects_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user
	createUserReq := jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"})
	createUserRR := doRequest(handler, createUserReq)
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	// Create two subjects
	doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/subjects", map[string]string{"name": "Math"}))
	doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/subjects", map[string]string{"name": "English"}))

	// List subjects
	listReq := jsonRequest("GET", "/v1/users/"+userID+"/subjects", nil)
	listRR := doRequest(handler, listReq)

	if listRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusOK, listRR.Code, listRR.Body.String())
	}

	var subjects []map[string]any
	parseJSON(t, listRR, &subjects)
	if len(subjects) != 2 {
		t.Fatalf("expected 2 subjects, got %d", len(subjects))
	}
}

func TestUpdateSubject_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user and subject
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createSubjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/subjects", map[string]string{"name": "Math"}))
	var subj map[string]any
	parseJSON(t, createSubjRR, &subj)
	subjectID := subj["id"].(string)

	// Update subject
	updateReq := jsonRequest("PUT", "/v1/subjects/"+subjectID, map[string]string{"name": "Advanced Math"})
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

func TestDeleteSubject_Success(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	// Create user and subject
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createSubjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/subjects", map[string]string{"name": "Math"}))
	var subj map[string]any
	parseJSON(t, createSubjRR, &subj)
	subjectID := subj["id"].(string)

	// Delete subject
	deleteReq := jsonRequest("DELETE", "/v1/subjects/"+subjectID, nil)
	deleteRR := doRequest(handler, deleteReq)

	if deleteRR.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNoContent, deleteRR.Code, deleteRR.Body.String())
	}
}

func TestListSubjects_UserNotFound(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	req := jsonRequest("GET", "/v1/users/nonexistent-user/subjects", nil)
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

	// Create subject
	createSubjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/subjects", map[string]string{"name": "Math"}))
	var subj map[string]any
	parseJSON(t, createSubjRR, &subj)
	subjectID := subj["id"].(string)

	// Create study log
	studiedAt := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	body := map[string]any{
		"subjectId": subjectID,
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
	if logResp["subjectId"] != subjectID {
		t.Errorf("expected subjectId '%s', got '%v'", subjectID, logResp["subjectId"])
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

	// Create user and subject
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createSubjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/subjects", map[string]string{"name": "Math"}))
	var subj map[string]any
	parseJSON(t, createSubjRR, &subj)
	subjectID := subj["id"].(string)

	// Create two study logs
	studiedAt1 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	studiedAt2 := time.Date(2024, 1, 16, 10, 0, 0, 0, time.UTC)
	doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/study-logs", map[string]any{
		"subjectId": subjectID, "studiedAt": studiedAt1.Format(time.RFC3339), "minutes": 60, "note": "session 1",
	}))
	doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/study-logs", map[string]any{
		"subjectId": subjectID, "studiedAt": studiedAt2.Format(time.RFC3339), "minutes": 45, "note": "session 2",
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

	// Create user and subject
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createSubjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/subjects", map[string]string{"name": "Math"}))
	var subj map[string]any
	parseJSON(t, createSubjRR, &subj)
	subjectID := subj["id"].(string)

	// Create logs on different dates
	jan15 := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	feb10 := time.Date(2024, 2, 10, 10, 0, 0, 0, time.UTC)
	doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/study-logs", map[string]any{
		"subjectId": subjectID, "studiedAt": jan15.Format(time.RFC3339), "minutes": 60, "note": "january",
	}))
	doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/study-logs", map[string]any{
		"subjectId": subjectID, "studiedAt": feb10.Format(time.RFC3339), "minutes": 45, "note": "february",
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

	// Create user, subject, and study log
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createSubjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/subjects", map[string]string{"name": "Math"}))
	var subj map[string]any
	parseJSON(t, createSubjRR, &subj)
	subjectID := subj["id"].(string)

	studiedAt := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	createLogRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/study-logs", map[string]any{
		"subjectId": subjectID, "studiedAt": studiedAt.Format(time.RFC3339), "minutes": 60, "note": "to delete",
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

	// Create user and subject
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createSubjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/subjects", map[string]string{"name": "Math"}))
	var subj map[string]any
	parseJSON(t, createSubjRR, &subj)
	subjectID := subj["id"].(string)

	// Upsert goal
	body := map[string]any{
		"targetMinutesPerWeek": 300,
		"startDate":            "2024-01-01",
	}
	upsertReq := jsonRequest("PUT", "/v1/users/"+userID+"/goals/"+subjectID, body)
	upsertRR := doRequest(handler, upsertReq)

	if upsertRR.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusOK, upsertRR.Code, upsertRR.Body.String())
	}

	var goalResp map[string]any
	parseJSON(t, upsertRR, &goalResp)
	if goalResp["userId"] != userID {
		t.Errorf("expected userId '%s', got '%v'", userID, goalResp["userId"])
	}
	if goalResp["subjectId"] != subjectID {
		t.Errorf("expected subjectId '%s', got '%v'", subjectID, goalResp["subjectId"])
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

	// Create user and subject
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createSubjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/subjects", map[string]string{"name": "Math"}))
	var subj map[string]any
	parseJSON(t, createSubjRR, &subj)
	subjectID := subj["id"].(string)

	endDate := "2024-06-30"
	body := map[string]any{
		"targetMinutesPerWeek": 200,
		"startDate":            "2024-01-01",
		"endDate":              endDate,
	}
	upsertReq := jsonRequest("PUT", "/v1/users/"+userID+"/goals/"+subjectID, body)
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

	// Create user and two subjects
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createSubj1RR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/subjects", map[string]string{"name": "Math"}))
	var subj1 map[string]any
	parseJSON(t, createSubj1RR, &subj1)
	subjectID1 := subj1["id"].(string)

	createSubj2RR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/subjects", map[string]string{"name": "English"}))
	var subj2 map[string]any
	parseJSON(t, createSubj2RR, &subj2)
	subjectID2 := subj2["id"].(string)

	// Upsert goals for both subjects
	doRequest(handler, jsonRequest("PUT", "/v1/users/"+userID+"/goals/"+subjectID1, map[string]any{
		"targetMinutesPerWeek": 300, "startDate": "2024-01-01",
	}))
	doRequest(handler, jsonRequest("PUT", "/v1/users/"+userID+"/goals/"+subjectID2, map[string]any{
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

	// Create user and subject
	createUserRR := doRequest(handler, jsonRequest("POST", "/v1/users", map[string]string{"name": "Alice"}))
	var user map[string]any
	parseJSON(t, createUserRR, &user)
	userID := user["id"].(string)

	createSubjRR := doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/subjects", map[string]string{"name": "Math"}))
	var subj map[string]any
	parseJSON(t, createSubjRR, &subj)
	subjectID := subj["id"].(string)

	// Create a study log within the week
	studiedAt := time.Date(2024, 1, 2, 10, 0, 0, 0, time.UTC) // Tuesday in the week of 2024-01-01
	doRequest(handler, jsonRequest("POST", "/v1/users/"+userID+"/study-logs", map[string]any{
		"subjectId": subjectID, "studiedAt": studiedAt.Format(time.RFC3339), "minutes": 90, "note": "study session",
	}))

	// Upsert a goal
	doRequest(handler, jsonRequest("PUT", "/v1/users/"+userID+"/goals/"+subjectID, map[string]any{
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

	subjects, ok := statsResp["subjects"].([]any)
	if !ok {
		t.Fatalf("expected subjects to be a list")
	}
	if len(subjects) != 1 {
		t.Fatalf("expected 1 subject, got %d", len(subjects))
	}

	mathStats := subjects[0].(map[string]any)
	if mathStats["subjectName"] != "Math" {
		t.Errorf("expected subjectName 'Math', got '%v'", mathStats["subjectName"])
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

func TestCreateSubject_UserNotFound(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	req := jsonRequest("POST", "/v1/users/nonexistent-user/subjects", map[string]string{"name": "Math"})
	rr := doRequest(handler, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}

func TestUpdateSubject_NotFound(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	req := jsonRequest("PUT", "/v1/subjects/nonexistent-id", map[string]string{"name": "Updated"})
	rr := doRequest(handler, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d; body: %s", http.StatusNotFound, rr.Code, rr.Body.String())
	}
}

func TestDeleteSubject_NotFound(t *testing.T) {
	handler, _, _, _, _ := setupRouter(t)

	req := jsonRequest("DELETE", "/v1/subjects/nonexistent-id", nil)
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
	req := jsonRequest("PUT", "/v1/users/nonexistent/goals/some-subject", body)
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
