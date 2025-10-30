package controllers

import (
	"context"

	"electric-circuit-web/server/internal/models"
	"electric-circuit-web/server/internal/services"

	"firebase.google.com/go/v4/auth"
)

// AuthController handles authentication-related business coordination
type AuthController struct {
	authService *services.AuthService
}

// NewAuthController creates a new auth controller
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// AuthRequest represents an authentication operation request
type AuthRequest struct {
	Token        string                 `json:"token,omitempty"`
	UID          string                 `json:"uid,omitempty"`          // Firebase UID
	Email        string                 `json:"email,omitempty"`
	Password     string                 `json:"password,omitempty"`
	DisplayName  string                 `json:"display_name,omitempty"`
	PhotoURL     string                 `json:"photo_url,omitempty"`
	Provider     string                 `json:"provider,omitempty"`     // Auth provider (google, email, etc.)
	CustomClaims map[string]interface{} `json:"custom_claims,omitempty"`
	IDToken      string                 `json:"id_token,omitempty"`     // Google OAuth ID token (deprecated, use UID instead)
}

// AuthResponse represents an authentication operation response
type AuthResponse struct {
	Success      bool                   `json:"success"`
	Message      string                 `json:"message"`
	User         *auth.UserRecord       `json:"user,omitempty"`          // Firebase user
	DBUser       *models.User           `json:"db_user,omitempty"`       // PostgreSQL user
	Token        string                 `json:"token,omitempty"`
	CustomClaims map[string]interface{} `json:"custom_claims,omitempty"`
	Error        string                 `json:"error,omitempty"`
}

// VerifyToken verifies a Firebase token
func (c *AuthController) VerifyToken(ctx context.Context, req *AuthRequest) (*AuthResponse, error) {
	// Validate request
	if req.Token == "" {
		return &AuthResponse{
			Success: false,
			Error:   "Token is required",
		}, nil
	}

	// Call service
	token, err := c.authService.VerifyIDToken(req.Token)
	if err != nil {
		return &AuthResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &AuthResponse{
		Success: true,
		Message: "Token verified successfully",
		Token:   token.UID, // Return UID as verification proof
	}, nil
}

// CreateUser creates a new user
func (c *AuthController) CreateUser(ctx context.Context, req *AuthRequest) (*AuthResponse, error) {
	// Validate request
	if req.Email == "" {
		return &AuthResponse{
			Success: false,
			Error:   "Email is required",
		}, nil
	}

	if req.Password == "" {
		return &AuthResponse{
			Success: false,
			Error:   "Password is required",
		}, nil
	}

	// Prepare user data
	userData := &auth.UserToCreate{}
	userData.Email(req.Email).Password(req.Password)

	if req.DisplayName != "" {
		userData.DisplayName(req.DisplayName)
	}

	if req.PhotoURL != "" {
		userData.PhotoURL(req.PhotoURL)
	}

	// Call service
	user, err := c.authService.CreateUser(userData)
	if err != nil {
		return &AuthResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &AuthResponse{
		Success: true,
		Message: "User created successfully",
		User:    user,
	}, nil
}

// GetUser retrieves user information
func (c *AuthController) GetUser(ctx context.Context, req *AuthRequest) (*AuthResponse, error) {
	// Validate request
	if req.UID == "" {
		return &AuthResponse{
			Success: false,
			Error:   "User UID is required",
		}, nil
	}

	// Call service
	user, err := c.authService.GetUser(req.UID)
	if err != nil {
		return &AuthResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &AuthResponse{
		Success: true,
		Message: "User retrieved successfully",
		User:    user,
	}, nil
}

// UpdateUser updates user information
func (c *AuthController) UpdateUser(ctx context.Context, req *AuthRequest) (*AuthResponse, error) {
	// Validate request
	if req.UID == "" {
		return &AuthResponse{
			Success: false,
			Error:   "User UID is required",
		}, nil
	}

	// Prepare update data
	updateData := &auth.UserToUpdate{}

	if req.Email != "" {
		updateData.Email(req.Email)
	}

	if req.DisplayName != "" {
		updateData.DisplayName(req.DisplayName)
	}

	if req.PhotoURL != "" {
		updateData.PhotoURL(req.PhotoURL)
	}

	if req.Password != "" {
		updateData.Password(req.Password)
	}

	// Call service
	user, err := c.authService.UpdateUser(req.UID, updateData)
	if err != nil {
		return &AuthResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &AuthResponse{
		Success: true,
		Message: "User updated successfully",
		User:    user,
	}, nil
}

// DeleteUser deletes a user
func (c *AuthController) DeleteUser(ctx context.Context, req *AuthRequest) (*AuthResponse, error) {
	// Validate request
	if req.UID == "" {
		return &AuthResponse{
			Success: false,
			Error:   "User UID is required",
		}, nil
	}

	// Call service
	err := c.authService.DeleteUser(req.UID)
	if err != nil {
		return &AuthResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &AuthResponse{
		Success: true,
		Message: "User deleted successfully",
	}, nil
}

// SetCustomClaims sets custom claims for a user
func (c *AuthController) SetCustomClaims(ctx context.Context, req *AuthRequest) (*AuthResponse, error) {
	// Validate request
	if req.UID == "" {
		return &AuthResponse{
			Success: false,
			Error:   "User UID is required",
		}, nil
	}

	if req.CustomClaims == nil {
		return &AuthResponse{
			Success: false,
			Error:   "Custom claims are required",
		}, nil
	}

	// Call service
	err := c.authService.SetCustomClaims(req.UID, req.CustomClaims)
	if err != nil {
		return &AuthResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &AuthResponse{
		Success:      true,
		Message:      "Custom claims set successfully",
		CustomClaims: req.CustomClaims,
	}, nil
}

// Register handles Google OAuth user registration
// Client sends: { id_token: "JWT-token", provider: "google" }
// 1. Verify ID token and extract UID
// 2. Get user info from Firebase using UID
// 3. Save to PostgreSQL
func (c *AuthController) Register(ctx context.Context, req *AuthRequest) (*AuthResponse, error) {
	// Note: Basic validation (required fields) is done in Handler layer

	// Verify ID token and extract user info
	token, err := c.authService.VerifyIDToken(req.IDToken)
	if err != nil {
		return &AuthResponse{
			Success: false,
			Error:   "Invalid ID token: " + err.Error(),
		}, err
	}

	// Extract UID from verified token
	uid := token.UID

	// Set default provider if not specified
	provider := req.Provider
	if provider == "" {
		provider = "google"
	}

	// Call service to register user (Firebase -> PostgreSQL)
	user, err := c.authService.RegisterUserByUID(uid, provider)
	if err != nil {
		return &AuthResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &AuthResponse{
		Success: true,
		Message: "User registered successfully",
		DBUser:  user,
	}, nil
}

// ListUsers lists all users from PostgreSQL database
func (c *AuthController) ListUsers(ctx context.Context) ([]map[string]interface{}, error) {
	users, err := c.authService.ListUsers()
	if err != nil {
		return nil, err
	}

	// Convert to map format for JSON response
	result := make([]map[string]interface{}, len(users))
	for i, user := range users {
		result[i] = map[string]interface{}{
			"uid":          user.UID,
			"email":        user.Email,
			"displayName":  user.DisplayName,
			"photoURL":     user.PhotoURL,
			"provider":     user.Provider,
			"createdAt":    user.CreatedAt,
			"lastLoginAt":  user.LastLoginAt,
		}
	}

	return result, nil
}
