package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"electric-circuit-web/server/internal/models"
)

// UserRepository handles user data operations in PostgreSQL
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create inserts a new user into PostgreSQL
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (id, email, display_name, photo_url, provider, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	now := time.Now()
	_, err := r.db.Exec(
		query,
		user.ID,
		user.Email,
		user.DisplayName,
		user.PhotoURL,
		user.Provider,
		now,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to create user in PostgreSQL: %v", err)
	}

	user.CreatedAt = now
	user.UpdatedAt = now

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id string) (*models.User, error) {
	query := `
		SELECT id, email, display_name, photo_url, provider, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PhotoURL,
		&user.Provider,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, display_name, photo_url, provider, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PhotoURL,
		&user.Provider,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	return user, nil
}

// Update updates user information
func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET display_name = $2, photo_url = $3, updated_at = $4
		WHERE id = $1
	`

	now := time.Now()
	_, err := r.db.Exec(
		query,
		user.ID,
		user.DisplayName,
		user.PhotoURL,
		now,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	user.UpdatedAt = now

	return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}

// Exists checks if a user exists by ID
func (r *UserRepository) Exists(id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`

	var exists bool
	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %v", err)
	}

	return exists, nil
}

// GetAll retrieves all users from the database
func (r *UserRepository) GetAll() ([]*models.User, error) {
	query := `
		SELECT id, email, display_name, photo_url, provider, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %v", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.DisplayName,
			&user.PhotoURL,
			&user.Provider,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating users: %v", err)
	}

	return users, nil
}
