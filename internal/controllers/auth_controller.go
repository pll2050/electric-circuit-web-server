package controllers

import (
	"context"

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
	UID          string                 `json:"uid,omitempty"`
	Email        string                 `json:"email,omitempty"`
	Password     string                 `json:"password,omitempty"`
	DisplayName  string                 `json:"display_name,omitempty"`
	PhotoURL     string                 `json:"photo_url,omitempty"`
	CustomClaims map[string]interface{} `json:"custom_claims,omitempty"`
}

// AuthResponse represents an authentication operation response
type AuthResponse struct {
	Success      bool                   `json:"success"`
	Message      string                 `json:"message"`
	User         *auth.UserRecord       `json:"user,omitempty"`
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
