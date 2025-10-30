package services

// TemplateService handles circuit template business logic
type TemplateService struct {
	// In the future, this could use a repository to load templates from a database
}

// NewTemplateService creates a new template service
func NewTemplateService() *TemplateService {
	return &TemplateService{}
}

// GetTemplateData returns template data by ID
func (s *TemplateService) GetTemplateData(templateID string) map[string]interface{} {
	// Note: In a real implementation, templates would be stored in a database
	// For now, we use hardcoded templates
	templates := map[string]map[string]interface{}{
		"template_basic_circuit": {
			"components": []map[string]interface{}{
				{"id": "battery1", "type": "battery", "x": 100, "y": 100},
				{"id": "resistor1", "type": "resistor", "x": 200, "y": 100},
				{"id": "led1", "type": "led", "x": 300, "y": 100},
			},
			"connections": []map[string]interface{}{
				{"from": "battery1", "to": "resistor1"},
				{"from": "resistor1", "to": "led1"},
			},
		},
		"template_amplifier": {
			"components": []map[string]interface{}{
				{"id": "opamp1", "type": "opamp", "x": 200, "y": 150},
				{"id": "r1", "type": "resistor", "x": 100, "y": 100},
				{"id": "r2", "type": "resistor", "x": 100, "y": 200},
			},
			"connections": []map[string]interface{}{
				{"from": "r1", "to": "opamp1"},
				{"from": "r2", "to": "opamp1"},
			},
		},
	}

	return templates[templateID]
}

// GetAllTemplates returns all available templates
func (s *TemplateService) GetAllTemplates() []map[string]interface{} {
	// Return list of available templates (metadata only)
	return []map[string]interface{}{
		{
			"id":          "template_basic_circuit",
			"name":        "Basic Circuit",
			"description": "A simple circuit with battery, resistor, and LED",
			"category":    "basic",
		},
		{
			"id":          "template_amplifier",
			"name":        "Op-Amp Circuit",
			"description": "An operational amplifier circuit with resistors",
			"category":    "amplifiers",
		},
	}
}
