package handlers

import (
	"encoding/json"
	"net/http"

	"electric-circuit-web/server/internal/controllers"
)

// StorageHandler handles storage-related HTTP requests
type StorageHandler struct {
	storageController controllers.StorageControllerInterface
}

// NewStorageHandler creates a new storage handler
func NewStorageHandler(storageController controllers.StorageControllerInterface) *StorageHandler {
	return &StorageHandler{
		storageController: storageController,
	}
}

// getUserIDFromContext extracts user ID from request context
// This would typically be set by authentication middleware
func (h *StorageHandler) getUserIDFromContext(r *http.Request) string {
	// For now, return a placeholder - this should be replaced with actual auth context
	return r.Header.Get("X-User-ID")
}

// writeError writes an error response
func (h *StorageHandler) writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// writeJSON writes a JSON response
func (h *StorageHandler) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// HandleUploadFile handles file upload requests
func (h *StorageHandler) HandleUploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32 MB max
	if err != nil {
		h.writeError(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.writeError(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	folder := r.FormValue("folder")

	req := &controllers.StorageRequest{
		UserID:     userID,
		File:       file,
		FileHeader: header,
		Folder:     folder,
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

// HandleGetFileURL handles getting file URL requests
func (h *StorageHandler) HandleGetFileURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	filePath := r.URL.Query().Get("filePath")
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

// HandleDeleteFile handles file deletion requests
func (h *StorageHandler) HandleDeleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	filePath := r.URL.Query().Get("filePath")
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

// HandleListFiles handles listing files requests
func (h *StorageHandler) HandleListFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

// HandleUploadCircuitImage handles circuit image upload requests
func (h *StorageHandler) HandleUploadCircuitImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := h.getUserIDFromContext(r)
	if userID == "" {
		h.writeError(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32 MB max
	if err != nil {
		h.writeError(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		h.writeError(w, "Unable to retrieve image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	circuitID := r.FormValue("circuitId")
	if circuitID == "" {
		h.writeError(w, "Circuit ID is required", http.StatusBadRequest)
		return
	}

	req := &controllers.StorageRequest{
		UserID:     userID,
		File:       file,
		FileHeader: header,
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
