package handlers

import (
	"electric-circuit-web/server/pkg/database"
	"encoding/json"
	"net/http"
)

type Handler struct {
	db *database.DB
}

func NewHandler(db *database.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"service": "electric-circuit-web",
	})
}

func (h *Handler) GetCircuits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// TODO: Implement circuit retrieval logic
	json.NewEncoder(w).Encode(map[string]interface{}{
		"circuits": []interface{}{},
	})
}
