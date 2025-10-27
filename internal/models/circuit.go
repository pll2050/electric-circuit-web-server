package models

import "time"

type Circuit struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Data        string    `json:"data"` // JSON string containing circuit diagram data
	UserID      int64     `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Component struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`     // resistor, capacitor, etc.
	Name     string `json:"name"`
	Value    string `json:"value"`
	Standard string `json:"standard"` // KEC, IEC, UL
}
