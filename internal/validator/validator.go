package validator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

// ConfigValidator представляет валидатор конфигурации
type ConfigValidator struct {
	schemas map[string]map[string]interface{}
}

// NewConfigValidator создает новый валидатор конфигурации
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{
		schemas: make(map[string]map[string]interface{}),
	}
}

// LoadSchema загружает схему из файла
func (v *ConfigValidator) LoadSchema(name, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read schema file %s: %w", filePath, err)
	}

	var schema map[string]interface{}
	if err := json.Unmarshal(data, &schema); err != nil {
		return fmt.Errorf("failed to parse schema %s: %w", filePath, err)
	}

	v.schemas[name] = schema
	return nil
}

// ValidationResult представляет результат валидации
type ValidationResult struct {
	Valid  bool
	Errors []string
}

// parseFileContent reads a file and parses it as YAML or JSON depending on extension.
func parseFileContent(filePath string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var config map[string]interface{}

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON configuration file: %w", err)
		}
	default:
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML configuration file: %w", err)
		}
	}

	return config, nil
}

// detectSchemaName determines the schema name from config content.
// Returns one of: "agents", "mcp-servers", "agent-systems".
func detectSchemaName(config map[string]interface{}) (string, error) {
	if _, hasAgents := config["agents"]; hasAgents {
		return "agents", nil
	}
	if _, hasMCPServers := config["mcp-servers"]; hasMCPServers {
		return "mcp-servers", nil
	}
	if _, hasSystems := config["agent-systems"]; hasSystems {
		return "agent-systems", nil
	}
	return "", fmt.Errorf("unknown configuration type: file must contain 'agents', 'mcp-servers', or 'agent-systems' key")
}

// ValidateFileWithResult валидирует файл и возвращает результат
func (v *ConfigValidator) ValidateFileWithResult(filePath string) (*ValidationResult, error) {
	config, err := parseFileContent(filePath)
	if err != nil {
		return nil, err
	}

	if config == nil {
		return nil, fmt.Errorf("configuration file is empty")
	}

	schemaName, err := detectSchemaName(config)
	if err != nil {
		return nil, err
	}

	schema, exists := v.schemas[schemaName]
	if !exists {
		return nil, fmt.Errorf("schema %s not found", schemaName)
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	schemaLoader := gojsonschema.NewGoLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(configJSON)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return nil, fmt.Errorf("validation error: %w", err)
	}

	var errors []string
	if !result.Valid() {
		for _, desc := range result.Errors() {
			errors = append(errors, desc.String())
		}
	}

	return &ValidationResult{
		Valid:  result.Valid(),
		Errors: errors,
	}, nil
}

// ValidateFile валидирует файл (старый API для совместимости)
func (v *ConfigValidator) ValidateFile(filePath string) (*ValidationResult, error) {
	return v.ValidateFileWithResult(filePath)
}

// PrintErrors выводит ошибки валидации
func (v *ConfigValidator) PrintErrors(result *ValidationResult) {
	if result.Valid {
		fmt.Println("✅ Файл валиден")
		return
	}

	fmt.Println("❌ Ошибки валидации:")
	for i, err := range result.Errors {
		fmt.Printf("  %d. %s\n", i+1, err)
	}
}

// ValidateConfig валидирует конфигурацию по JSON схеме
func ValidateConfig(config map[string]interface{}, schemaPath string) error {
	schemaData, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file %s: %w", schemaPath, err)
	}

	var schema map[string]interface{}
	if err := json.Unmarshal(schemaData, &schema); err != nil {
		return fmt.Errorf("failed to parse schema: %w", err)
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config to JSON: %w", err)
	}

	schemaLoader := gojsonschema.NewGoLoader(schema)
	documentLoader := gojsonschema.NewBytesLoader(configJSON)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if !result.Valid() {
		var errors string
		for _, desc := range result.Errors() {
			errors += fmt.Sprintf("- %s\n", desc)
		}
		return fmt.Errorf("validation failed:\n%s", errors)
	}

	return nil
}

// ValidateFile валидирует файл конфигурации
func ValidateFile(filePath string, schemaPath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file %s does not exist", filePath)
	}

	config, err := parseFileContent(filePath)
	if err != nil {
		return err
	}

	return ValidateConfig(config, schemaPath)
}
