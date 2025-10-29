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