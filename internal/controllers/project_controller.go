package controllers

import (
	"context"

	"electric-circuit-web/server/internal/models"
	"electric-circuit-web/server/internal/services"
)

// ProjectController handles project-related business coordination
type ProjectController struct {
	projectService *services.ProjectService
}

// NewProjectController creates a new project controller
func NewProjectController(projectService *services.ProjectService) *ProjectController {
	return &ProjectController{
		projectService: projectService,
	}
}

// ProjectRequest represents a project operation request
type ProjectRequest struct {
	UserID      string `json:"user_id"`
	ProjectID   string `json:"project_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ProjectResponse represents a project operation response
type ProjectResponse struct {
	Success  bool             `json:"success"`
	Message  string           `json:"message"`
	Project  *models.Project  `json:"project,omitempty"`
	Projects []models.Project `json:"projects,omitempty"`
	Error    string           `json:"error,omitempty"`
}

// GetUserProjects retrieves all projects for a user
func (c *ProjectController) GetUserProjects(ctx context.Context, req *ProjectRequest) (*ProjectResponse, error) {
	// Validate request
	if req.UserID == "" {
		return &ProjectResponse{
			Success: false,
			Error:   "User ID is required",
		}, nil
	}

	// Call service
	projects, err := c.projectService.GetUserProjects(req.UserID)
	if err != nil {
		return &ProjectResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &ProjectResponse{
		Success:  true,
		Message:  "Projects retrieved successfully",
		Projects: projects,
	}, nil
}

// GetProject retrieves a specific project
func (c *ProjectController) GetProject(ctx context.Context, req *ProjectRequest) (*ProjectResponse, error) {
	// Validate request
	if req.ProjectID == "" {
		return &ProjectResponse{
			Success: false,
			Error:   "Project ID is required",
		}, nil
	}

	if req.UserID == "" {
		return &ProjectResponse{
			Success: false,
			Error:   "User ID is required",
		}, nil
	}

	// Call service
	project, err := c.projectService.GetProject(req.ProjectID, req.UserID)
	if err != nil {
		return &ProjectResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &ProjectResponse{
		Success: true,
		Message: "Project retrieved successfully",
		Project: project,
	}, nil
}

// CreateProject creates a new project
func (c *ProjectController) CreateProject(ctx context.Context, req *ProjectRequest) (*ProjectResponse, error) {
	// Validate request
	if req.Name == "" {
		return &ProjectResponse{
			Success: false,
			Error:   "Project name is required",
		}, nil
	}

	if req.UserID == "" {
		return &ProjectResponse{
			Success: false,
			Error:   "User ID is required",
		}, nil
	}

	// Create project model
	project := &models.Project{
		Name:        req.Name,
		Description: req.Description,
		UserID:      req.UserID,
	}

	// Call service
	projectID, err := c.projectService.CreateProject(req.UserID, project)
	if err != nil {
		return &ProjectResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// Set ID for response
	project.ID = projectID

	return &ProjectResponse{
		Success: true,
		Message: "Project created successfully",
		Project: project,
	}, nil
}

// UpdateProject updates an existing project
func (c *ProjectController) UpdateProject(ctx context.Context, req *ProjectRequest) (*ProjectResponse, error) {
	// Validate request
	if req.ProjectID == "" {
		return &ProjectResponse{
			Success: false,
			Error:   "Project ID is required",
		}, nil
	}

	if req.UserID == "" {
		return &ProjectResponse{
			Success: false,
			Error:   "User ID is required",
		}, nil
	}

	// Prepare updates
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	// Call service
	err := c.projectService.UpdateProject(req.ProjectID, req.UserID, updates)
	if err != nil {
		return &ProjectResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &ProjectResponse{
		Success: true,
		Message: "Project updated successfully",
	}, nil
}

// DeleteProject deletes a project
func (c *ProjectController) DeleteProject(ctx context.Context, req *ProjectRequest) (*ProjectResponse, error) {
	// Validate request
	if req.ProjectID == "" {
		return &ProjectResponse{
			Success: false,
			Error:   "Project ID is required",
		}, nil
	}

	if req.UserID == "" {
		return &ProjectResponse{
			Success: false,
			Error:   "User ID is required",
		}, nil
	}

	// Call service
	err := c.projectService.DeleteProject(req.ProjectID, req.UserID)
	if err != nil {
		return &ProjectResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	return &ProjectResponse{
		Success: true,
		Message: "Project deleted successfully",
	}, nil
}

// DuplicateProject duplicates an existing project
func (c *ProjectController) DuplicateProject(ctx context.Context, req *ProjectRequest) (*ProjectResponse, error) {
	// Validate request
	if req.ProjectID == "" {
		return &ProjectResponse{
			Success: false,
			Error:   "Project ID is required",
		}, nil
	}

	if req.UserID == "" {
		return &ProjectResponse{
			Success: false,
			Error:   "User ID is required",
		}, nil
	}

	// Get original project
	originalProject, err := c.projectService.GetProject(req.ProjectID, req.UserID)
	if err != nil {
		return &ProjectResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// Create new project with duplicated data
	newName := req.Name
	if newName == "" {
		newName = originalProject.Name + " (Copy)"
	}

	duplicatedProject := &models.Project{
		Name:        newName,
		Description: originalProject.Description,
		UserID:      req.UserID,
	}

	// Call service
	projectID, err := c.projectService.CreateProject(req.UserID, duplicatedProject)
	if err != nil {
		return &ProjectResponse{
			Success: false,
			Error:   err.Error(),
		}, err
	}

	// Set ID for response
	duplicatedProject.ID = projectID

	return &ProjectResponse{
		Success: true,
		Message: "Project duplicated successfully",
		Project: duplicatedProject,
	}, nil
}
