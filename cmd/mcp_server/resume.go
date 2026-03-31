package mcp_server

import (
	"context"
	"fmt"

	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
	"github.com/spf13/cobra"
)

// resumeCmd represents the resume command
var resumeCmd = &cobra.Command{
	Use:   "resume <server-id>",
	Short: "Возобновить работу MCP сервера",
	Long:  "Возобновляет работу приостановленного MCP сервера",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		serverID := args[0]

		errorHandler := errors.NewHandler()
		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			appErr := errorHandler.WrapAPIError(err, "API_CLIENT_ERROR", "Ошибка получения API клиента")
			appErr = appErr.WithSuggestions(mcpAPISuggestions()...)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		if err := apiClient.MCPServers.Resume(ctx, serverID); err != nil {
			appErr := errorHandler.WrapAPIError(err, "MCP_SERVER_RESUME_FAILED", "Ошибка возобновления MCP сервера")
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

		// Выводим результат
		fmt.Println(successStyle.Render("✅ MCP сервер возобновлен"))
		fmt.Printf("ID: %s\n", serverID)
	},
}

func init() {
	RootCMD.AddCommand(resumeCmd)
}
