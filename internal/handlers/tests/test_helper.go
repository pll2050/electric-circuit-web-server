package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"electric-circuit-web/server/internal/controllers"
)

// TestHelper provides common testing utilities
type TestHelper struct {
	t *testing.T
}

// NewTestHelper creates a new test helper
func NewTestHelper(t *testing.T) *TestHelper {
	return &TestHelper{t: t}
}

// AssertStatusCode checks if the response has the expected status code
func (h *TestHelper) AssertStatusCode(w *httptest.ResponseRecorder, expected int) {
	if w.Code != expected {
		h.t.Errorf("Expected status code %d, got %d", expected, w.Code)
	}
}

// AssertContentType checks if the response has the expected content type
func (h *TestHelper) AssertContentType(w *httptest.ResponseRecorder, expected string) {
	contentType := w.Header().Get("Content-Type")
	if contentType != expected {
		h.t.Errorf("Expected content type '%s', got '%s'", expected, contentType)
	}
}

// CreateTestRequest creates a new HTTP test request with common headers
func (h *TestHelper) CreateTestRequest(method, url, userID string) *http.Request {
	req := httptest.NewRequest(method, url, nil)
	if userID != "" {
		req.Header.Set("X-User-ID", userID)
	}
	return req
}

// MockProjectController for completeness
type MockProjectController struct {
	getUserProjectsResponse *controllers.ProjectResponse
	getUserProjectsError    error
}

func (m *MockProjectController) GetUserProjects(ctx context.Context, req *controllers.ProjectRequest) (*controllers.ProjectResponse, error) {
	return m.getUserProjectsResponse, m.getUserProjectsError
}

func (m *MockProjectController) CreateProject(ctx context.Context, req *controllers.ProjectRequest) (*controllers.ProjectResponse, error) {
	return &controllers.ProjectResponse{Success: true, Message: "Project created"}, nil
}

func (m *MockProjectController) GetProject(ctx context.Context, req *controllers.ProjectRequest) (*controllers.ProjectResponse, error) {
	return &controllers.ProjectResponse{Success: true, Message: "Project retrieved"}, nil
}

func (m *MockProjectController) UpdateProject(ctx context.Context, req *controllers.ProjectRequest) (*controllers.ProjectResponse, error) {
	return &controllers.ProjectResponse{Success: true, Message: "Project updated"}, nil
}

func (m *MockProjectController) DeleteProject(ctx context.Context, req *controllers.ProjectRequest) (*controllers.ProjectResponse, error) {
	return &controllers.ProjectResponse{Success: true, Message: "Project deleted"}, nil
}

func (m *MockProjectController) DuplicateProject(ctx context.Context, req *controllers.ProjectRequest) (*controllers.ProjectResponse, error) {
	return &controllers.ProjectResponse{Success: true, Message: "Project duplicated"}, nil
}
