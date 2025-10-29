package services

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/auth"
)

// AuthService handles Firebase authentication operations
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
func (s *AuthService) VerifyIDToken(idToken string) (*auth.Token, error) {
	token, err := s.client.VerifyIDToken(s.ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("error verifying ID token: %v", err)
	}
	return token, nil
}

// GetUser retrieves user information by UID
func (s *AuthService) GetUser(uid string) (*auth.UserRecord, error) {
	user, err := s.client.GetUser(s.ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}
	return user, nil
}

// CreateUser creates a new user
func (s *AuthService) CreateUser(params *auth.UserToCreate) (*auth.UserRecord, error) {
	user, err := s.client.CreateUser(s.ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %v", err)
	}
	return user, nil
}

// UpdateUser updates an existing user
func (s *AuthService) UpdateUser(uid string, params *auth.UserToUpdate) (*auth.UserRecord, error) {
	user, err := s.client.UpdateUser(s.ctx, uid, params)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %v", err)
	}
	return user, nil
}

// DeleteUser deletes a user
func (s *AuthService) DeleteUser(uid string) error {
	err := s.client.DeleteUser(s.ctx, uid)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}
	return nil
}

// SetCustomClaims sets custom claims for a user
func (s *AuthService) SetCustomClaims(uid string, claims map[string]interface{}) error {
	err := s.client.SetCustomUserClaims(s.ctx, uid, claims)
	if err != nil {
		return fmt.Errorf("error setting custom claims: %v", err)
	}
	return nil
}
