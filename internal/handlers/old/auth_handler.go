package handlers

import (
	"encoding/json"
	"net/http"

	"electric-circuit-web/server/internal/services"

	"firebase.google.com/go/v4/auth"
)

// AuthHandler handles Firebase authentication-related HTTP requests
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// VerifyToken verifies a Firebase ID token
func (h *AuthHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := h.authService.VerifyIDToken(request.Token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"valid":  true,
		"uid":    token.UID,
		"email":  token.Claims["email"],
		"name":   token.Claims["name"],
		"claims": token.Claims,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateUser creates a new user
func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"displayName,omitempty"`
		PhotoURL    string `json:"photoURL,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params := (&auth.UserToCreate{}).
		Email(request.Email).
		Password(request.Password)

	if request.DisplayName != "" {
		params = params.DisplayName(request.DisplayName)
	}
	if request.PhotoURL != "" {
		params = params.PhotoURL(request.PhotoURL)
	}

	user, err := h.authService.CreateUser(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"uid":   user.UID,
		"email": user.Email,
		"name":  user.DisplayName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUser retrieves user information
func (h *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uid := r.URL.Query().Get("uid")
	if uid == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	user, err := h.authService.GetUser(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"uid":           user.UID,
		"email":         user.Email,
		"displayName":   user.DisplayName,
		"photoURL":      user.PhotoURL,
		"disabled":      user.Disabled,
		"emailVerified": user.EmailVerified,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateUser updates user information
func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uid := r.URL.Query().Get("uid")
	if uid == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var request struct {
		Email       string `json:"email,omitempty"`
		DisplayName string `json:"displayName,omitempty"`
		PhotoURL    string `json:"photoURL,omitempty"`
		Disabled    *bool  `json:"disabled,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params := &auth.UserToUpdate{}

	if request.Email != "" {
		params = params.Email(request.Email)
	}
	if request.DisplayName != "" {
		params = params.DisplayName(request.DisplayName)
	}
	if request.PhotoURL != "" {
		params = params.PhotoURL(request.PhotoURL)
	}
	if request.Disabled != nil {
		params = params.Disabled(*request.Disabled)
	}

	user, err := h.authService.UpdateUser(uid, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"uid":         user.UID,
		"email":       user.Email,
		"displayName": user.DisplayName,
		"photoURL":    user.PhotoURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteUser deletes a user
func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uid := r.URL.Query().Get("uid")
	if uid == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	err := h.authService.DeleteUser(uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "User deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SetCustomClaims sets custom claims for a user
func (h *AuthHandler) SetCustomClaims(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uid := r.URL.Query().Get("uid")
	if uid == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var request struct {
		Claims map[string]interface{} `json:"claims"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.authService.SetCustomClaims(uid, request.Claims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Custom claims set successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
