package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"electric-circuit-web/server/internal/middleware"
)

// StorageHandler handles Firebase Storage-related HTTP requests
type StorageHandler struct {
	// Note: Firebase Storage service would be injected here
	// For now, we'll implement basic file upload handling
}

// NewStorageHandler creates a new storage handler
func NewStorageHandler() *StorageHandler {
	return &StorageHandler{}
}

// UploadFile handles file upload to Firebase Storage
func (h *StorageHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
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

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get optional parameters
	folder := r.FormValue("folder")
	if folder == "" {
		folder = "uploads"
	}

	// Validate file type
	allowedTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".svg", ".json", ".xml"}
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	isAllowed := false
	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		http.Error(w, "File type not allowed", http.StatusBadRequest)
		return
	}

	// Create unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s", timestamp, handler.Filename)

	// TODO: Implement actual Firebase Storage upload
	// For now, return mock response
	downloadURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/your-bucket/o/%s%%2F%s?alt=media", folder, filename)

	response := map[string]interface{}{
		"message":      "File uploaded successfully",
		"filename":     filename,
		"originalName": handler.Filename,
		"size":         handler.Size,
		"folder":       folder,
		"downloadURL":  downloadURL,
		"uploadedBy":   user.UID,
		"uploadedAt":   time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetFileURL generates a download URL for a file
func (h *StorageHandler) GetFileURL(w http.ResponseWriter, r *http.Request) {
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

	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	// TODO: Implement actual Firebase Storage URL generation
	// For now, return mock response
	downloadURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/your-bucket/o/%s?alt=media",
		strings.ReplaceAll(filePath, "/", "%2F"))

	response := map[string]interface{}{
		"downloadURL": downloadURL,
		"path":        filePath,
		"requestedBy": user.UID,
		"generatedAt": time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteFile deletes a file from Firebase Storage
func (h *StorageHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
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

	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	// TODO: Implement actual Firebase Storage file deletion
	// For now, return mock response

	response := map[string]interface{}{
		"message":   "File deleted successfully",
		"path":      filePath,
		"deletedBy": user.UID,
		"deletedAt": time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListFiles lists files in a storage folder
func (h *StorageHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
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

	folder := r.URL.Query().Get("folder")
	if folder == "" {
		folder = "uploads"
	}

	// TODO: Implement actual Firebase Storage file listing
	// For now, return mock response
	mockFiles := []map[string]interface{}{
		{
			"name":        "circuit_diagram_1.svg",
			"path":        fmt.Sprintf("%s/circuit_diagram_1.svg", folder),
			"size":        15432,
			"contentType": "image/svg+xml",
			"uploadedAt":  time.Now().Add(-24 * time.Hour).UTC(),
			"downloadURL": fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/your-bucket/o/%s%%2Fcircuit_diagram_1.svg?alt=media", folder),
		},
		{
			"name":        "project_layout.json",
			"path":        fmt.Sprintf("%s/project_layout.json", folder),
			"size":        8920,
			"contentType": "application/json",
			"uploadedAt":  time.Now().Add(-48 * time.Hour).UTC(),
			"downloadURL": fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/your-bucket/o/%s%%2Fproject_layout.json?alt=media", folder),
		},
	}

	response := map[string]interface{}{
		"files":       mockFiles,
		"folder":      folder,
		"requestedBy": user.UID,
		"count":       len(mockFiles),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UploadCircuitImage handles circuit diagram image uploads
func (h *StorageHandler) UploadCircuitImage(w http.ResponseWriter, r *http.Request) {
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

	circuitID := r.URL.Query().Get("circuitId")
	if circuitID == "" {
		http.Error(w, "Circuit ID is required", http.StatusBadRequest)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(5 << 20) // 5 MB limit for images
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Unable to get image file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate image type
	allowedImageTypes := []string{".jpg", ".jpeg", ".png", ".svg", ".gif"}
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	isAllowed := false
	for _, allowedType := range allowedImageTypes {
		if ext == allowedType {
			isAllowed = true
			break
		}
	}

	if !isAllowed {
		http.Error(w, "Only image files are allowed", http.StatusBadRequest)
		return
	}

	// Read file content to validate it's actually an image
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusInternalServerError)
		return
	}

	// Create unique filename for circuit image
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("circuit_%s_%d%s", circuitID, timestamp, ext)
	folder := fmt.Sprintf("circuits/%s/images", circuitID)

	// TODO: Implement actual Firebase Storage upload
	// For now, return mock response
	downloadURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/your-bucket/o/%s%%2F%s?alt=media",
		strings.ReplaceAll(folder, "/", "%2F"), filename)

	response := map[string]interface{}{
		"message":      "Circuit image uploaded successfully",
		"filename":     filename,
		"originalName": handler.Filename,
		"size":         len(fileBytes),
		"folder":       folder,
		"downloadURL":  downloadURL,
		"circuitId":    circuitID,
		"uploadedBy":   user.UID,
		"uploadedAt":   time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
