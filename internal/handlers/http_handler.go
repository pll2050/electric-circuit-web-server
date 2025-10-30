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
	authHandler    *AuthHandler
	circuitHandler *CircuitHandler
	projectHandler *ProjectHandler
	storageHandler *StorageHandler
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(
	circuitController *controllers.CircuitController,
	projectController *controllers.ProjectController,
	authController *controllers.AuthController,
	storageController *controllers.StorageController,
) *HTTPHandler {
	return &HTTPHandler{
		authHandler:    NewAuthHandler(authController),
		circuitHandler: NewCircuitHandler(circuitController),
		projectHandler: NewProjectHandler(projectController),
		storageHandler: NewStorageHandler(storageController),
	}
}

// SetupRoutes configures all HTTP routes
func (h *HTTPHandler) SetupRoutes() {
	// Health check (no auth required)
	http.HandleFunc("/api/health", middleware.CORS(h.healthCheck))

	// Public auth routes (no auth required)
	http.HandleFunc("/api/auth/", middleware.CORS(h.routeAuth))

	// Admin routes (require auth)
	http.HandleFunc("/api/users", middleware.CORS(h.withAuth(h.routeUsers)))

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
	// Handle OPTIONS preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/auth")

	switch {
	case path == "/register" && r.Method == "POST":
		h.authHandler.HandleRegister(w, r)
	case path == "/verify" && r.Method == "POST":
		h.authHandler.HandleVerifyToken(w, r)
	case path == "/create-user" && r.Method == "POST":
		h.authHandler.HandleCreateUser(w, r)
	case path == "/get-user" && r.Method == "GET":
		h.authHandler.HandleGetUser(w, r)
	case path == "/update-user" && r.Method == "PUT":
		h.authHandler.HandleUpdateUser(w, r)
	case path == "/delete-user" && r.Method == "DELETE":
		h.authHandler.HandleDeleteUser(w, r)
	case path == "/set-custom-claims" && r.Method == "POST":
		h.authHandler.HandleSetCustomClaims(w, r)
	default:
		http.NotFound(w, r)
	}
}

// routeCircuits handles circuit-related HTTP routing
func (h *HTTPHandler) routeCircuits(w http.ResponseWriter, r *http.Request) {
	// Handle OPTIONS preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := r.URL.Path

	switch {
	case path == "/api/circuits" && r.Method == "GET":
		h.circuitHandler.HandleGetCircuits(w, r)
	case path == "/api/circuits/create" && r.Method == "POST":
		h.circuitHandler.HandleCreateCircuit(w, r)
	case path == "/api/circuits/get" && r.Method == "GET":
		h.circuitHandler.HandleGetCircuit(w, r)
	case path == "/api/circuits/update" && r.Method == "PUT":
		h.circuitHandler.HandleUpdateCircuit(w, r)
	case path == "/api/circuits/delete" && r.Method == "DELETE":
		h.circuitHandler.HandleDeleteCircuit(w, r)
	case path == "/api/circuits/templates" && r.Method == "GET":
		h.circuitHandler.HandleGetTemplates(w, r)
	case path == "/api/circuits/create-from-template" && r.Method == "POST":
		h.circuitHandler.HandleCreateFromTemplate(w, r)
	default:
		http.NotFound(w, r)
	}
}

// routeProjects handles project-related HTTP routing
func (h *HTTPHandler) routeProjects(w http.ResponseWriter, r *http.Request) {
	// Handle OPTIONS preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := r.URL.Path

	switch {
	case path == "/api/projects" && r.Method == "GET":
		h.projectHandler.HandleGetProjects(w, r)
	case path == "/api/projects/create" && r.Method == "POST":
		h.projectHandler.HandleCreateProject(w, r)
	case path == "/api/projects/get" && r.Method == "GET":
		h.projectHandler.HandleGetProject(w, r)
	case path == "/api/projects/update" && r.Method == "PUT":
		h.projectHandler.HandleUpdateProject(w, r)
	case path == "/api/projects/delete" && r.Method == "DELETE":
		h.projectHandler.HandleDeleteProject(w, r)
	case path == "/api/projects/duplicate" && r.Method == "POST":
		h.projectHandler.HandleDuplicateProject(w, r)
	default:
		http.NotFound(w, r)
	}
}

// routeStorage handles storage-related HTTP routing
func (h *HTTPHandler) routeStorage(w http.ResponseWriter, r *http.Request) {
	// Handle OPTIONS preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/storage")

	switch {
	case path == "/upload" && r.Method == "POST":
		h.storageHandler.HandleUploadFile(w, r)
	case path == "/url" && r.Method == "GET":
		h.storageHandler.HandleGetFileURL(w, r)
	case path == "/delete" && r.Method == "DELETE":
		h.storageHandler.HandleDeleteFile(w, r)
	case path == "/list" && r.Method == "GET":
		h.storageHandler.HandleListFiles(w, r)
	case path == "/upload-circuit-image" && r.Method == "POST":
		h.storageHandler.HandleUploadCircuitImage(w, r)
	default:
		http.NotFound(w, r)
	}
}

// routeUsers handles user-related HTTP routing
func (h *HTTPHandler) routeUsers(w http.ResponseWriter, r *http.Request) {
	// Handle OPTIONS preflight request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := r.URL.Path

	switch {
	case path == "/api/users" && r.Method == "GET":
		h.authHandler.HandleListUsers(w, r)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
