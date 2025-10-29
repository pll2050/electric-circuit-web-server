package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"electric-circuit-web/server/internal/controllers"
	"electric-circuit-web/server/internal/handlers"
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

// MockFullStack creates a complete mock stack for integration testing
type MockFullStack struct {
	AuthController    *MockAuthController
	CircuitController *MockCircuitController
	ProjectController *MockProjectController
	HTTPHandler       *handlers.HTTPHandler
}

// NewMockFullStack creates a new mock full stack for testing
func NewMockFullStack() *MockFullStack {
	authController := &MockAuthController{
		verifyTokenResponse: &controllers.AuthResponse{Success: true, Message: "Token verified"},
	}

	circuitController := &MockCircuitController{
		getProjectCircuitsResponse: &controllers.CircuitResponse{Success: true, Message: "Circuits retrieved"},
	}

	projectController := &MockProjectController{
		getUserProjectsResponse: &controllers.ProjectResponse{Success: true, Message: "Projects retrieved"},
	}

	// Note: For HTTPHandler, you would need to implement the constructor with all controllers
	// httpHandler := handlers.NewHTTPHandler(circuitController, projectController, authController, storageController)

	return &MockFullStack{
		AuthController:    authController,
		CircuitController: circuitController,
		ProjectController: projectController,
		// HTTPHandler:       httpHandler,
	}
}

// MockProjectController for completeness
type MockProjectController struct {
	getUserProjectsResponse *controllers.ProjectResponse
	getUserProjectsError    error
}

func (m *MockProjectController) GetUserProjects(ctx interface{}, req *controllers.ProjectRequest) (*controllers.ProjectResponse, error) {
	return m.getUserProjectsResponse, m.getUserProjectsError
}

func (m *MockProjectController) CreateProject(ctx interface{}, req *controllers.ProjectRequest) (*controllers.ProjectResponse, error) {
	return &controllers.ProjectResponse{Success: true, Message: "Project created"}, nil
}

func (m *MockProjectController) GetProject(ctx interface{}, req *controllers.ProjectRequest) (*controllers.ProjectResponse, error) {
	return &controllers.ProjectResponse{Success: true, Message: "Project retrieved"}, nil
}

func (m *MockProjectController) UpdateProject(ctx interface{}, req *controllers.ProjectRequest) (*controllers.ProjectResponse, error) {
	return &controllers.ProjectResponse{Success: true, Message: "Project updated"}, nil
}

func (m *MockProjectController) DeleteProject(ctx interface{}, req *controllers.ProjectRequest) (*controllers.ProjectResponse, error) {
	return &controllers.ProjectResponse{Success: true, Message: "Project deleted"}, nil
}

func (m *MockProjectController) DuplicateProject(ctx interface{}, req *controllers.ProjectRequest) (*controllers.ProjectResponse, error) {
	return &controllers.ProjectResponse{Success: true, Message: "Project duplicated"}, nil
}
