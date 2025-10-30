package controllers

import (
	"context"
)

// CircuitControllerInterface defines the contract for circuit-related operations
type CircuitControllerInterface interface {
	GetProjectCircuits(ctx context.Context, req *CircuitRequest) (*CircuitResponse, error)
	GetCircuit(ctx context.Context, req *CircuitRequest) (*CircuitResponse, error)
	CreateCircuit(ctx context.Context, req *CircuitRequest) (*CircuitResponse, error)
	UpdateCircuit(ctx context.Context, req *CircuitRequest) (*CircuitResponse, error)
	DeleteCircuit(ctx context.Context, req *CircuitRequest) (*CircuitResponse, error)
	CreateFromTemplate(ctx context.Context, req *CircuitRequest) (*CircuitResponse, error)
}

// ProjectControllerInterface defines the contract for project-related operations
type ProjectControllerInterface interface {
	GetUserProjects(ctx context.Context, req *ProjectRequest) (*ProjectResponse, error)
	GetProject(ctx context.Context, req *ProjectRequest) (*ProjectResponse, error)
	CreateProject(ctx context.Context, req *ProjectRequest) (*ProjectResponse, error)
	UpdateProject(ctx context.Context, req *ProjectRequest) (*ProjectResponse, error)
	DeleteProject(ctx context.Context, req *ProjectRequest) (*ProjectResponse, error)
	DuplicateProject(ctx context.Context, req *ProjectRequest) (*ProjectResponse, error)
}

// StorageControllerInterface defines the contract for storage-related operations
type StorageControllerInterface interface {
	UploadFile(ctx context.Context, req *StorageRequest) (*StorageResponse, error)
	GetFileURL(ctx context.Context, req *StorageRequest) (*StorageResponse, error)
	DeleteFile(ctx context.Context, req *StorageRequest) (*StorageResponse, error)
	ListFiles(ctx context.Context, req *StorageRequest) (*StorageResponse, error)
	UploadCircuitImage(ctx context.Context, req *StorageRequest) (*StorageResponse, error)
}

// AuthControllerInterface defines the contract for authentication-related operations
type AuthControllerInterface interface {
	VerifyToken(ctx context.Context, req *AuthRequest) (*AuthResponse, error)
	CreateUser(ctx context.Context, req *AuthRequest) (*AuthResponse, error)
	GetUser(ctx context.Context, req *AuthRequest) (*AuthResponse, error)
	UpdateUser(ctx context.Context, req *AuthRequest) (*AuthResponse, error)
	DeleteUser(ctx context.Context, req *AuthRequest) (*AuthResponse, error)
	SetCustomClaims(ctx context.Context, req *AuthRequest) (*AuthResponse, error)
	Register(ctx context.Context, req *AuthRequest) (*AuthResponse, error)
	ListUsers(ctx context.Context) ([]map[string]interface{}, error)
}
