package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cloud-ru/evo-ai-agents-cli/internal/validator"
)

func init() {
	// Устанавливаем переменные окружения для тестов
	os.Setenv("API_KEY", "test-api-key")
	os.Setenv("PROJECT_ID", "test-project")
}

func TestValidateCommand_Integration(t *testing.T) {
	// Создаем временную директорию для тестов
	tempDir := t.TempDir()

	// Создаем валидный конфигурационный файл
	validConfig := `agents:
  - name: "test-agent"
    description: "Test agent"
    llm_options:
      provider: "openai"
      model: "gpt-4"
`
	validFile := filepath.Join(tempDir, "valid.yaml")
	if err := os.WriteFile(validFile, []byte(validConfig), 0644); err != nil {
		t.Fatalf("Failed to create valid config file: %v", err)
	}

	// Создаем невалидный конфигурационный файл
	invalidConfig := `agents:
  - description: "Test agent without name"
    llm_options:
      provider: "openai"
`
	invalidFile := filepath.Join(tempDir, "invalid.yaml")
	if err := os.WriteFile(invalidFile, []byte(invalidConfig), 0644); err != nil {
		t.Fatalf("Failed to create invalid config file: %v", err)
	}

	// Тестируем валидацию валидного файла
	configValidator := validator.NewConfigValidator()
	schemaPath := "schemas/schema.json"
	for _, name := range []string{"agents", "mcp-servers", "agent-systems"} {
		if err := configValidator.LoadSchema(name, schemaPath); err != nil {
			// Try relative path from cmd package
			schemaPath = "../schemas/schema.json"
			if err2 := configValidator.LoadSchema(name, schemaPath); err2 != nil {
				t.Skipf("Schema file not found, skipping: %v", err)
			}
		}
	}
	result, err := configValidator.ValidateFile(validFile)
	if err != nil {
		t.Fatalf("Unexpected error validating valid file: %v", err)
	}
	if !result.Valid {
		t.Errorf("Expected valid file to be valid, but got errors: %v", result.Errors)
	}

	// Тестируем валидацию невалидного файла
	result, err = configValidator.ValidateFile(invalidFile)
	if err != nil {
		t.Fatalf("Unexpected error validating invalid file: %v", err)
	}
	if result.Valid {
		t.Errorf("Expected invalid file to be invalid, but validation passed")
	}
	if len(result.Errors) == 0 {
		t.Errorf("Expected validation errors for invalid file")
	}
}

func TestFindConfigFiles_Integration(t *testing.T) {
	// Создаем временную директорию с различными файлами
	tempDir := t.TempDir()

	// Создаем структуру директорий
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Создаем файлы различных типов
	files := map[string]string{
		"config.yaml":      "yaml content",
		"config.yml":       "yml content",
		"config.json":      "json content",
		"other.txt":        "text content",
		"subdir/test.yaml": "subdir yaml content",
		"subdir/test.yml":  "subdir yml content",
	}

	for filePath, content := range files {
		fullPath := filepath.Join(tempDir, filePath)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", filePath, err)
		}
	}

	// Тестируем поиск файлов
	foundFiles, err := findConfigFiles(tempDir)
	if err != nil {
		t.Fatalf("Unexpected error finding config files: %v", err)
	}

	// Проверяем, что найдены только конфигурационные файлы
	expectedCount := 4 // 2 yaml, 1 yml, 1 json
	if len(foundFiles) != expectedCount {
		t.Errorf("Expected %d config files, found %d", expectedCount, len(foundFiles))
		for _, file := range foundFiles {
			t.Logf("Found file: %s", file)
		}
	}

	// Проверяем, что txt файл не найден
	for _, file := range foundFiles {
		if filepath.Ext(file) == ".txt" {
			t.Errorf("Found unexpected .txt file: %s", file)
		}
	}
}

func TestValidateCommand_WithExamples(t *testing.T) {
	// Тестируем валидацию примеров из examples/
	examplesDir := "examples"
	if _, err := os.Stat(examplesDir); os.IsNotExist(err) {
		t.Skip("Examples directory not found, skipping test")
	}

	// Находим все конфигурационные файлы в examples/
	files, err := findConfigFiles(examplesDir)
	if err != nil {
		t.Fatalf("Failed to find config files in examples: %v", err)
	}

	if len(files) == 0 {
		t.Skip("No config files found in examples directory")
	}

	// Валидируем каждый файл
	configValidator := validator.NewConfigValidator()
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			result, err := configValidator.ValidateFile(file)
			if err != nil {
				t.Fatalf("Unexpected error validating %s: %v", file, err)
			}

			if !result.Valid {
				t.Errorf("Example file %s is not valid:", file)
				for _, validationError := range result.Errors {
					t.Errorf("  %s", validationError)
				}
			}
		})
	}
}

func TestValidateCommand_EdgeCases(t *testing.T) {
	tempDir := t.TempDir()
	configValidator := validator.NewConfigValidator()

	tests := []struct {
		name        string
		content     string
		filename    string
		expectValid bool
	}{
		{
			name:        "empty file",
			content:     "",
			filename:    "empty.yaml",
			expectValid: false,
		},
		{
			name:        "invalid yaml syntax",
			content:     "invalid: yaml: [",
			filename:    "invalid.yaml",
			expectValid: false,
		},
		{
			name:        "valid json",
			content:     `{"agents": [{"name": "test", "llm_options": {"provider": "openai"}}]}`,
			filename:    "valid.json",
			expectValid: true,
		},
		{
			name:        "invalid json",
			content:     `{"agents": [{"name": "test"}]}`,
			filename:    "invalid.json",
			expectValid: false,
		},
		{
			name: "file with only comments",
			content: `# This is a comment
# Another comment`,
			filename:    "comments.yaml",
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(tempDir, tt.filename)
			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			result, err := configValidator.ValidateFile(filePath)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result.Valid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.Valid)
				if !result.Valid {
					for _, validationError := range result.Errors {
						t.Logf("Validation error: %s", validationError)
					}
				}
			}
		})
	}
}
