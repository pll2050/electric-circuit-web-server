package services

import (
	"context"
	"fmt"
	"time"

	"electric-circuit-web/server/internal/models"
	"electric-circuit-web/server/internal/repositories"

	"firebase.google.com/go/v4/auth"
)

// AuthService handles Firebase authentication operations
type AuthService struct {
	client         *auth.Client
	ctx            context.Context
	userRepository *repositories.UserRepository
}

// NewAuthService creates a new auth service
func NewAuthService(client *auth.Client, ctx context.Context, userRepository *repositories.UserRepository) *AuthService {
	return &AuthService{
		client:         client,
		ctx:            ctx,
		userRepository: userRepository,
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

// RegisterUser registers a user with Google OAuth
// 1. Verify Firebase token from Google OAuth
// 2. Get or create user in Firebase
// 3. Save user to PostgreSQL
func (s *AuthService) RegisterUser(idToken string) (*models.User, error) {
	// Step 1: Verify Firebase ID token from Google OAuth
	token, err := s.VerifyIDToken(idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token: %v", err)
	}

	// Step 2: Get user info from Firebase
	firebaseUser, err := s.GetUser(token.UID)
	if err != nil {
		return nil, fmt.Errorf("failed to get Firebase user: %v", err)
	}

	// Check if user already exists in PostgreSQL
	existingUser, err := s.userRepository.GetByID(firebaseUser.UID)
	if err == nil && existingUser != nil {
		// User already exists, return existing user
		return existingUser, nil
	}

	// Step 3: Save user to PostgreSQL
	user := &models.User{
		ID:          firebaseUser.UID,
		Email:       firebaseUser.Email,
		DisplayName: firebaseUser.DisplayName,
		PhotoURL:    firebaseUser.PhotoURL,
		Provider:    "google", // Google OAuth provider
	}

	err = s.userRepository.Create(user)
	if err != nil {
		// If PostgreSQL save fails, we should consider rolling back Firebase user
		// For now, we'll just return the error
		return nil, fmt.Errorf("failed to save user to PostgreSQL: %v", err)
	}

	return user, nil
}

// RegisterUserByUID registers a user using Firebase UID
// Client sends: { uid: "firebase-uid", provider: "google" }
// 1. Get user info from Firebase using UID
// 2. Save to PostgreSQL
func (s *AuthService) RegisterUserByUID(uid string, provider string) (*models.User, error) {
	// Step 1: Get user info from Firebase
	firebaseUser, err := s.GetUser(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get Firebase user: %v", err)
	}

	// Check if user already exists in PostgreSQL
	existingUser, err := s.userRepository.GetByID(firebaseUser.UID)
	if err == nil && existingUser != nil {
		// User already exists, return existing user
		return existingUser, nil
	}

	// Set default provider if not provided
	if provider == "" {
		provider = "google"
	}

	// Step 2: Save user to PostgreSQL
	user := &models.User{
		ID:          firebaseUser.UID,
		Email:       firebaseUser.Email,
		DisplayName: firebaseUser.DisplayName,
		PhotoURL:    firebaseUser.PhotoURL,
		Provider:    provider,
	}

	err = s.userRepository.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to save user to PostgreSQL: %v", err)
	}

	return user, nil
}

// GetUserFromDB retrieves user from PostgreSQL by ID
func (s *AuthService) GetUserFromDB(uid string) (*models.User, error) {
	user, err := s.userRepository.GetByID(uid)
	if err != nil {
		return nil, fmt.Errorf("failed to get user from database: %v", err)
	}
	return user, nil
}

// ListUsers retrieves all users from Firebase Authentication
func (s *AuthService) ListUsers() ([]*models.User, error) {
	// Create iterator to list all users
	iter := s.client.Users(s.ctx, "")

	var users []*models.User

	// Iterate through all users
	for {
		user, err := iter.Next()
		if err != nil {
			// Check if we've reached the end
			if err.Error() == "no more items in iterator" {
				break
			}
			return nil, fmt.Errorf("failed to iterate users: %v", err)
		}

		// Convert Firebase user to our User model
		modelUser := &models.User{
			ID:          user.UID,
			UID:         user.UID,
			Email:       user.Email,
			DisplayName: user.DisplayName,
			PhotoURL:    user.PhotoURL,
			Provider:    getProvider(user.ProviderUserInfo),
		}

		// Set timestamps if available
		if user.UserMetadata != nil {
			// Convert Unix milliseconds to time.Time
			modelUser.CreatedAt = convertUnixMillisToTime(user.UserMetadata.CreationTimestamp)

			if user.UserMetadata.LastLogInTimestamp > 0 {
				lastLogin := convertUnixMillisToTime(user.UserMetadata.LastLogInTimestamp)
				modelUser.UpdatedAt = lastLogin
				modelUser.LastLoginAt = &lastLogin
			} else {
				modelUser.UpdatedAt = modelUser.CreatedAt
			}
		}

		users = append(users, modelUser)
	}

	return users, nil
}

// getProvider extracts the primary auth provider from user info
func getProvider(providerData []*auth.UserInfo) string {
	if len(providerData) == 0 {
		return "email"
	}

	// Return the first provider ID (e.g., "google.com", "password", etc.)
	providerID := providerData[0].ProviderID

	// Simplify provider name
	switch providerID {
	case "google.com":
		return "google"
	case "password":
		return "email"
	default:
		return providerID
	}
}

// convertUnixMillisToTime converts Unix milliseconds timestamp to time.Time
func convertUnixMillisToTime(millis int64) time.Time {
	return time.Unix(millis/1000, (millis%1000)*1000000)
}
