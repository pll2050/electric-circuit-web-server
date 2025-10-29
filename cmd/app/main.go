package main

import (
	"fmt"
	"log"
	"net/http"

	"electric-circuit-web/server/internal/controllers"
	"electric-circuit-web/server/internal/handlers"
	"electric-circuit-web/server/internal/middleware"
	"electric-circuit-web/server/internal/repositories"
	"electric-circuit-web/server/internal/services"
	"electric-circuit-web/server/pkg/config"
	"electric-circuit-web/server/pkg/database"
	"electric-circuit-web/server/pkg/firebase"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Firebase
	var firebaseApp *firebase.FirebaseApp
	if cfg.FirebaseProjectID != "" {
		firebaseConfig := firebase.Config{
			ProjectID:             cfg.FirebaseProjectID,
			ServiceAccountKeyPath: cfg.FirebaseServiceAccountKeyPath,
			DatabaseURL:           cfg.FirebaseDatabaseURL,
		}

		firebaseApp, err = firebase.NewFirebaseApp(firebaseConfig)
		if err != nil {
			log.Printf("Warning: Failed to initialize Firebase: %v", err)
		} else {
			log.Println("Firebase initialized successfully")
			defer firebaseApp.Close()
		}
	}

	// Initialize components with dependency injection
	var httpHandler *handlers.HTTPHandler

	if firebaseApp != nil {
		// Initialize Repository layer
		firestoreRepository := repositories.NewFirestoreService(firebaseApp.Firestore, firebaseApp.GetContext())

		// Initialize Service layer
		authService := services.NewAuthService(firebaseApp.Auth, firebaseApp.GetContext())
		projectService := services.NewProjectService(firestoreRepository)
		circuitService := services.NewFirebaseCircuitService(firestoreRepository, projectService)

		// Initialize Controller layer
		authController := controllers.NewAuthController(authService)
		projectController := controllers.NewProjectController(projectService)
		circuitController := controllers.NewCircuitController(circuitService)
		storageController := controllers.NewStorageController()

		// Initialize Handler layer (HTTP Protocol)
		httpHandler = handlers.NewHTTPHandler(circuitController, projectController, authController, storageController)

		// Initialize Middleware (for future use)
		_ = middleware.NewAuthMiddleware(authService)
	}

	// Setup routes
	if httpHandler != nil {
		httpHandler.SetupRoutes()
	} else {
		// Fallback health check if Firebase is not available
		http.HandleFunc("/api/health", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "healthy", "firebase": "disabled"}`))
		}))
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
