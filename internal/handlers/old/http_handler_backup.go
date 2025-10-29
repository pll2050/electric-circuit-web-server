package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"electric-circuit-web/server/internal/controllers"
	"electric-circuit-web/server/internal/middleware"
)

// HTTPHandler handles HTTP routing and protocol-specific concerns
type HTTPHandler struct {
	circuitController *controllers.CircuitController
	projectController *controllers.ProjectController
	authController    *controllers.AuthController
	storageController *controllers.StorageController
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(
	circuitController *controllers.CircuitController,
	projectController *controllers.ProjectController,
	authController *controllers.AuthController,
	storageController *controllers.StorageController,
) *HTTPHandler {
	return &HTTPHandler{
		circuitController: circuitController,
		projectController: projectController,
		authController:    authController,
		storageController: storageController,
	}
}

// SetupRoutes configures all HTTP routes
func (h *HTTPHandler) SetupRoutes() {
	// Health check (no auth required)
	http.HandleFunc("/api/health", middleware.CORS(h.healthCheck))

	// Auth routes
	http.HandleFunc("/api/auth/", middleware.CORS(h.routeAuth))

	// Protected routes (require auth)
	http.HandleFunc("/api/circuits", middleware.CORS(h.withAuth(h.routeCircuits)))
	http.HandleFunc("/api/projects", middleware.CORS(h.withAuth(h.routeProjects)))
	http.HandleFunc("/api/storage", middleware.CORS(h.withAuth(h.routeStorage)))
}

// withAuth wraps a handler with authentication middleware
func (h *HTTPHandler) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// This would be replaced with actual auth middleware
		// For now, we'll just call the next handler
		next.ServeHTTP(w, r)
	}
}

// routeAuth handles all authentication-related routes
func (h *HTTPHandler) routeAuth(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/auth")

	switch {
	case path == "/verify" && r.Method == "POST":
		h.handleAuthVerify(w, r)
	case path == "/create-user" && r.Method == "POST":
		h.handleAuthCreateUser(w, r)
	case path == "/get-user" && r.Method == "GET":
		h.handleAuthGetUser(w, r)
	case path == "/update-user" && r.Method == "PUT":
		h.handleAuthUpdateUser(w, r)
	case path == "/delete-user" && r.Method == "DELETE":
		h.handleAuthDeleteUser(w, r)
	case path == "/set-custom-claims" && r.Method == "POST":
		h.handleAuthSetCustomClaims(w, r)
	default:
		http.NotFound(w, r)
	}
}

// routeCircuits handles circuit-related HTTP routing
func (h *HTTPHandler) routeCircuits(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == "/api/circuits" && r.Method == "GET":
		h.handleCircuitsGet(w, r)
	case path == "/api/circuits/create" && r.Method == "POST":
		h.handleCircuitsCreate(w, r)
	case path == "/api/circuits/get" && r.Method == "GET":
		h.handleCircuitsGetOne(w, r)
	case path == "/api/circuits/update" && r.Method == "PUT":
		h.handleCircuitsUpdate(w, r)
	case path == "/api/circuits/delete" && r.Method == "DELETE":
		h.handleCircuitsDelete(w, r)
	case path == "/api/circuits/templates" && r.Method == "GET":
		h.handleCircuitsTemplates(w, r)
	case path == "/api/circuits/create-from-template" && r.Method == "POST":
		h.handleCircuitsCreateFromTemplate(w, r)
	default:
		http.NotFound(w, r)
	}
}

// routeProjects handles project-related HTTP routing
func (h *HTTPHandler) routeProjects(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == "/api/projects" && r.Method == "GET":
		h.handleProjectsGet(w, r)
	case path == "/api/projects/create" && r.Method == "POST":
		h.handleProjectsCreate(w, r)
	case path == "/api/projects/get" && r.Method == "GET":
		h.handleProjectsGetOne(w, r)
	case path == "/api/projects/update" && r.Method == "PUT":
		h.handleProjectsUpdate(w, r)
	case path == "/api/projects/delete" && r.Method == "DELETE":
		h.handleProjectsDelete(w, r)
	case path == "/api/projects/duplicate" && r.Method == "POST":
		h.handleProjectsDuplicate(w, r)
	default:
		http.NotFound(w, r)
	}
}

// routeStorage handles storage-related HTTP routing
func (h *HTTPHandler) routeStorage(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/storage")

	switch {
	case path == "/upload" && r.Method == "POST":
		h.handleStorageUpload(w, r)
	case path == "/url" && r.Method == "GET":
		h.handleStorageGetURL(w, r)
	case path == "/delete" && r.Method == "DELETE":
		h.handleStorageDelete(w, r)
	case path == "/list" && r.Method == "GET":
		h.handleStorageList(w, r)
	case path == "/upload-circuit-image" && r.Method == "POST":
		h.handleStorageUploadCircuitImage(w, r)
	default:
		http.NotFound(w, r)
	}
}

// healthCheck handles health check requests
func (h *HTTPHandler) healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	h.writeJSON(w, map[string]string{"status": "healthy"})
}

// Helper methods for HTTP response handling
func (h *HTTPHandler) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func (h *HTTPHandler) writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func (h *HTTPHandler) getUserIDFromContext(r *http.Request) string {
	// Get user from context (this would be set by auth middleware)
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		return ""
	}
	return user.UID
}

// Auth Handler Methods
func (h *HTTPHandler) handleAuthVerify(w http.ResponseWriter, r *http.Request) {
	var req controllers.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response, err := h.authController.VerifyToken(r.Context(), &req)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !response.Success {
		h.writeError(w, response.Error, http.StatusUnauthorized)
		return
	}

	h.writeJSON(w, response)
}

func (h *HTTPHandler) handleAuthCreateUser(w http.ResponseWriter, r *http.Request) {
	var req controllers.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response, err := h.authController.CreateUser(r.Context(), &req)
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

func (h *HTTPHandler) handleAuthGetUser(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	req := &controllers.AuthRequest{UID: userID}
	response, err := h.authController.GetUser(r.Context(), req)
	if err != nil {
		h.writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !response.Success {
		h.writeError(w, response.Error, http.StatusNotFound)
		return
	}

	h.writeJSON(w, response)
}

func (h *HTTPHandler) handleAuthUpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req controllers.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	req.UID = userID

	response, err := h.authController.UpdateUser(r.Context(), &req)
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

func (h *HTTPHandler) handleAuthDeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	req := &controllers.AuthRequest{UID: userID}
	response, err := h.authController.DeleteUser(r.Context(), req)
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

func (h *HTTPHandler) handleAuthSetCustomClaims(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	var req controllers.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	req.UID = userID

	response, err := h.authController.SetCustomClaims(r.Context(), &req)
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

// Circuit Handler Methods
func (h *HTTPHandler) handleCircuitsGet(w http.ResponseWriter, r *http.Request) {
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

func (h *HTTPHandler) handleCircuitsCreate(w http.ResponseWriter, r *http.Request) {
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

func (h *HTTPHandler) handleCircuitsGetOne(w http.ResponseWriter, r *http.Request) {
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
		h.writeError(w, response.Error, http.StatusNotFound)
		return
	}

	h.writeJSON(w, response)
}

func (h *HTTPHandler) handleCircuitsUpdate(w http.ResponseWriter, r *http.Request) {
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

func (h *HTTPHandler) handleCircuitsDelete(w http.ResponseWriter, r *http.Request) {
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

func (h *HTTPHandler) handleCircuitsTemplates(w http.ResponseWriter, r *http.Request) {
	// Templates are public, so no auth required for this endpoint
	templates := []map[string]interface{}{
		{
			"id":          "template_basic_circuit",
			"name":        "Basic Circuit",
			"description": "A simple circuit with resistor, LED, and battery",
			"category":    "basic",
		},
		{
			"id":          "template_amplifier",
			"name":        "Op-Amp Circuit",
			"description": "Basic operational amplifier configuration",
			"category":    "analog",
		},
	}

	h.writeJSON(w, map[string]interface{}{
		"success":   true,
		"templates": templates,
	})
}

func (h *HTTPHandler) handleCircuitsCreateFromTemplate(w http.ResponseWriter, r *http.Request) {
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

// Project Handler Methods
func (h *HTTPHandler) handleProjectsGet(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	req := &controllers.ProjectRequest{UserID: userID}
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

func (h *HTTPHandler) handleProjectsCreate(w http.ResponseWriter, r *http.Request) {
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

func (h *HTTPHandler) handleProjectsGetOne(w http.ResponseWriter, r *http.Request) {
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
		h.writeError(w, response.Error, http.StatusNotFound)
		return
	}

	h.writeJSON(w, response)
}

func (h *HTTPHandler) handleProjectsUpdate(w http.ResponseWriter, r *http.Request) {
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

func (h *HTTPHandler) handleProjectsDelete(w http.ResponseWriter, r *http.Request) {
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

func (h *HTTPHandler) handleProjectsDuplicate(w http.ResponseWriter, r *http.Request) {
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

// Storage Handler Methods
func (h *HTTPHandler) handleStorageUpload(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		h.writeError(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		h.writeError(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	req := &controllers.StorageRequest{
		UserID:     userID,
		File:       file,
		FileHeader: handler,
		Folder:     r.FormValue("folder"),
	}

	response, err := h.storageController.UploadFile(r.Context(), req)
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

func (h *HTTPHandler) handleStorageGetURL(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		h.writeError(w, "File path is required", http.StatusBadRequest)
		return
	}

	req := &controllers.StorageRequest{
		UserID:   userID,
		FilePath: filePath,
	}

	response, err := h.storageController.GetFileURL(r.Context(), req)
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

func (h *HTTPHandler) handleStorageDelete(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		h.writeError(w, "File path is required", http.StatusBadRequest)
		return
	}

	req := &controllers.StorageRequest{
		UserID:   userID,
		FilePath: filePath,
	}

	response, err := h.storageController.DeleteFile(r.Context(), req)
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

func (h *HTTPHandler) handleStorageList(w http.ResponseWriter, r *http.Request) {
	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	folder := r.URL.Query().Get("folder")

	req := &controllers.StorageRequest{
		UserID: userID,
		Folder: folder,
	}

	response, err := h.storageController.ListFiles(r.Context(), req)
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

func (h *HTTPHandler) handleStorageUploadCircuitImage(w http.ResponseWriter, r *http.Request) {
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

	// Parse multipart form
	err := r.ParseMultipartForm(5 << 20) // 5 MB limit for images
	if err != nil {
		h.writeError(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		h.writeError(w, "Unable to get image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	req := &controllers.StorageRequest{
		UserID:     userID,
		File:       file,
		FileHeader: handler,
		CircuitID:  circuitID,
	}

	response, err := h.storageController.UploadCircuitImage(r.Context(), req)
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
