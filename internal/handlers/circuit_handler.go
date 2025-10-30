package handlers

import (
	"encoding/json"
	"net/http"

	"electric-circuit-web/server/internal/controllers"
)

// CircuitHandler handles circuit-related HTTP requests
type CircuitHandler struct {
	circuitController controllers.CircuitControllerInterface
}

// NewCircuitHandler creates a new circuit handler
func NewCircuitHandler(circuitController controllers.CircuitControllerInterface) *CircuitHandler {
	return &CircuitHandler{
		circuitController: circuitController,
	}
}

// getUserIDFromContext extracts user ID from request context
// This would typically be set by authentication middleware
func (h *CircuitHandler) getUserIDFromContext(r *http.Request) string {
	// For now, return a placeholder - this should be replaced with actual auth context
	return r.Header.Get("X-User-ID")
}

// writeError writes an error response
func (h *CircuitHandler) writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// writeJSON writes a JSON response
func (h *CircuitHandler) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// HandleGetCircuits handles getting circuits for a project
func (h *CircuitHandler) HandleGetCircuits(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		h.writeError(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	req := &controllers.CircuitRequest{
		UserID:    userID,
		ProjectID: projectID,
	}

	response, err := h.circuitController.GetProjectCircuits(r.Context(), req)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !response.Success {
		h.writeError(w, response.Error, http.StatusBadRequest)
		return
	}

	h.writeJSON(w, response)
}

// HandleCreateCircuit handles circuit creation
func (h *CircuitHandler) HandleCreateCircuit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req controllers.CircuitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	req.UserID = userID

	// HTTP-level validation: Check required fields
	if req.Name == "" {
		h.writeError(w, "Circuit name is required", http.StatusBadRequest)
		return
	}
	if req.ProjectID == "" {
		h.writeError(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.circuitController.CreateCircuit(r.Context(), &req)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !response.Success {
		h.writeError(w, response.Error, http.StatusBadRequest)
		return
	}

	h.writeJSON(w, response)
}

// HandleGetCircuit handles getting a specific circuit
func (h *CircuitHandler) HandleGetCircuit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	circuitID := r.URL.Query().Get("circuitId")
	if circuitID == "" {
		h.writeError(w, "Circuit ID is required", http.StatusBadRequest)
		return
	}

	req := &controllers.CircuitRequest{
		UserID:    userID,
		CircuitID: circuitID,
	}

	response, err := h.circuitController.GetCircuit(r.Context(), req)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !response.Success {
		h.writeError(w, response.Error, http.StatusBadRequest)
		return
	}

	h.writeJSON(w, response)
}

// HandleUpdateCircuit handles circuit updates
func (h *CircuitHandler) HandleUpdateCircuit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req controllers.CircuitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	req.UserID = userID

	response, err := h.circuitController.UpdateCircuit(r.Context(), &req)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !response.Success {
		h.writeError(w, response.Error, http.StatusBadRequest)
		return
	}

	h.writeJSON(w, response)
}

// HandleDeleteCircuit handles circuit deletion
func (h *CircuitHandler) HandleDeleteCircuit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	circuitID := r.URL.Query().Get("circuitId")
	if circuitID == "" {
		h.writeError(w, "Circuit ID is required", http.StatusBadRequest)
		return
	}

	req := &controllers.CircuitRequest{
		UserID:    userID,
		CircuitID: circuitID,
	}

	response, err := h.circuitController.DeleteCircuit(r.Context(), req)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !response.Success {
		h.writeError(w, response.Error, http.StatusBadRequest)
		return
	}

	h.writeJSON(w, response)
}

// HandleGetTemplates handles getting circuit templates
func (h *CircuitHandler) HandleGetTemplates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// For now, return a simple templates response
	// This should be implemented in the circuit controller
	h.writeJSON(w, map[string]interface{}{
		"success":   true,
		"message":   "Templates feature not yet implemented",
		"templates": []interface{}{},
	})
}

// HandleCreateFromTemplate handles creating a circuit from a template
func (h *CircuitHandler) HandleCreateFromTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req controllers.CircuitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	req.UserID = userID

	// HTTP-level validation: Check required fields
	if req.TemplateID == "" {
		h.writeError(w, "Template ID is required", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		h.writeError(w, "Circuit name is required", http.StatusBadRequest)
		return
	}
	if req.ProjectID == "" {
		h.writeError(w, "Project ID is required", http.StatusBadRequest)
		return
	}

	response, err := h.circuitController.CreateFromTemplate(r.Context(), &req)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !response.Success {
		h.writeError(w, response.Error, http.StatusBadRequest)
		return
	}

	h.writeJSON(w, response)
}
