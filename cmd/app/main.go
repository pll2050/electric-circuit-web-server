package main

import (
	"fmt"
	"log"
	"net/http"

	"electric-circuit-web/server/internal/handlers"
	"electric-circuit-web/server/pkg/config"
	"electric-circuit-web/server/pkg/database"
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

	// Initialize handlers
	handler := handlers.NewHandler(db)

	// Setup routes
	http.HandleFunc("/api/health", handler.HealthCheck)
	http.HandleFunc("/api/circuits", handler.GetCircuits)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
