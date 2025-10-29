package controllers

import (
	"context"
	"fmt"

	"electric-circuit-web/server/internal/models"
	"electric-circuit-web/server/internal/services"
)

// CircuitController handles circuit-related business coordination
type CircuitController struct {
	circuitService *services.FirebaseCircuitService
}

// NewCircuitController creates a new circuit controller
func NewCircuitController(circuitService *services.FirebaseCircuitService) *CircuitController {
	return &CircuitController{
		circuitService: circuitService,
	}
}

// CircuitRequest represents a circuit operation request
type CircuitRequest struct {
	UserID      string                 `json:"user_id"`
	ProjectID   string                 `json:"project_id"`
	CircuitID   string                 `json:"circuit_id,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
	TemplateID  string                 `json:"template_id,omitempty"`
}

// CircuitResponse represents a circuit operation response
type CircuitResponse struct {
	Success  bool                      `json:"success"`
	Message  string                    `json:"message"`
	Circuit  *models.CircuitFirestore  `json:"circuit,omitempty"`
	Circuits []models.CircuitFirestore `json:"circuits,omitempty"`
	Error    string                    `json:"error,omitempty"`
}

// GetProjectCircuits retrieves circuits for a project
func (c *CircuitController) GetProjectCircuits(ctx context.Context, req *CircuitRequest) (*CircuitResponse, error) {
	// Validate request
	if req.ProjectID == "" {
		return &CircuitResponse{
			Success: false,
			Error:   "Project ID is required",
		}, nil
	}

	if req.UserID == "" {
		return &CircuitResponse{
			Success: false,
			Error:   "User ID is required",
		}, nil
	}

	// Call service
	circuits, err := c.circuitService.GetProjectCircuits(req.ProjectID, req.UserID)
	if err != nil {
		return &CircuitResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &CircuitResponse{
		Success:  true,
		Message:  "Circuits retrieved successfully",
		Circuits: circuits,
	}, nil
}

// GetCircuit retrieves a specific circuit
func (c *CircuitController) GetCircuit(ctx context.Context, req *CircuitRequest) (*CircuitResponse, error) {
	// Validate request
	if req.CircuitID == "" {
		return &CircuitResponse{
			Success: false,
			Error:   "Circuit ID is required",
		}, nil
	}

	if req.UserID == "" {
		return &CircuitResponse{
			Success: false,
			Error:   "User ID is required",
		}, nil
	}

	// Call service
	circuit, err := c.circuitService.GetCircuit(req.CircuitID, req.UserID)
	if err != nil {
		return &CircuitResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &CircuitResponse{
		Success: true,
		Message: "Circuit retrieved successfully",
		Circuit: circuit,
	}, nil
}

// CreateCircuit creates a new circuit
func (c *CircuitController) CreateCircuit(ctx context.Context, req *CircuitRequest) (*CircuitResponse, error) {
	// Validate request
	if req.Name == "" {
		return &CircuitResponse{
			Success: false,
			Error:   "Circuit name is required",
		}, nil
	}

	if req.ProjectID == "" {
		return &CircuitResponse{
			Success: false,
			Error:   "Project ID is required",
		}, nil
	}

	if req.UserID == "" {
		return &CircuitResponse{
			Success: false,
			Error:   "User ID is required",
		}, nil
	}

	// Create circuit model
	circuit := &models.CircuitFirestore{
		Name:        req.Name,
		Description: req.Description,
		ProjectID:   req.ProjectID,
		UserID:      req.UserID,
		Data:        req.Data,
		Version:     1,
		IsTemplate:  false,
		Tags:        []string{},
	}

	// Call service
	circuitID, err := c.circuitService.CreateCircuit(req.UserID, circuit)
	if err != nil {
		return &CircuitResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// Set ID for response
	circuit.ID = circuitID

	return &CircuitResponse{
		Success: true,
		Message: "Circuit created successfully",
		Circuit: circuit,
	}, nil
}

// UpdateCircuit updates an existing circuit
func (c *CircuitController) UpdateCircuit(ctx context.Context, req *CircuitRequest) (*CircuitResponse, error) {
	// Validate request
	if req.CircuitID == "" {
		return &CircuitResponse{
			Success: false,
			Error:   "Circuit ID is required",
		}, nil
	}

	if req.UserID == "" {
		return &CircuitResponse{
			Success: false,
			Error:   "User ID is required",
		}, nil
	}

	// Prepare updates
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Data != nil {
		updates["data"] = req.Data
	}

	// Call service
	err := c.circuitService.UpdateCircuit(req.CircuitID, req.UserID, updates)
	if err != nil {
		return &CircuitResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &CircuitResponse{
		Success: true,
		Message: "Circuit updated successfully",
	}, nil
}

// DeleteCircuit deletes a circuit
func (c *CircuitController) DeleteCircuit(ctx context.Context, req *CircuitRequest) (*CircuitResponse, error) {
	// Validate request
	if req.CircuitID == "" {
		return &CircuitResponse{
			Success: false,
			Error:   "Circuit ID is required",
		}, nil
	}

	if req.UserID == "" {
		return &CircuitResponse{
			Success: false,
			Error:   "User ID is required",
		}, nil
	}

	// Call service
	err := c.circuitService.DeleteCircuit(req.CircuitID, req.UserID)
	if err != nil {
		return &CircuitResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &CircuitResponse{
		Success: true,
		Message: "Circuit deleted successfully",
	}, nil
}

// CreateFromTemplate creates a circuit from template
func (c *CircuitController) CreateFromTemplate(ctx context.Context, req *CircuitRequest) (*CircuitResponse, error) {
	// Validate request
	if req.TemplateID == "" {
		return &CircuitResponse{
			Success: false,
			Error:   "Template ID is required",
		}, nil
	}

	if req.Name == "" {
		return &CircuitResponse{
			Success: false,
			Error:   "Circuit name is required",
		}, nil
	}

	if req.ProjectID == "" {
		return &CircuitResponse{
			Success: false,
			Error:   "Project ID is required",
		}, nil
	}

	if req.UserID == "" {
		return &CircuitResponse{
			Success: false,
			Error:   "User ID is required",
		}, nil
	}

	// Get template data (this would be from a template service in the future)
	templateData := c.getTemplateData(req.TemplateID)
	if templateData == nil {
		return &CircuitResponse{
			Success: false,
			Error:   "Template not found",
		}, fmt.Errorf("template not found: %s", req.TemplateID)
	}

	// Create circuit from template
	circuit := &models.CircuitFirestore{
		Name:        req.Name,
		Description: req.Description,
		ProjectID:   req.ProjectID,
		UserID:      req.UserID,
		Data:        templateData,
		Version:     1,
		IsTemplate:  false,
		Tags:        []string{},
	}

	// Call service
	circuitID, err := c.circuitService.CreateCircuit(req.UserID, circuit)
	if err != nil {
		return &CircuitResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// Set ID for response
	circuit.ID = circuitID

	return &CircuitResponse{
		Success: true,
		Message: fmt.Sprintf("Circuit created from template %s successfully", req.TemplateID),
		Circuit: circuit,
	}, nil
}

// getTemplateData returns template data (mock implementation)
func (c *CircuitController) getTemplateData(templateID string) map[string]interface{} {
	templates := map[string]map[string]interface{}{
		"template_basic_circuit": {
			"components": []map[string]interface{}{
				{"id": "battery1", "type": "battery", "x": 100, "y": 100},
				{"id": "resistor1", "type": "resistor", "x": 200, "y": 100},
				{"id": "led1", "type": "led", "x": 300, "y": 100},
			},
			"connections": []map[string]interface{}{
				{"from": "battery1", "to": "resistor1"},
				{"from": "resistor1", "to": "led1"},
			},
		},
		"template_amplifier": {
			"components": []map[string]interface{}{
				{"id": "opamp1", "type": "opamp", "x": 200, "y": 150},
				{"id": "r1", "type": "resistor", "x": 100, "y": 100},
				{"id": "r2", "type": "resistor", "x": 100, "y": 200},
			},
			"connections": []map[string]interface{}{
				{"from": "r1", "to": "opamp1"},
				{"from": "r2", "to": "opamp1"},
			},
		},
	}

	return templates[templateID]
}
