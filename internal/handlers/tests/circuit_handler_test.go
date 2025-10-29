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
	"electric-circuit-web/server/internal/models"
)

// MockCircuitController implements the CircuitControllerInterface for testing
type MockCircuitController struct {
	getProjectCircuitsResponse *controllers.CircuitResponse
	getProjectCircuitsError    error
	createCircuitResponse      *controllers.CircuitResponse
	createCircuitError         error
	getCircuitResponse         *controllers.CircuitResponse
	getCircuitError            error
}

func (m *MockCircuitController) GetProjectCircuits(ctx context.Context, req *controllers.CircuitRequest) (*controllers.CircuitResponse, error) {
	return m.getProjectCircuitsResponse, m.getProjectCircuitsError
}

func (m *MockCircuitController) CreateCircuit(ctx context.Context, req *controllers.CircuitRequest) (*controllers.CircuitResponse, error) {
	return m.createCircuitResponse, m.createCircuitError
}

func (m *MockCircuitController) GetCircuit(ctx context.Context, req *controllers.CircuitRequest) (*controllers.CircuitResponse, error) {
	return m.getCircuitResponse, m.getCircuitError
}

func (m *MockCircuitController) UpdateCircuit(ctx context.Context, req *controllers.CircuitRequest) (*controllers.CircuitResponse, error) {
	return &controllers.CircuitResponse{Success: true, Message: "Circuit updated"}, nil
}

func (m *MockCircuitController) DeleteCircuit(ctx context.Context, req *controllers.CircuitRequest) (*controllers.CircuitResponse, error) {
	return &controllers.CircuitResponse{Success: true, Message: "Circuit deleted"}, nil
}

func (m *MockCircuitController) CreateFromTemplate(ctx context.Context, req *controllers.CircuitRequest) (*controllers.CircuitResponse, error) {
	return &controllers.CircuitResponse{Success: true, Message: "Circuit created from template"}, nil
}

func TestCircuitHandler_HandleGetCircuits_Success(t *testing.T) {
	// Arrange
	mockController := &MockCircuitController{
		getProjectCircuitsResponse: &controllers.CircuitResponse{
			Success: true,
			Message: "Circuits retrieved successfully",
			Circuits: []models.CircuitFirestore{
				{
					ID:          "circuit-1",
					Name:        "Test Circuit 1",
					Description: "Test Description 1",
				},
				{
					ID:          "circuit-2",
					Name:        "Test Circuit 2",
					Description: "Test Description 2",
				},
			},
		},
		getProjectCircuitsError: nil,
	}

	handler := handlers.NewCircuitHandler(mockController)

	req := httptest.NewRequest("GET", "/api/circuits?projectId=test-project", nil)
	req.Header.Set("X-User-ID", "test-user")
	w := httptest.NewRecorder()

	// Act
	handler.HandleGetCircuits(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response controllers.CircuitResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success to be true, got %v", response.Success)
	}

	if len(response.Circuits) != 2 {
		t.Errorf("Expected 2 circuits, got %d", len(response.Circuits))
	}
}

func TestCircuitHandler_HandleGetCircuits_MissingUserID(t *testing.T) {
	// Arrange
	mockController := &MockCircuitController{}
	handler := handlers.NewCircuitHandler(mockController)

	req := httptest.NewRequest("GET", "/api/circuits?projectId=test-project", nil)
	// No X-User-ID header
	w := httptest.NewRecorder()

	// Act
	handler.HandleGetCircuits(w, req)

	// Assert
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestCircuitHandler_HandleGetCircuits_MissingProjectID(t *testing.T) {
	// Arrange
	mockController := &MockCircuitController{}
	handler := handlers.NewCircuitHandler(mockController)

	req := httptest.NewRequest("GET", "/api/circuits", nil) // No projectId parameter
	req.Header.Set("X-User-ID", "test-user")
	w := httptest.NewRecorder()

	// Act
	handler.HandleGetCircuits(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCircuitHandler_HandleGetCircuits_InvalidMethod(t *testing.T) {
	// Arrange
	mockController := &MockCircuitController{}
	handler := handlers.NewCircuitHandler(mockController)

	req := httptest.NewRequest("POST", "/api/circuits?projectId=test-project", nil)
	req.Header.Set("X-User-ID", "test-user")
	w := httptest.NewRecorder()

	// Act
	handler.HandleGetCircuits(w, req)

	// Assert
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestCircuitHandler_HandleCreateCircuit_Success(t *testing.T) {
	// Arrange
	mockController := &MockCircuitController{
		createCircuitResponse: &controllers.CircuitResponse{
			Success: true,
			Message: "Circuit created successfully",
			Circuit: &models.CircuitFirestore{
				ID:          "new-circuit-id",
				Name:        "New Circuit",
				Description: "New Circuit Description",
			},
		},
		createCircuitError: nil,
	}

	handler := handlers.NewCircuitHandler(mockController)

	requestBody := map[string]interface{}{
		"project_id":  "test-project",
		"name":        "New Circuit",
		"description": "New Circuit Description",
		"data":        map[string]interface{}{"elements": []interface{}{}},
	}
	jsonBody, _ := json.Marshal(requestBody)

	req := httptest.NewRequest("POST", "/api/circuits/create", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "test-user")
	w := httptest.NewRecorder()

	// Act
	handler.HandleCreateCircuit(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response controllers.CircuitResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success to be true, got %v", response.Success)
	}

	if response.Circuit == nil {
		t.Error("Expected circuit in response, got nil")
	}
}

func TestCircuitHandler_HandleCreateCircuit_InvalidJSON(t *testing.T) {
	// Arrange
	mockController := &MockCircuitController{}
	handler := handlers.NewCircuitHandler(mockController)

	req := httptest.NewRequest("POST", "/api/circuits/create", bytes.NewReader([]byte("invalid-json")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "test-user")
	w := httptest.NewRecorder()

	// Act
	handler.HandleCreateCircuit(w, req)

	// Assert
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCircuitHandler_HandleGetCircuit_Success(t *testing.T) {
	// Arrange
	mockController := &MockCircuitController{
		getCircuitResponse: &controllers.CircuitResponse{
			Success: true,
			Message: "Circuit retrieved successfully",
			Circuit: &models.CircuitFirestore{
				ID:          "circuit-1",
				Name:        "Test Circuit",
				Description: "Test Description",
			},
		},
		getCircuitError: nil,
	}

	handler := handlers.NewCircuitHandler(mockController)

	req := httptest.NewRequest("GET", "/api/circuits/get?circuitId=circuit-1", nil)
	req.Header.Set("X-User-ID", "test-user")
	w := httptest.NewRecorder()

	// Act
	handler.HandleGetCircuit(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response controllers.CircuitResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success to be true, got %v", response.Success)
	}

	if response.Circuit == nil {
		t.Error("Expected circuit in response, got nil")
	}
}