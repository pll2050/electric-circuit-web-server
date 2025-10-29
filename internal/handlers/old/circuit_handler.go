package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"electric-circuit-web/server/internal/middleware"
	"electric-circuit-web/server/internal/models"
	"electric-circuit-web/server/internal/services"
)

// CircuitHandler handles circuit-related HTTP requests
type CircuitHandler struct {
	circuitService *services.FirebaseCircuitService
}

// NewCircuitHandler creates a new circuit handler
func NewCircuitHandler(circuitService *services.FirebaseCircuitService) *CircuitHandler {
	return &CircuitHandler{
		circuitService: circuitService,
	}
}

// GetCircuits retrieves circuits for a project
func (h *CircuitHandler) GetCircuits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		http.Error(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	// Get circuits for this project
	circuits, err := h.circuitService.GetProjectCircuits(projectID, user.UID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(circuits)
}

// GetCircuit retrieves a specific circuit
func (h *CircuitHandler) GetCircuit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	circuitID := r.URL.Query().Get("id")
	if circuitID == "" {
		http.Error(w, "Circuit ID is required", http.StatusBadRequest)
		return
	}

	circuit, err := h.circuitService.GetCircuit(circuitID, user.UID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(circuit)
}

// CreateCircuit creates a new circuit
func (h *CircuitHandler) CreateCircuit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var circuit models.CircuitFirestore
	if err := json.NewDecoder(r.Body).Decode(&circuit); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create circuit using service
	circuitID, err := h.circuitService.CreateCircuit(user.UID, &circuit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"id":      circuitID,
		"message": "Circuit created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateCircuit updates a circuit
func (h *CircuitHandler) UpdateCircuit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	circuitID := r.URL.Query().Get("id")
	if circuitID == "" {
		http.Error(w, "Circuit ID is required", http.StatusBadRequest)
		return
	}

	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.circuitService.UpdateCircuit(circuitID, user.UID, updateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Circuit updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteCircuit deletes a circuit
func (h *CircuitHandler) DeleteCircuit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	circuitID := r.URL.Query().Get("id")
	if circuitID == "" {
		http.Error(w, "Circuit ID is required", http.StatusBadRequest)
		return
	}

	err := h.circuitService.DeleteCircuit(circuitID, user.UID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Circuit deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DuplicateCircuit creates a copy of an existing circuit
func (h *CircuitHandler) DuplicateCircuit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	circuitID := r.URL.Query().Get("id")
	if circuitID == "" {
		http.Error(w, "Circuit ID is required", http.StatusBadRequest)
		return
	}

	var request struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get original circuit
	originalCircuit, err := h.circuitService.GetCircuit(circuitID, user.UID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Create new circuit with copied data
	newCircuit := &models.CircuitFirestore{
		Name:        request.Name,
		Description: originalCircuit.Description + " (Copy)",
		ProjectID:   originalCircuit.ProjectID,
		Data:        originalCircuit.Data,
		Version:     1,
		IsTemplate:  false,
		Tags:        originalCircuit.Tags,
	}

	newCircuitID, err := h.circuitService.CreateCircuit(user.UID, newCircuit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":      newCircuitID,
		"message": "Circuit duplicated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// SaveCircuitTemplate saves a circuit as a template
func (h *CircuitHandler) SaveCircuitTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	circuitID := r.URL.Query().Get("id")
	if circuitID == "" {
		http.Error(w, "Circuit ID is required", http.StatusBadRequest)
		return
	}

	var request struct {
		TemplateName string   `json:"templateName"`
		Description  string   `json:"description"`
		Tags         []string `json:"tags"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get original circuit
	originalCircuit, err := h.circuitService.GetCircuit(circuitID, user.UID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Create template circuit
	templateCircuit := &models.CircuitFirestore{
		Name:        request.TemplateName,
		Description: request.Description,
		ProjectID:   "templates", // Special project for templates
		Data:        originalCircuit.Data,
		Version:     1,
		IsTemplate:  true,
		Tags:        request.Tags,
	}

	templateID, err := h.circuitService.CreateCircuit(user.UID, templateCircuit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":      templateID,
		"message": "Circuit template saved successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetTemplates handles GET /api/circuits/templates - get available circuit templates
func (h *CircuitHandler) GetTemplates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Mock templates for now - in real implementation, these would come from a templates collection
	templates := []map[string]interface{}{
		{
			"id":          "template_basic_circuit",
			"name":        "Basic Circuit",
			"description": "A simple circuit with resistor, LED, and battery",
			"category":    "basic",
			"thumbnail":   "https://example.com/thumbnails/basic_circuit.png",
			"data": map[string]interface{}{
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
		},
		{
			"id":          "template_amplifier",
			"name":        "Op-Amp Circuit",
			"description": "Basic operational amplifier configuration",
			"category":    "analog",
			"thumbnail":   "https://example.com/thumbnails/opamp_circuit.png",
			"data": map[string]interface{}{
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
		},
	}

	response := map[string]interface{}{
		"templates":   templates,
		"count":       len(templates),
		"requestedBy": user.UID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateFromTemplate handles POST /api/circuits/create-from-template - create circuit from template
func (h *CircuitHandler) CreateFromTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get user from context
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req struct {
		TemplateID  string `json:"templateId"`
		ProjectID   string `json:"projectId"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.TemplateID == "" || req.ProjectID == "" || req.Name == "" {
		http.Error(w, "Template ID, project ID, and name are required", http.StatusBadRequest)
		return
	}

	// Get template data (mock for now)
	var templateData map[string]interface{}
	switch req.TemplateID {
	case "template_basic_circuit":
		templateData = map[string]interface{}{
			"components": []map[string]interface{}{
				{"id": "battery1", "type": "battery", "x": 100, "y": 100},
				{"id": "resistor1", "type": "resistor", "x": 200, "y": 100},
				{"id": "led1", "type": "led", "x": 300, "y": 100},
			},
			"connections": []map[string]interface{}{
				{"from": "battery1", "to": "resistor1"},
				{"from": "resistor1", "to": "led1"},
			},
		}
	case "template_amplifier":
		templateData = map[string]interface{}{
			"components": []map[string]interface{}{
				{"id": "opamp1", "type": "opamp", "x": 200, "y": 150},
				{"id": "r1", "type": "resistor", "x": 100, "y": 100},
				{"id": "r2", "type": "resistor", "x": 100, "y": 200},
			},
			"connections": []map[string]interface{}{
				{"from": "r1", "to": "opamp1"},
				{"from": "r2", "to": "opamp1"},
			},
		}
	default:
		http.Error(w, "Template not found", http.StatusNotFound)
		return
	}

	// Create circuit from template
	circuit := &models.CircuitFirestore{
		Name:        req.Name,
		Description: req.Description,
		ProjectID:   req.ProjectID,
		Data:        templateData,
		UserID:      user.UID,
		Version:     1,
		IsTemplate:  false,
		Tags:        []string{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	circuitID, err := h.circuitService.CreateCircuit(user.UID, circuit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the ID for response
	circuit.ID = circuitID

	response := map[string]interface{}{
		"message":    "Circuit created from template successfully",
		"circuit":    circuit,
		"templateId": req.TemplateID,
		"createdBy":  user.UID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
