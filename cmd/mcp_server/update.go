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
	updateName        string
	updateDescription string
	updateConfigFile  string
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update <server-id>",
	Short: "Обновить MCP сервер",
	Long:  "Обновляет существующий MCP сервер с новыми параметрами",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		serverID := args[0]
		errorHandler := errors.NewHandler()

		var req *api.MCPServerUpdateRequest

		if updateConfigFile != "" {
			data, err := os.ReadFile(updateConfigFile)
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

			req = &api.MCPServerUpdateRequest{
				Name:        config.Name,
				Description: config.Description,
				Options:     config.Options,
			}
		} else {
			req = &api.MCPServerUpdateRequest{}
			if updateName != "" {
				req.Name = updateName
			}
			if updateDescription != "" {
				req.Description = updateDescription
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

		server, err := apiClient.MCPServers.Update(ctx, serverID, req)
		if err != nil {
			appErr := errorHandler.WrapAPIError(err, "MCP_SERVER_UPDATE_FAILED", "Ошибка обновления MCP сервера")
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
		fmt.Println(successStyle.Render("✅ MCP сервер успешно обновлен"))
		fmt.Println()
		fmt.Printf("%s: %s\n", labelStyle.Render("ID"), valueStyle.Render(server.ID))
		fmt.Printf("%s: %s\n", labelStyle.Render("Название"), valueStyle.Render(server.Name))

		if server.Description != "" {
			fmt.Printf("%s: %s\n", labelStyle.Render("Описание"), valueStyle.Render(server.Description))
		}

		fmt.Printf("%s: %s\n", labelStyle.Render("Статус"), valueStyle.Render(server.Status))
		fmt.Printf("%s: %s\n", labelStyle.Render("Обновлен"), valueStyle.Render(server.UpdatedAt.Time.Format("02.01.2006 15:04:05")))
	},
}

func init() {
	RootCMD.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&updateName, "name", "n", "", "Новое название MCP сервера")
	updateCmd.Flags().StringVarP(&updateDescription, "description", "d", "", "Новое описание MCP сервера")
	updateCmd.Flags().StringVarP(&updateConfigFile, "config", "c", "", "Путь к файлу конфигурации (JSON)")
}
