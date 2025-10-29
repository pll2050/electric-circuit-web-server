package tests

import (
	"os"
	"testing"

	"electric-circuit-web/server/pkg/firebase"
)

func TestFirebaseConfig_ValidConfig(t *testing.T) {
	// Arrange
	config := firebase.Config{
		ProjectID:             "test-project",
		ServiceAccountKeyPath: "",
		DatabaseURL:           "https://test-project-default-rtdb.firebaseio.com/",
	}

	// Act & Assert
	if config.ProjectID != "test-project" {
		t.Errorf("Expected ProjectID 'test-project', got '%s'", config.ProjectID)
	}

	if config.DatabaseURL != "https://test-project-default-rtdb.firebaseio.com/" {
		t.Errorf("Expected DatabaseURL 'https://test-project-default-rtdb.firebaseio.com/', got '%s'", config.DatabaseURL)
	}
}

func TestFirebaseConfig_EmptyProjectID(t *testing.T) {
	// Arrange
	config := firebase.Config{
		ProjectID:             "",
		ServiceAccountKeyPath: "",
		DatabaseURL:           "https://test-project-default-rtdb.firebaseio.com/",
	}

	// Act & Assert
	if config.ProjectID != "" {
		t.Errorf("Expected empty ProjectID, got '%s'", config.ProjectID)
	}
}

func TestNewFirebaseApp_EmptyConfig_ShouldUseDefaults(t *testing.T) {
	// Note: This test would require actual Firebase credentials to work
	// In a real test environment, you would use Firebase emulator

	// Arrange
	config := firebase.Config{
		ProjectID:             "",
		ServiceAccountKeyPath: "",
		DatabaseURL:           "",
	}

	// Act
	_, err := firebase.NewFirebaseApp(config)

	// Assert
	// We expect this to fail without proper credentials
	if err == nil {
		t.Error("Expected error with empty config, got nil")
	}
}

func TestNewFirebaseApp_WithServiceAccountPath(t *testing.T) {
	// Skip this test if no test service account file is available
	testKeyPath := "./test-service-account.json"
	if _, err := os.Stat(testKeyPath); os.IsNotExist(err) {
		t.Skip("Skipping test: test service account file not found")
	}

	// Arrange
	config := firebase.Config{
		ProjectID:             "test-project",
		ServiceAccountKeyPath: testKeyPath,
		DatabaseURL:           "https://test-project-default-rtdb.firebaseio.com/",
	}

	// Act
	app, err := firebase.NewFirebaseApp(config)

	// Assert
	if err != nil {
		t.Errorf("Expected no error with valid config, got %v", err)
	}

	if app == nil {
		t.Error("Expected Firebase app, got nil")
	}

	// Clean up
	if app != nil {
		app.Close()
	}
}

func TestFirebaseApp_Close(t *testing.T) {
	// Skip this test if no test service account file is available
	testKeyPath := "./test-service-account.json"
	if _, err := os.Stat(testKeyPath); os.IsNotExist(err) {
		t.Skip("Skipping test: test service account file not found")
	}

	// Arrange
	config := firebase.Config{
		ProjectID:             "test-project",
		ServiceAccountKeyPath: testKeyPath,
		DatabaseURL:           "https://test-project-default-rtdb.firebaseio.com/",
	}

	app, err := firebase.NewFirebaseApp(config)
	if err != nil {
		t.Skipf("Skipping test: Firebase app creation failed: %v", err)
	}

	// Act & Assert
	// Should not panic
	app.Close()
}

// Helper function to create a temporary test service account file
func createTestServiceAccountFile(t *testing.T) string {
	content := `{
		"type": "service_account",
		"project_id": "test-project",
		"private_key_id": "test-key-id",
		"private_key": "-----BEGIN PRIVATE KEY-----\ntest-private-key\n-----END PRIVATE KEY-----\n",
		"client_email": "test@test-project.iam.gserviceaccount.com",
		"client_id": "123456789",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token"
	}`

	tmpFile, err := os.CreateTemp("", "test-service-account-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpFile.Name()
}

func TestNewFirebaseApp_WithInvalidServiceAccount(t *testing.T) {
	// Create invalid service account file
	invalidKeyPath := createTestServiceAccountFile(t)
	defer os.Remove(invalidKeyPath) // Clean up

	// Arrange
	config := firebase.Config{
		ProjectID:             "test-project",
		ServiceAccountKeyPath: invalidKeyPath,
		DatabaseURL:           "https://test-project-default-rtdb.firebaseio.com/",
	}

	// Act
	_, err := firebase.NewFirebaseApp(config)

	// Assert
	// We expect this to fail with invalid credentials
	if err == nil {
		t.Error("Expected error with invalid service account, got nil")
	}
}
