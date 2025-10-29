package firebase

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
)

// AuthService handles Firebase authentication
type AuthService struct {
	client *auth.Client
	ctx    context.Context
}

// NewAuthService creates a new auth service
func NewAuthService(client *auth.Client, ctx context.Context) *AuthService {
	return &AuthService{
		client: client,
		ctx:    ctx,
	}
}

// VerifyIDToken verifies a Firebase ID token
func (a *AuthService) VerifyIDToken(idToken string) (*auth.Token, error) {
	token, err := a.client.VerifyIDToken(a.ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("error verifying ID token: %v", err)
	}
	return token, nil
}

// GetUser retrieves user information by UID
func (a *AuthService) GetUser(uid string) (*auth.UserRecord, error) {
	user, err := a.client.GetUser(a.ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return user, nil
}

// CreateUser creates a new user
func (a *AuthService) CreateUser(params *auth.UserToCreate) (*auth.UserRecord, error) {
	user, err := a.client.CreateUser(a.ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}
	return user, nil
}

// UpdateUser updates an existing user
func (a *AuthService) UpdateUser(uid string, params *auth.UserToUpdate) (*auth.UserRecord, error) {
	user, err := a.client.UpdateUser(a.ctx, uid, params)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}
	return user, nil
}

// DeleteUser deletes a user
func (a *AuthService) DeleteUser(uid string) error {
	err := a.client.DeleteUser(a.ctx, uid)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}

// SetCustomClaims sets custom claims for a user
func (a *AuthService) SetCustomClaims(uid string, claims map[string]interface{}) error {
	err := a.client.SetCustomUserClaims(a.ctx, uid, claims)
	if err != nil {
		return fmt.Errorf("error setting custom claims: %v", err)
	}
	return nil
}

// AuthMiddleware is a middleware for verifying Firebase tokens
func (a *AuthService) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Verify token
		decodedToken, err := a.VerifyIDToken(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user info to request context
		ctx := context.WithValue(r.Context(), "user", decodedToken)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

// GetUserFromContext extracts user info from request context
func GetUserFromContext(ctx context.Context) (*auth.Token, bool) {
	user, ok := ctx.Value("user").(*auth.Token)
	return user, ok
}
