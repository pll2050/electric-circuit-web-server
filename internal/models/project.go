package models

import (
	"time"
)

// Project represents a circuit design project
type Project struct {
	ID          string                 `json:"id" firestore:"id,omitempty"`
	Name        string                 `json:"name" firestore:"name"`
	Description string                 `json:"description" firestore:"description"`
	UserID      string                 `json:"user_id" firestore:"userId"`
	Status      string                 `json:"status" firestore:"status"` // active, archived, deleted
	Settings    map[string]interface{} `json:"settings" firestore:"settings,omitempty"`
	CreatedAt   time.Time              `json:"created_at" firestore:"createdAt"`
	UpdatedAt   time.Time              `json:"updated_at" firestore:"updatedAt"`
}

// User represents a Firebase user
type User struct {
	UID         string    `json:"uid"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	PhotoURL    string    `json:"photo_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CircuitFirestore represents a circuit in Firestore
type CircuitFirestore struct {
	ID          string                 `json:"id" firestore:"id,omitempty"`
	Name        string                 `json:"name" firestore:"name"`
	Description string                 `json:"description" firestore:"description"`
	ProjectID   string                 `json:"project_id" firestore:"projectId"`
	UserID      string                 `json:"user_id" firestore:"userId"`
	Data        map[string]interface{} `json:"data" firestore:"data"` // Circuit diagram data
	Version     int                    `json:"version" firestore:"version"`
	IsTemplate  bool                   `json:"is_template" firestore:"isTemplate"`
	Tags        []string               `json:"tags" firestore:"tags,omitempty"`
	CreatedAt   time.Time              `json:"created_at" firestore:"createdAt"`
	UpdatedAt   time.Time              `json:"updated_at" firestore:"updatedAt"`
}
