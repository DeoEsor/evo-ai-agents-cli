package validator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigValidator_ValidateFile(t *testing.T) {
	validator := NewConfigValidator()

	// Загружаем реальную схему из проекта
	schemaPath := "../../schemas/schema.json"
	for _, name := range []string{"agents", "mcp-servers", "agent-systems"} {
		if err := validator.LoadSchema(name, schemaPath); err != nil {
			t.Skipf("Schema file not found at %s, skipping test: %v", schemaPath, err)
		}
	}

	// Создаем временные файлы для тестирования
	tempDir := t.TempDir()

	tests := []struct {
		name          string
		fileContent   string
		fileName      string
		expectValid   bool
		expectErrors  int
		errorContains []string
	}{
		{
			name: "valid agents yaml",
			fileContent: `agents:
  - name: "test-agent"
    description: "Test agent"
    llm_options:
      provider: "openai"
      model: "gpt-4"
`,
			fileName:     "valid_agents.yaml",
			expectValid:  true,
			expectErrors: 0,
		},
		{
			name: "valid mcp-servers yaml",
			fileContent: `mcp-servers:
  - name: "test-mcp"
    description: "Test MCP server"
    options:
      host: "localhost"
      port: 8080
`,
			fileName:     "valid_mcp.yaml",
			expectValid:  true,
			expectErrors: 0,
		},
		{
			name: "valid agent-systems yaml",
			fileContent: `agent-systems:
  - name: "test-system"
    description: "Test system"
    agents:
      - "agent1"
      - "agent2"
`,
			fileName:     "valid_systems.yaml",
			expectValid:  true,
			expectErrors: 0,
		},
		{
			name: "agent without name passes permissive schema (YAML)",
			fileContent: `agents:
  - description: "Test agent"
    llm_options:
      provider: "openai"
`,
			fileName:     "agent_no_name.yaml",
			expectValid:  true,
			expectErrors: 0,
		},
		{
			name: "mcp-server without name passes permissive schema (YAML)",
			fileContent: `mcp-servers:
  - description: "Test MCP server"
    options:
      host: "localhost"
`,
			fileName:     "mcp_no_name.yaml",
			expectValid:  true,
			expectErrors: 0,
		},
		{
			name: "agent-system without agents passes permissive schema (YAML)",
			fileContent: `agent-systems:
  - name: "test-system"
    description: "Test system"
`,
			fileName:     "system_no_agents.yaml",
			expectValid:  true,
			expectErrors: 0,
		},
		{
			name:         "agent without name passes permissive schema (JSON)",
			fileContent:  `{"agents": [{"description": "Test agent", "llm_options": {"provider": "openai"}}]}`,
			fileName:     "agent_no_name.json",
			expectValid:  true,
			expectErrors: 0,
		},
		{
			name:         "mcp-server without name passes permissive schema (JSON)",
			fileContent:  `{"mcp-servers": [{"description": "Test MCP server", "options": {"host": "localhost"}}]}`,
			fileName:     "mcp_no_name.json",
			expectValid:  true,
			expectErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(tempDir, tt.fileName)
			err := os.WriteFile(filePath, []byte(tt.fileContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			result, err := validator.ValidateFile(filePath)
			if err != nil {
				if !tt.expectValid {
					return
				}
				t.Fatalf("Unexpected error during validation: %v", err)
			}

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
			}

			if len(result.Errors) != tt.expectErrors {
				t.Errorf("Expected %d errors, got %d", tt.expectErrors, len(result.Errors))
				for i, err := range result.Errors {
					t.Logf("Error %d: %s", i, err)
				}
			}

			if len(tt.errorContains) > 0 {
				found := false
				for _, expectedError := range tt.errorContains {
					for _, actualError := range result.Errors {
						if contains(actualError, expectedError) {
							found = true
							break
						}
					}
					if found {
						break
					}
				}
				if !found {
					t.Errorf("Expected error to contain one of: %v, but got errors: %v", tt.errorContains, result.Errors)
				}
			}
		})
	}
}


func TestConfigValidator_LoadSchema(t *testing.T) {
	validator := NewConfigValidator()

	// Создаем временный файл схемы
	tempDir := t.TempDir()
	schemaFile := filepath.Join(tempDir, "test_schema.json")
	schemaContent := `{
		"type": "object",
		"properties": {
			"name": {
				"type": "string"
			}
		}
	}`

	err := os.WriteFile(schemaFile, []byte(schemaContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create schema file: %v", err)
	}

	// Тестируем загрузку схемы
	err = validator.LoadSchema("test", schemaFile)
	if err != nil {
		t.Errorf("Failed to load schema: %v", err)
	}

	// Проверяем, что схема загружена
	if validator.schemas["test"] == nil {
		t.Errorf("Schema not loaded")
	}

	// Тестируем загрузку несуществующего файла
	err = validator.LoadSchema("nonexistent", "nonexistent.json")
	if err == nil {
		t.Errorf("Expected error for nonexistent file")
	}
}


// Вспомогательная функция для проверки содержания строки
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && s[len(s)-len(substr):] == substr ||
		len(s) >= len(substr) && containsSubstring(s, substr)
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
