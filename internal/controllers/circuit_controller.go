package controllers

import (
	"context"
	"fmt"

	"electric-circuit-web/server/internal/models"
	"electric-circuit-web/server/internal/services"
)

// CircuitController handles circuit-related business coordination
type CircuitController struct {
	circuitService  *services.FirebaseCircuitService
	templateService *services.TemplateService
}

// NewCircuitController creates a new circuit controller
func NewCircuitController(circuitService *services.FirebaseCircuitService, templateService *services.TemplateService) *CircuitController {
	return &CircuitController{
		circuitService:  circuitService,
		templateService: templateService,
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
	// Note: Basic validation (required fields) is done in Handler layer
	// Controller focuses on business logic coordination

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
	// Note: Basic validation (required fields) is done in Handler layer

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
	// Note: Basic validation (required fields) is done in Handler layer
	// Controller handles business logic validation and DTO mapping

	// Business validation: Circuit name length
	if len(req.Name) > 255 {
		return &CircuitResponse{
			Success: false,
			Error:   "Circuit name must not exceed 255 characters",
		}, nil
	}

	// Create circuit model (DTO to Domain model)
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
	// Note: Basic validation (required fields) is done in Handler layer

	// Business validation: Name length if provided
	if req.Name != "" && len(req.Name) > 255 {
		return &CircuitResponse{
			Success: false,
			Error:   "Circuit name must not exceed 255 characters",
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
	// Note: Basic validation (required fields) is done in Handler layer

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
	// Note: Basic validation (required fields) is done in Handler layer

	// Business validation: Circuit name length
	if len(req.Name) > 255 {
		return &CircuitResponse{
			Success: false,
			Error:   "Circuit name must not exceed 255 characters",
		}, nil
	}

	// Get template data from TemplateService
	templateData := c.templateService.GetTemplateData(req.TemplateID)
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
