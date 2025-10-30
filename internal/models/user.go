package models

import "time"

// User represents a user in both Firebase and PostgreSQL database
type User struct {
	ID          string    `json:"id,omitempty" db:"id"`           // Firebase UID (primary key in PostgreSQL)
	UID         string    `json:"uid,omitempty"`                  // Firebase UID alias for compatibility
	Email       string    `json:"email" db:"email"`               // User email
	DisplayName string    `json:"display_name" db:"display_name"` // User display name
	PhotoURL    string    `json:"photo_url" db:"photo_url"`       // User photo URL
	Provider    string    `json:"provider,omitempty" db:"provider"` // Auth provider (google, email, etc.)
	CreatedAt   time.Time `json:"created_at" db:"created_at"`     // Creation timestamp
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`     // Last update timestamp
}
