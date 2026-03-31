package mcp_server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/ui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	outputFormat string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <server-id>",
	Short: "Получить информацию о MCP сервере",
	Long:  "Показывает подробную информацию о конкретном MCP сервере",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		serverID := args[0]

		// Создаем обработчик ошибок
		errorHandler := errors.NewHandler()

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			appErr := errorHandler.WrapAPIError(err, "API_CLIENT_ERROR", "Ошибка получения API клиента")
			appErr = appErr.WithSuggestions(
				"Проверьте переменные окружения: IAM_KEY_ID, IAM_SECRET, IAM_ENDPOINT",
				"Убедитесь что вы авторизованы: ai-agents-cli auth login или в папке выполнения команды лежит .env файл с перемнными выше",
				"Проверьте доступность API: curl -I $IAM_ENDPOINT",
				"Обратитесь к администратору для получения учетных данных",
				"📚 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",
			)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		// Получаем информацию о MCP сервере
		server, err := apiClient.MCPServers.Get(ctx, serverID)
		if err != nil {
			appErr := errorHandler.WrapAPIError(err, "MCP_SERVER_GET_FAILED", "Ошибка получения MCP сервера")
			appErr = appErr.WithSuggestions(
				"Проверьте правильность ID сервера: "+serverID,
				"Убедитесь что сервер существует: ai-agents-cli mcp-servers list",
				"Проверьте переменные окружения: IAM_KEY_ID, IAM_SECRET, IAM_ENDPOINT",
				"Убедитесь что вы авторизованы: ai-agents-cli auth login или в папке выполнения команды лежит .env файл с перемнными выше",
				"📚 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",
			)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		if outputFormat == "json" {
			// Выводим в JSON формате
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(server); err != nil {
				appErr := errorHandler.WrapFileSystemError(err, "JSON_ENCODE_ERROR", "Ошибка кодирования JSON")
				appErr = appErr.WithSuggestions(
					"Проверьте доступность stdout",
					"Попробуйте перенаправить вывод в файл",
					"Проверьте размер данных для вывода",
					"📚 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",
				)
				fmt.Println(errorHandler.HandlePlain(appErr))
				os.Exit(1)
			}
			return
		}

		// Показываем детальную информацию с табами
		if isTerminal() {
			// Интерактивная версия с табами
			program := ui.NewMCPDetailViewModel(ui.NewMCPDetailModel(server))
			if err := program.Start(); err != nil {
				log.Fatal("Failed to start detail view", "error", err)
			}
		} else {
			// Простая версия для не-терминала
			fmt.Printf("🔧 MCP Сервер: %s\n", server.Name)
			fmt.Printf("🆔 ID: %s\n", server.ID)
			fmt.Printf("📊 Статус: %s\n", server.Status)
		}
	},
}

// isTerminal проверяет, является ли терминал терминалом
func isTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func init() {
	RootCMD.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&outputFormat, "output", "o", "table", "Формат вывода (table, json)")
}
