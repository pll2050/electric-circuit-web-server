package firebase

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"
)

// FirebaseApp wraps Firebase services
type FirebaseApp struct {
	App       *firebase.App
	Auth      *auth.Client
	Firestore *firestore.Client
	Database  *db.Client
	ctx       context.Context
}

// Config holds Firebase configuration
type Config struct {
	ServiceAccountKeyPath string
	ProjectID             string
	DatabaseURL           string
}

// NewFirebaseApp initializes a new Firebase application
func NewFirebaseApp(config Config) (*FirebaseApp, error) {
	ctx := context.Background()

	// Initialize Firebase App with service account
	var app *firebase.App
	var err error

	if config.ServiceAccountKeyPath != "" {
		// Use service account key file
		opt := option.WithCredentialsFile(config.ServiceAccountKeyPath)
		appConfig := &firebase.Config{
			ProjectID:   config.ProjectID,
			DatabaseURL: config.DatabaseURL,
		}
		app, err = firebase.NewApp(ctx, appConfig, opt)
	} else {
		// Use default credentials (for Cloud Run, App Engine, etc.)
		appConfig := &firebase.Config{
			ProjectID:   config.ProjectID,
			DatabaseURL: config.DatabaseURL,
		}
		app, err = firebase.NewApp(ctx, appConfig)
	}

	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	// Initialize Auth client
	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase auth: %v", err)
	}

	// Initialize Firestore client
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing firestore: %v", err)
	}

	// Initialize Realtime Database client (optional)
	var dbClient *db.Client
	if config.DatabaseURL != "" {
		dbClient, err = app.Database(ctx)
		if err != nil {
			log.Printf("Warning: error initializing realtime database: %v", err)
			dbClient = nil
		}
	}

	return &FirebaseApp{
		App:       app,
		Auth:      authClient,
		Firestore: firestoreClient,
		Database:  dbClient,
		ctx:       ctx,
	}, nil
}

// Close closes all Firebase connections
func (f *FirebaseApp) Close() error {
	if f.Firestore != nil {
		return f.Firestore.Close()
	}
	return nil
}

// GetContext returns the context used by Firebase services
func (f *FirebaseApp) GetContext() context.Context {
	return f.ctx
}
