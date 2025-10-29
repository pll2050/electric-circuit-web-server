package services

import (
	"electric-circuit-web/server/internal/models"
	"electric-circuit-web/server/internal/repositories"
	"fmt"
)

// ProjectService handles project-related business logic
type ProjectService struct {
	firestoreRepository *repositories.FirestoreService
}

// NewProjectService creates a new project service
func NewProjectService(firestoreRepository *repositories.FirestoreService) *ProjectService {
	return &ProjectService{
		firestoreRepository: firestoreRepository,
	}
}

// GetUserProjects retrieves projects for a specific user
func (s *ProjectService) GetUserProjects(userID string) ([]models.Project, error) {
	docs, err := s.firestoreRepository.QueryCollection("projects", "userId", "==", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user projects: %v", err)
	}

	projects := make([]models.Project, 0, len(docs))
	for _, doc := range docs {
		project, err := s.mapToProject(doc)
		if err != nil {
			continue // Skip invalid projects
		}
		projects = append(projects, *project)
	}

	return projects, nil
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(userID string, project *models.Project) (string, error) {
	// Validate project data
	if project.Name == "" {
		return "", fmt.Errorf("project name is required")
	}

	// Set user ID and default values
	project.UserID = userID
	if project.Status == "" {
		project.Status = "active"
	}

	// Convert to map for Firestore
	data := s.projectToMap(project)

	// Create in Firestore
	projectID, err := s.firestoreRepository.CreateDocumentWithAutoID("projects", data)
	if err != nil {
		return "", fmt.Errorf("failed to create project: %v", err)
	}

	return projectID, nil
}

// GetProject retrieves a project by ID
func (s *ProjectService) GetProject(projectID, userID string) (*models.Project, error) {
	doc, err := s.firestoreRepository.GetDocument("projects", projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %v", err)
	}

	// Check ownership
	if doc["userId"] != userID {
		return nil, fmt.Errorf("access denied: project does not belong to user")
	}

	return s.mapToProject(doc)
}

// UpdateProject updates an existing project
func (s *ProjectService) UpdateProject(projectID, userID string, updates map[string]interface{}) error {
	// Verify ownership first
	_, err := s.GetProject(projectID, userID)
	if err != nil {
		return err
	}

	// Remove fields that shouldn't be updated
	delete(updates, "userId")
	delete(updates, "id")
	delete(updates, "createdAt")

	err = s.firestoreRepository.UpdateDocument("projects", projectID, updates)
	if err != nil {
		return fmt.Errorf("failed to update project: %v", err)
	}

	return nil
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(projectID, userID string) error {
	// Verify ownership first
	_, err := s.GetProject(projectID, userID)
	if err != nil {
		return err
	}

	// TODO: Also delete related circuits
	err = s.firestoreRepository.DeleteDocument("projects", projectID)
	if err != nil {
		return fmt.Errorf("failed to delete project: %v", err)
	}

	return nil
}

// Helper methods
func (s *ProjectService) mapToProject(doc map[string]interface{}) (*models.Project, error) {
	project := &models.Project{}

	if id, ok := doc["id"].(string); ok {
		project.ID = id
	}

	if name, ok := doc["name"].(string); ok {
		project.Name = name
	} else {
		return nil, fmt.Errorf("invalid project: name is required")
	}

	if description, ok := doc["description"].(string); ok {
		project.Description = description
	}

	if userID, ok := doc["userId"].(string); ok {
		project.UserID = userID
	} else {
		return nil, fmt.Errorf("invalid project: userId is required")
	}

	if status, ok := doc["status"].(string); ok {
		project.Status = status
	}

	// Handle optional fields
	if settings, ok := doc["settings"].(map[string]interface{}); ok {
		project.Settings = settings
	}

	return project, nil
}

func (s *ProjectService) projectToMap(project *models.Project) map[string]interface{} {
	return map[string]interface{}{
		"name":        project.Name,
		"description": project.Description,
		"userId":      project.UserID,
		"status":      project.Status,
		"settings":    project.Settings,
	}
}
