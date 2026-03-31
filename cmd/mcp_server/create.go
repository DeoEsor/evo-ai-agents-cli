package mcp_server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
	"github.com/spf13/cobra"
)

var (
	name        string
	description string
	configFile  string
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать новый MCP сервер",
	Long:  "Создает новый MCP сервер с указанными параметрами",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		errorHandler := errors.NewHandler()

		var req *api.MCPServerCreateRequest

		if configFile != "" {
			data, err := os.ReadFile(configFile)
			if err != nil {
				appErr := errorHandler.WrapFileSystemError(err, "CONFIG_READ_ERROR", "Ошибка чтения файла конфигурации")
				fmt.Println(errorHandler.HandlePlain(appErr))
				os.Exit(1)
			}

			var config struct {
				Name        string                 `json:"name"`
				Description string                 `json:"description"`
				Options     map[string]interface{} `json:"options"`
			}

			if err := json.Unmarshal(data, &config); err != nil {
				appErr := errorHandler.WrapValidationError(err, "CONFIG_PARSE_ERROR", "Ошибка парсинга файла конфигурации")
				fmt.Println(errorHandler.HandlePlain(appErr))
				os.Exit(1)
			}

			req = &api.MCPServerCreateRequest{
				Name:        config.Name,
				Description: config.Description,
				Options:     config.Options,
			}
		} else {
			if name == "" {
				appErr := errors.ValidationError("MISSING_NAME", "Название обязательно").
					WithSuggestions("Используйте --name флаг или --config файл")
				fmt.Println(errorHandler.HandlePlain(appErr))
				os.Exit(1)
			}

			req = &api.MCPServerCreateRequest{
				Name:        name,
				Description: description,
				Options:     make(map[string]interface{}),
			}
		}

		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			appErr := errorHandler.WrapAPIError(err, "API_CLIENT_ERROR", "Ошибка получения API клиента")
			appErr = appErr.WithSuggestions(mcpAPISuggestions()...)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		server, err := apiClient.MCPServers.Create(ctx, req)
		if err != nil {
			appErr := errorHandler.WrapAPIError(err, "MCP_SERVER_CREATE_FAILED", "Ошибка создания MCP сервера")
			appErr = appErr.WithSuggestions(mcpAPISuggestions()...)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		// Создаем стили для вывода
		successStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		labelStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99"))

		valueStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

		// Выводим результат
		fmt.Println(successStyle.Render("✅ MCP сервер успешно создан"))
		fmt.Println()
		fmt.Printf("%s: %s\n", labelStyle.Render("ID"), valueStyle.Render(server.ID))
		fmt.Printf("%s: %s\n", labelStyle.Render("Название"), valueStyle.Render(server.Name))

		if server.Description != "" {
			fmt.Printf("%s: %s\n", labelStyle.Render("Описание"), valueStyle.Render(server.Description))
		}

		fmt.Printf("%s: %s\n", labelStyle.Render("Статус"), valueStyle.Render(server.Status))
		fmt.Printf("%s: %s\n", labelStyle.Render("Создан"), valueStyle.Render(server.CreatedAt.Time.Format("02.01.2006 15:04:05")))
	},
}

func init() {
	RootCMD.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&name, "name", "n", "", "Название MCP сервера")
	createCmd.Flags().StringVarP(&description, "description", "d", "", "Описание MCP сервера")
	createCmd.Flags().StringVarP(&configFile, "config", "c", "", "Путь к файлу конфигурации (JSON)")
}
