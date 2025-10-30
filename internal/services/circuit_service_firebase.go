package services

import (
	"electric-circuit-web/server/internal/models"
	"electric-circuit-web/server/internal/repositories"
	"fmt"
)

// FirebaseCircuitService handles circuit-related business logic for Firebase
type FirebaseCircuitService struct {
	firestoreRepository *repositories.FirestoreService
	projectService      *ProjectService
}

// NewFirebaseCircuitService creates a new Firebase circuit service
func NewFirebaseCircuitService(firestoreRepository *repositories.FirestoreService, projectService *ProjectService) *FirebaseCircuitService {
	return &FirebaseCircuitService{
		firestoreRepository: firestoreRepository,
		projectService:      projectService,
	}
}

// GetProjectCircuits retrieves circuits for a specific project
func (s *FirebaseCircuitService) GetProjectCircuits(projectID, userID string) ([]models.CircuitFirestore, error) {
	// Verify user owns the project
	_, err := s.projectService.GetProject(projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify project ownership: %v", err)
	}

	docs, err := s.firestoreRepository.QueryCollection("circuits", "projectId", "==", projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get project circuits: %v", err)
	}

	circuits := make([]models.CircuitFirestore, 0, len(docs))
	for _, doc := range docs {
		circuit, err := s.mapToCircuit(doc)
		if err != nil {
			continue // Skip invalid circuits
		}
		circuits = append(circuits, *circuit)
	}

	return circuits, nil
}

// CreateCircuit creates a new circuit
func (s *FirebaseCircuitService) CreateCircuit(userID string, circuit *models.CircuitFirestore) (string, error) {
	// Note: Input validation is done in Handler and Controller layers
	// Service focuses on business logic

	// Business logic: Verify user owns the project
	_, err := s.projectService.GetProject(circuit.ProjectID, userID)
	if err != nil {
		return "", fmt.Errorf("failed to verify project ownership: %v", err)
	}

	// Set user ID and default values
	circuit.UserID = userID
	if circuit.Version == 0 {
		circuit.Version = 1
	}

	// Convert to map for Firestore
	data := s.circuitToMap(circuit)

	// Create in Firestore
	circuitID, err := s.firestoreRepository.CreateDocumentWithAutoID("circuits", data)
	if err != nil {
		return "", fmt.Errorf("failed to create circuit: %v", err)
	}

	return circuitID, nil
}

// GetCircuit retrieves a circuit by ID
func (s *FirebaseCircuitService) GetCircuit(circuitID, userID string) (*models.CircuitFirestore, error) {
	doc, err := s.firestoreRepository.GetDocument("circuits", circuitID)
	if err != nil {
		return nil, fmt.Errorf("failed to get circuit: %v", err)
	}

	circuit, err := s.mapToCircuit(doc)
	if err != nil {
		return nil, err
	}

	// Verify user owns the project that contains this circuit
	_, err = s.projectService.GetProject(circuit.ProjectID, userID)
	if err != nil {
		return nil, fmt.Errorf("access denied: circuit does not belong to user")
	}

	return circuit, nil
}

// UpdateCircuit updates an existing circuit
func (s *FirebaseCircuitService) UpdateCircuit(circuitID, userID string, updates map[string]interface{}) error {
	// Verify ownership first
	_, err := s.GetCircuit(circuitID, userID)
	if err != nil {
		return err
	}

	// Remove fields that shouldn't be updated
	delete(updates, "userId")
	delete(updates, "projectId")
	delete(updates, "id")
	delete(updates, "createdAt")

	// Increment version if data is being updated
	if _, hasData := updates["data"]; hasData {
		circuit, _ := s.GetCircuit(circuitID, userID)
		if circuit != nil {
			updates["version"] = circuit.Version + 1
		}
	}

	err = s.firestoreRepository.UpdateDocument("circuits", circuitID, updates)
	if err != nil {
		return fmt.Errorf("failed to update circuit: %v", err)
	}

	return nil
}

// DeleteCircuit deletes a circuit
func (s *FirebaseCircuitService) DeleteCircuit(circuitID, userID string) error {
	// Verify ownership first
	_, err := s.GetCircuit(circuitID, userID)
	if err != nil {
		return err
	}

	err = s.firestoreRepository.DeleteDocument("circuits", circuitID)
	if err != nil {
		return fmt.Errorf("failed to delete circuit: %v", err)
	}

	return nil
}

// Helper methods
func (s *FirebaseCircuitService) mapToCircuit(doc map[string]interface{}) (*models.CircuitFirestore, error) {
	circuit := &models.CircuitFirestore{}

	if id, ok := doc["id"].(string); ok {
		circuit.ID = id
	}

	if name, ok := doc["name"].(string); ok {
		circuit.Name = name
	} else {
		return nil, fmt.Errorf("invalid circuit: name is required")
	}

	if description, ok := doc["description"].(string); ok {
		circuit.Description = description
	}

	if projectID, ok := doc["projectId"].(string); ok {
		circuit.ProjectID = projectID
	} else {
		return nil, fmt.Errorf("invalid circuit: projectId is required")
	}

	if userID, ok := doc["userId"].(string); ok {
		circuit.UserID = userID
	} else {
		return nil, fmt.Errorf("invalid circuit: userId is required")
	}

	if data, ok := doc["data"].(map[string]interface{}); ok {
		circuit.Data = data
	}

	if version, ok := doc["version"].(int); ok {
		circuit.Version = version
	}

	if isTemplate, ok := doc["isTemplate"].(bool); ok {
		circuit.IsTemplate = isTemplate
	}

	if tags, ok := doc["tags"].([]interface{}); ok {
		circuit.Tags = make([]string, len(tags))
		for i, tag := range tags {
			if tagStr, ok := tag.(string); ok {
				circuit.Tags[i] = tagStr
			}
		}
	}

	return circuit, nil
}

func (s *FirebaseCircuitService) circuitToMap(circuit *models.CircuitFirestore) map[string]interface{} {
	return map[string]interface{}{
		"name":        circuit.Name,
		"description": circuit.Description,
		"projectId":   circuit.ProjectID,
		"userId":      circuit.UserID,
		"data":        circuit.Data,
		"version":     circuit.Version,
		"isTemplate":  circuit.IsTemplate,
		"tags":        circuit.Tags,
	}
}
