package scaffolder

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
)

// Use TemplatesFS from templates.go

// Scaffolder represents the project scaffolding functionality
type Scaffolder struct {
	templates embed.FS
	config    *ScaffolderConfig
}

// ScaffolderConfig holds configuration for the scaffolder
type ScaffolderConfig struct {
	Author      string
	DefaultCICD string
}

// ProjectData holds the data for template rendering
type ProjectData struct {
	ProjectName     string
	ProjectType     string
	Framework       string // New field for agent framework
	Author          string
	Year            string
	CICDType        string
	DatabaseType    string // New field for database selection
	ExternalAPIKeys string // New field for external API keys selection
	Description     string // New field for project description
}

// NewScaffolder creates a new scaffolder instance
func NewScaffolder() *Scaffolder {
	config := &ScaffolderConfig{
		Author:      "", // Will be determined dynamically from git config
		DefaultCICD: getEnvOrDefault("SCAFFOLDER_DEFAULT_CICD", "both"),
	}

	return &Scaffolder{
		templates: TemplatesFS,
		config:    config,
	}
}

// NewScaffolderWithConfig creates a new scaffolder instance with custom config
func NewScaffolderWithConfig(config *ScaffolderConfig) *Scaffolder {
	if config == nil {
		config = &ScaffolderConfig{
			Author:      "", // Will be determined dynamically from git config
			DefaultCICD: "both",
		}
	}
	return &Scaffolder{
		templates: TemplatesFS,
		config:    config,
	}
}

// CreateProject creates a new project from template
func (s *Scaffolder) CreateProject(projectType, projectName, targetPath, cicdType string) error {
	// log.Info("Creating project", "type", projectType, "name", projectName, "path", targetPath)

	// Validate inputs
	if err := s.validateInputs(projectType, projectName, targetPath); err != nil {
		return errors.Wrap(err, errors.ErrorTypeValidation, errors.SeverityMedium, "VALIDATION_FAILED", fmt.Sprintf("Ошибка валидации входных данных: %v", err))
	}

	// Prepare template data
	data := &ProjectData{
		ProjectName: projectName,
		ProjectType: projectType,
		Author:      s.getAuthor(),
		Year:        fmt.Sprintf("%d", time.Now().Year()),
		CICDType:    cicdType,
	}

	// Get template directory; for agents without a framework, default to adk
	var templateDir string
	if projectType == "agent" {
		templateDir = filepath.Join("templates", "agent-frameworks", "adk")
	} else {
		templateDir = filepath.Join("templates", projectType)
	}

	// Create target directory
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return errors.Wrap(err, errors.ErrorTypeFileSystem, errors.SeverityMedium, "DIRECTORY_CREATION_FAILED", "Ошибка создания директории")
	}

	// Process template files
	if err := s.processTemplateFiles(templateDir, targetPath, data); err != nil {
		return errors.Wrap(err, errors.ErrorTypeTemplate, errors.SeverityMedium, "TEMPLATE_PROCESSING_FAILED", fmt.Sprintf("Ошибка обработки шаблонов: %v", err))
	}

	// log.Info("Project created successfully", "path", targetPath)
	return nil
}

// CreateProjectWithOptions creates a new project from template with additional options
func (s *Scaffolder) CreateProjectWithOptions(projectType, projectName, targetPath, cicdType, framework, databaseType, externalAPIKeys string, options []string) error {
	// log.Info("Creating project", "type", projectType, "name", projectName, "path", targetPath)

	// Validate inputs
	if err := s.validateInputs(projectType, projectName, targetPath); err != nil {
		return errors.Wrap(err, errors.ErrorTypeValidation, errors.SeverityMedium, "VALIDATION_FAILED", fmt.Sprintf("Ошибка валидации входных данных: %v", err))
	}

	// Prepare template data
	data := &ProjectData{
		ProjectName:     projectName,
		ProjectType:     projectType,
		Framework:       framework,
		Author:          s.getAuthor(),
		Year:            fmt.Sprintf("%d", time.Now().Year()),
		CICDType:        cicdType,
		DatabaseType:    databaseType,
		ExternalAPIKeys: externalAPIKeys,
		Description:     s.getProjectDescription(projectType, framework),
	}

	// Get template directory based on project type and framework
	var templateDir string
	if projectType == "agent" && framework != "" {
		templateDir = filepath.Join("templates", "agent-frameworks", framework)
	} else {
		templateDir = filepath.Join("templates", projectType)
	}

	// Create target directory
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return errors.Wrap(err, errors.ErrorTypeFileSystem, errors.SeverityMedium, "DIRECTORY_CREATION_FAILED", "Ошибка создания директории")
	}

	// Cleanup function to remove project directory on error
	cleanup := func() {
		if err := os.RemoveAll(targetPath); err != nil {
			// Log cleanup error but don't fail the main operation
			fmt.Printf("Warning: failed to cleanup project directory %s: %v\n", targetPath, err)
		}
	}

	// Process template files
	if err := s.processTemplateFiles(templateDir, targetPath, data); err != nil {
		cleanup()
		return errors.Wrap(err, errors.ErrorTypeTemplate, errors.SeverityMedium, "TEMPLATE_PROCESSING_FAILED", fmt.Sprintf("Ошибка обработки шаблонов: %v", err))
	}

	// Apply additional options
	if err := s.applyProjectOptions(targetPath, options); err != nil {
		cleanup()
		return fmt.Errorf("failed to apply project options: %w", err)
	}

	// log.Info("Project created successfully", "path", targetPath)
	return nil
}

// validateInputs validates the input parameters
func (s *Scaffolder) validateInputs(projectType, projectName, targetPath string) error {
	// Validate project type
	if projectType != "mcp" && projectType != "agent" {
		return fmt.Errorf("invalid project type: %s (must be 'mcp' or 'agent')", projectType)
	}

	// Validate project name
	if projectName == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Check if target path already exists
	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("target path already exists: %s", targetPath)
	}

	return nil
}

// processTemplateFiles processes all template files in the template directory
func (s *Scaffolder) processTemplateFiles(templateDir, targetPath string, data *ProjectData) error {
	// Debug: check if template directory exists
	if _, err := s.templates.ReadDir(templateDir); err != nil {
		return fmt.Errorf("template directory '%s' not found in embedded FS: %w", templateDir, err)
	}

	return fs.WalkDir(s.templates, templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the template directory itself
		if path == templateDir {
			return nil
		}

		// Calculate relative path within template directory
		relPath, err := filepath.Rel(templateDir, path)
		if err != nil {
			return err
		}

		// Skip hidden files and directories, but allow specific ones
		if strings.HasPrefix(relPath, ".") &&
			relPath != ".gitignore" &&
			relPath != ".editorconfig" &&
			relPath != ".env.example" &&
			relPath != ".gitlab-ci.yml" &&
			!strings.HasPrefix(relPath, ".github/") {
			return nil
		}

		// Skip CI/CD files based on CICDType
		if relPath == ".gitlab-ci.yml" && data.CICDType != "gitlab" && data.CICDType != "both" {
			return nil
		}
		if strings.HasPrefix(relPath, ".github/") && data.CICDType != "github" && data.CICDType != "both" {
			return nil
		}

		if d.IsDir() {
			// Create directory
			dirPath := filepath.Join(targetPath, s.processPath(relPath, data))
			return os.MkdirAll(dirPath, 0755)
		} else {
			// Process file
			return s.processFile(path, targetPath, relPath, data)
		}
	})
}

// processFile processes a single template file
func (s *Scaffolder) processFile(templatePath, targetPath, relPath string, data *ProjectData) error {
	// Read template content
	content, err := s.templates.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", templatePath, err)
	}

	// Process path (remove .tmpl extension)
	processedPath := s.processPath(relPath, data)

	// Remove .tmpl extension from filename
	processedPath = strings.TrimSuffix(processedPath, ".tmpl")

	// Create target file path
	targetFilePath := filepath.Join(targetPath, processedPath)

	// Create directory if needed
	if err := os.MkdirAll(filepath.Dir(targetFilePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", targetFilePath, err)
	}

	// Process template content
	processedContent, err := s.processTemplate(string(content), data)
	if err != nil {
		return fmt.Errorf("failed to process template content for %s: %w", templatePath, err)
	}

	// Write processed content to target file
	if err := os.WriteFile(targetFilePath, []byte(processedContent), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", targetFilePath, err)
	}

	// log.Debug("Processed template file", "template", templatePath, "target", targetFilePath)
	return nil
}

// processPath processes a file path with template variables
func (s *Scaffolder) processPath(path string, data *ProjectData) string {
	// For now, just return the path as-is
	// In the future, we could add path template processing here
	return path
}

// processTemplate processes template content with data
func (s *Scaffolder) processTemplate(content string, data *ProjectData) (string, error) {
	tmpl, err := template.New("template").Parse(content)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return result.String(), nil
}

// GetAvailableTemplates returns the list of available project templates
func (s *Scaffolder) GetAvailableTemplates() ([]string, error) {
	seen := make(map[string]bool)
	var templates []string

	err := fs.WalkDir(s.templates, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && path != "templates" {
			parts := strings.Split(path, "/")
			if len(parts) >= 2 {
				templateName := parts[1]
				if !seen[templateName] {
					seen[templateName] = true
					templates = append(templates, templateName)
				}
				// Expose "agent" as a virtual template backed by agent-frameworks
				if templateName == "agent-frameworks" && !seen["agent"] {
					seen["agent"] = true
					templates = append(templates, "agent")
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk templates directory: %w", err)
	}

	return templates, nil
}

// ValidateTemplate validates that a template exists
func (s *Scaffolder) ValidateTemplate(templateName string) error {
	templatePath := filepath.Join("templates", templateName)

	if _, err := s.templates.ReadDir(templatePath); err != nil {
		// For agent type, try agent-frameworks/adk as default
		if templateName == "agent" {
			altPath := filepath.Join("templates", "agent-frameworks", "adk")
			if _, err2 := s.templates.ReadDir(altPath); err2 == nil {
				return nil
			}
		}
		return fmt.Errorf("template '%s' not found", templateName)
	}

	return nil
}

// getAuthor returns the configured author
func (s *Scaffolder) getAuthor() string {
	// Priority: config > git config > environment > default
	if s.config != nil && s.config.Author != "" {
		return s.config.Author
	}

	// Try to get from git config
	if gitAuthor := getGitAuthor(); gitAuthor != "" {
		return gitAuthor
	}

	// Try to get from environment variable
	if envAuthor := getEnvOrDefault("SCAFFOLDER_AUTHOR", ""); envAuthor != "" {
		return envAuthor
	}

	// Default fallback
	return "Cloud.ru Team"
}

// getDefaultCICD returns the configured default CI/CD type
func (s *Scaffolder) getDefaultCICD() string {
	if s.config != nil && s.config.DefaultCICD != "" {
		return s.config.DefaultCICD
	}
	return "both"
}

// getProjectDescription returns a description based on project type and framework
func (s *Scaffolder) getProjectDescription(projectType, framework string) string {
	if projectType == "agent" {
		switch framework {
		case "adk":
			return "AI агент, созданный с использованием ADK (Agent Development Kit) - современного фреймворка для разработки AI агентов"
		case "langgraph":
			return "AI агент, созданный с использованием LangGraph - фреймворка для создания stateful, multi-actor applications с LLM"
		case "crewai":
			return "AI агент, созданный с использованием CrewAI - фреймворка для создания команд AI агентов, которые могут работать вместе"
		default:
			return "AI агент для автоматизации задач и интеллектуальной обработки данных"
		}
	} else if projectType == "mcp" {
		return "MCP (Model Context Protocol) сервер для интеграции внешних инструментов и ресурсов с AI агентами"
	}
	return "Проект, созданный с помощью ai-agents-cli"
}

// getEnvOrDefault returns environment variable value or default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getGitAuthor returns the author from git config
func getGitAuthor() string {
	// Try to get user.name from git config
	cmd := exec.Command("git", "config", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	name := strings.TrimSpace(string(output))
	if name == "" {
		return ""
	}

	// Try to get user.email from git config
	cmd = exec.Command("git", "config", "user.email")
	output, err = cmd.Output()
	if err != nil {
		return name // Return just name if no email
	}

	email := strings.TrimSpace(string(output))
	if email == "" {
		return name // Return just name if no email
	}

	return fmt.Sprintf("%s <%s>", name, email)
}

// applyProjectOptions applies additional options to the created project
func (s *Scaffolder) applyProjectOptions(targetPath string, options []string) error {
	for _, option := range options {
		switch option {
		case "git_init":
			if err := s.initializeGitRepo(targetPath); err != nil {
				return fmt.Errorf("failed to initialize git repository: %w", err)
			}
		case "create_env":
			if err := s.createEnvFile(targetPath); err != nil {
				return fmt.Errorf("failed to create .env file: %w", err)
			}
		case "install_deps":
			if err := s.installDependencies(targetPath); err != nil {
				return fmt.Errorf("failed to install dependencies: %w", err)
			}
		}
	}
	return nil
}

// initializeGitRepo initializes a git repository in the project directory
func (s *Scaffolder) initializeGitRepo(targetPath string) error {
	// Change to project directory
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(targetPath); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	// Initialize git repository
	cmd := exec.Command("git", "init")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run git init: %w", err)
	}

	// Add all files
	cmd = exec.Command("git", "add", ".")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run git add: %w", err)
	}

	// Make initial commit
	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run git commit: %w", err)
	}

	return nil
}

// createEnvFile creates .env file from .env.example with user credentials
func (s *Scaffolder) createEnvFile(targetPath string) error {
	envExamplePath := filepath.Join(targetPath, "env.example")
	envPath := filepath.Join(targetPath, ".env")

	// Check if .env.example exists
	if _, err := os.Stat(envExamplePath); os.IsNotExist(err) {
		return nil // No .env.example file, skip
	}

	// Read .env.example content
	content, err := os.ReadFile(envExamplePath)
	if err != nil {
		return fmt.Errorf("failed to read env.example: %w", err)
	}

	// Get user credentials from environment or auth system
	envContent := string(content)

	// Replace placeholder values with actual values from environment
	envContent = strings.ReplaceAll(envContent, "your_iam_key_id_here", os.Getenv("IAM_KEY_ID"))
	envContent = strings.ReplaceAll(envContent, "your_iam_secret_key_here", os.Getenv("IAM_SECRET_KEY"))
	envContent = strings.ReplaceAll(envContent, "your_project_id_here", os.Getenv("PROJECT_ID"))
	envContent = strings.ReplaceAll(envContent, "your_customer_id_here", os.Getenv("CUSTOMER_ID"))

	// Set EVOLUTION_PROJECT_ID to PROJECT_ID value
	if projectID := os.Getenv("PROJECT_ID"); projectID != "" {
		envContent = strings.ReplaceAll(envContent, "EVOLUTION_PROJECT_ID=${PROJECT_ID}", "EVOLUTION_PROJECT_ID="+projectID)
	}

	// Write .env file
	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	return nil
}

// installDependencies installs Python dependencies
func (s *Scaffolder) installDependencies(targetPath string) error {
	// Change to project directory
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(targetPath); err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	// Check if pyproject.toml exists (preferred for uv)
	pyprojectPath := filepath.Join(targetPath, "pyproject.toml")
	requirementsPath := filepath.Join(targetPath, "requirements.txt")

	var cmd *exec.Cmd

	if _, err := os.Stat(pyprojectPath); err == nil {
		// Use uv sync if pyproject.toml exists
		cmd = exec.Command("uv", "sync")
	} else if _, err := os.Stat(requirementsPath); err == nil {
		// Fallback to pip if only requirements.txt exists
		cmd = exec.Command("pip", "install", "-r", "requirements.txt")
	} else {
		return nil // No dependency files, skip
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}
