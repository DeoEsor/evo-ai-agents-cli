package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/validator"
	"github.com/spf13/cobra"
)

var (
	validateFile string
	validateDir  string
)

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate [config-file|directory]",
	Short: "Валидация конфигурационных файлов",
	Long: `Валидация конфигурационных файлов для AI Agents.

Эта команда проверяет корректность YAML и JSON файлов конфигурации
перед их развертыванием в облачной платформе.

Поддерживаемые форматы:
• YAML файлы (.yaml, .yml)
• JSON файлы (.json)

Поддерживаемые типы конфигураций:
• Агенты (agents)
• MCP серверы (mcp-servers) 
• Системы агентов (agent-systems)

Примеры использования:
  ai-agents-cli validate examples/agents.yaml
  ai-agents-cli validate examples/
  ai-agents-cli validate --file config.yaml`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Запуск валидации конфигурационных файлов")

		// Создаем валидатор
		configValidator := validator.NewConfigValidator()
		log.Debug("Валидатор создан")

		// Загружаем схемы
		log.Info("Загрузка схем валидации")
		if err := loadSchemas(configValidator); err != nil {
			log.Warn("Failed to load schemas", "error", err)
		} else {
			log.Info("Схемы валидации загружены успешно")
		}

		// Определяем файлы для валидации
		var files []string

		if len(args) > 0 {
			// Проверяем, является ли аргумент директорией
			if info, err := os.Stat(args[0]); err == nil && info.IsDir() {
				var err error
				files, err = findConfigFiles(args[0])
				if err != nil {
					log.Fatal("Failed to find config files", "error", err, "dir", args[0])
				}
			} else {
				files = []string{args[0]}
			}
		} else if validateFile != "" {
			files = []string{validateFile}
		} else if validateDir != "" {
			// Находим все конфигурационные файлы в директории
			var err error
			files, err = findConfigFiles(validateDir)
			if err != nil {
				log.Fatal("Failed to find config files", "error", err, "dir", validateDir)
			}
		} else {
			// Ищем файлы в текущей директории
			var err error
			files, err = findConfigFiles(".")
			if err != nil {
				log.Fatal("Failed to find config files", "error", err)
			}
		}

		if len(files) == 0 {
			log.Error("Конфигурационные файлы не найдены")
			log.Fatal("No configuration files found")
		}

		log.Info("Найдено файлов для валидации", "count", len(files), "files", files)

		// Валидируем каждый файл
		allValid := true
		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		fmt.Println(headerStyle.Render("🔍 Валидация конфигурационных файлов"))
		fmt.Println()

		for _, file := range files {
			fmt.Printf("Проверка %s...\n", file)
			log.Debug("Валидация файла", "file", file)

			result, err := configValidator.ValidateFile(file)
			if err != nil {
				log.Error("Ошибка при валидации файла", "file", file, "error", err)
				allValid = false
				continue
			}

			if !result.Valid {
				log.Warn("Файл не прошел валидацию", "file", file, "errors", len(result.Errors))
				allValid = false
				configValidator.PrintErrors(result)
			} else {
				log.Info("Файл валиден", "file", file)
			}
			fmt.Println()
		}

		// Выводим итоговый результат
		if allValid {
			fmt.Println(headerStyle.Copy().Foreground(lipgloss.Color("2")).Render("🎉 Все файлы валидны!"))
			os.Exit(0)
		} else {
			fmt.Println(headerStyle.Copy().Foreground(lipgloss.Color("1")).Render("❌ Обнаружены ошибки валидации"))
			os.Exit(1)
		}
	},
}

func loadSchemas(validator *validator.ConfigValidator) error {
	schemas := map[string]string{
		"mcp-servers":   "schemas/schema.json",
		"agents":        "schemas/schema.json",
		"agent-systems": "schemas/schema.json",
	}

	for name, path := range schemas {
		if err := validator.LoadSchema(name, path); err != nil {
			log.Warn("Failed to load schema", "schema", name, "path", path, "error", err)
		}
	}

	return nil
}

func findConfigFiles(dir string) ([]string, error) {
	var files []string

	extensions := []string{".yaml", ".yml", ".json"}

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		name := d.Name()
		for _, ext := range extensions {
			if len(name) > len(ext) && name[len(name)-len(ext):] == ext {
				files = append(files, path)
				break
			}
		}
		return nil
	})

	return files, err
}

func init() {
	RootCMD.AddCommand(validateCmd)

	validateCmd.Flags().StringVarP(&validateFile, "file", "f", "", "Файл для валидации")
	validateCmd.Flags().StringVarP(&validateDir, "dir", "d", "", "Директория с файлами для валидации")
}
