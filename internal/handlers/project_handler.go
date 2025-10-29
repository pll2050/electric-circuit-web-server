package handlers

import (
	"encoding/json"
	"net/http"

	"electric-circuit-web/server/internal/controllers"
)

// ProjectHandler handles project-related HTTP requests
type ProjectHandler struct {
	projectController *controllers.ProjectController
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(projectController *controllers.ProjectController) *ProjectHandler {
	return &ProjectHandler{
		projectController: projectController,
	}
}

// getUserIDFromContext extracts user ID from request context
// This would typically be set by authentication middleware
func (h *ProjectHandler) getUserIDFromContext(r *http.Request) string {
	// For now, return a placeholder - this should be replaced with actual auth context
	return r.Header.Get("X-User-ID")
}

// writeError writes an error response
func (h *ProjectHandler) writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// writeJSON writes a JSON response
func (h *ProjectHandler) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// HandleGetProjects handles getting all projects for a user
func (h *ProjectHandler) HandleGetProjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	req := &controllers.ProjectRequest{
		UserID: userID,
	}

	response, err := h.projectController.GetUserProjects(r.Context(), req)
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

// HandleCreateProject handles project creation
func (h *ProjectHandler) HandleCreateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req controllers.ProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	req.UserID = userID

	response, err := h.projectController.CreateProject(r.Context(), &req)
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

// HandleGetProject handles getting a specific project
func (h *ProjectHandler) HandleGetProject(w http.ResponseWriter, r *http.Request) {
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

	req := &controllers.ProjectRequest{
		UserID:    userID,
		ProjectID: projectID,
	}

	response, err := h.projectController.GetProject(r.Context(), req)
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

// HandleUpdateProject handles project updates
func (h *ProjectHandler) HandleUpdateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req controllers.ProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	req.UserID = userID

	response, err := h.projectController.UpdateProject(r.Context(), &req)
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

// HandleDeleteProject handles project deletion
func (h *ProjectHandler) HandleDeleteProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
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

	req := &controllers.ProjectRequest{
		UserID:    userID,
		ProjectID: projectID,
	}

	response, err := h.projectController.DeleteProject(r.Context(), req)
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

// HandleDuplicateProject handles project duplication
func (h *ProjectHandler) HandleDuplicateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req controllers.ProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	req.UserID = userID

	response, err := h.projectController.DuplicateProject(r.Context(), &req)
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
