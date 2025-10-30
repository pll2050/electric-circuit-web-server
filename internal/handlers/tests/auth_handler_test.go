package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"electric-circuit-web/server/internal/controllers"
	"electric-circuit-web/server/internal/handlers"
)

// MockAuthController implements the AuthControllerInterface for testing
type MockAuthController struct {
	verifyTokenResponse *controllers.AuthResponse
	verifyTokenError    error
	createUserResponse  *controllers.AuthResponse
	createUserError     error
}

func (m *MockAuthController) VerifyToken(ctx context.Context, req *controllers.AuthRequest) (*controllers.AuthResponse, error) {
	return m.verifyTokenResponse, m.verifyTokenError
}

func (m *MockAuthController) CreateUser(ctx context.Context, req *controllers.AuthRequest) (*controllers.AuthResponse, error) {
	return m.createUserResponse, m.createUserError
}

func (m *MockAuthController) GetUser(ctx context.Context, req *controllers.AuthRequest) (*controllers.AuthResponse, error) {
	return &controllers.AuthResponse{Success: true, Message: "User found"}, nil
}

func (m *MockAuthController) UpdateUser(ctx context.Context, req *controllers.AuthRequest) (*controllers.AuthResponse, error) {
	return &controllers.AuthResponse{Success: true, Message: "User updated"}, nil
}

func (m *MockAuthController) DeleteUser(ctx context.Context, req *controllers.AuthRequest) (*controllers.AuthResponse, error) {
	return &controllers.AuthResponse{Success: true, Message: "User deleted"}, nil
}

func (m *MockAuthController) SetCustomClaims(ctx context.Context, req *controllers.AuthRequest) (*controllers.AuthResponse, error) {
	return &controllers.AuthResponse{Success: true, Message: "Claims set"}, nil
}

func (m *MockAuthController) Register(ctx context.Context, req *controllers.AuthRequest) (*controllers.AuthResponse, error) {
	return &controllers.AuthResponse{Success: true, Message: "User registered"}, nil
}

func (m *MockAuthController) ListUsers(ctx context.Context) ([]map[string]interface{}, error) {
	return []map[string]interface{}{
		{
			"uid":         "test-uid-1",
			"email":       "test1@example.com",
			"displayName": "Test User 1",
		},
		{
			"uid":         "test-uid-2",
			"email":       "test2@example.com",
			"displayName": "Test User 2",
		},
	}, nil
}

func TestAuthHandler_HandleVerifyToken_Success(t *testing.T) {
	// Arrange
	mockController := &MockAuthController{
		verifyTokenResponse: &controllers.AuthResponse{
			Success: true,
			Message: "Token verified successfully",
			Token:   "valid-token",
		},
		verifyTokenError: nil,
	}

	handler := handlers.NewAuthHandler(mockController)

	requestBody := map[string]string{
		"token": "test-token",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/api/auth/verify", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.HandleVerifyToken(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response controllers.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success to be true, got %v", response.Success)
	}

	if response.Message != "Token verified successfully" {
		t.Errorf("Expected message 'Token verified successfully', got '%s'", response.Message)
	}
}

func TestAuthHandler_HandleVerifyToken_InvalidMethod(t *testing.T) {
	// Arrange
	mockController := &MockAuthController{}
	handler := handlers.NewAuthHandler(mockController)

	req := httptest.NewRequest("GET", "/api/auth/verify", nil)
	w := httptest.NewRecorder()

	// Act
	handler.HandleVerifyToken(w, req)

	// Assert
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestAuthHandler_HandleVerifyToken_InvalidJSON(t *testing.T) {
	// Arrange
	mockController := &MockAuthController{}
	handler := handlers.NewAuthHandler(mockController)

	req := httptest.NewRequest("POST", "/api/auth/verify", bytes.NewReader([]byte("invalid-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.HandleVerifyToken(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAuthHandler_HandleCreateUser_Success(t *testing.T) {
	// Arrange
	mockController := &MockAuthController{
		createUserResponse: &controllers.AuthResponse{
			Success: true,
			Message: "User created successfully",
		},
		createUserError: nil,
	}

	handler := handlers.NewAuthHandler(mockController)

	requestBody := map[string]string{
		"email":        "test@example.com",
		"password":     "password123",
		"display_name": "Test User",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/api/auth/create-user", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.HandleCreateUser(w, req)

	// Assert
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response controllers.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success to be true, got %v", response.Success)
	}
}

func TestAuthHandler_HandleGetUser_Success(t *testing.T) {
	// Arrange
	mockController := &MockAuthController{}
	handler := handlers.NewAuthHandler(mockController)

	req := httptest.NewRequest("GET", "/api/auth/get-user?uid=test-uid", nil)
	w := httptest.NewRecorder()

	// Act
	handler.HandleGetUser(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestAuthHandler_HandleGetUser_MissingUID(t *testing.T) {
	// Arrange
	mockController := &MockAuthController{}
	handler := handlers.NewAuthHandler(mockController)

	req := httptest.NewRequest("GET", "/api/auth/get-user", nil)
	w := httptest.NewRecorder()

	// Act
	handler.HandleGetUser(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}