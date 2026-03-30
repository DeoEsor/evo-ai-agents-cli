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

var (
	force bool
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <server-id>",
	Short: "Удалить MCP сервер",
	Long:  "Удаляет существующий MCP сервер",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		serverID := args[0]

		if !force {
			// Спрашиваем подтверждение
			fmt.Printf("Вы уверены, что хотите удалить MCP сервер %s? (y/N): ", serverID)
			var response string
			fmt.Scanln(&response)

			if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
				fmt.Println("❌ Операция отменена")
				return
			}
		}

		errorHandler := errors.NewHandler()
		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			appErr := errorHandler.WrapAPIError(err, "API_CLIENT_ERROR", "Ошибка получения API клиента")
			appErr = appErr.WithSuggestions(mcpAPISuggestions()...)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		if err := apiClient.MCPServers.Delete(ctx, serverID); err != nil {
			appErr := errorHandler.WrapAPIError(err, "MCP_SERVER_DELETE_FAILED", "Ошибка удаления MCP сервера")
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
		fmt.Println(successStyle.Render("✅ MCP сервер успешно удален"))
		fmt.Printf("ID: %s\n", serverID)
	},
}

func init() {
	RootCMD.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVarP(&force, "force", "f", false, "Принудительное удаление без подтверждения")
}
