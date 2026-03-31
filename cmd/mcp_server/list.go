package mcp_server

import (
	"context"
	"fmt"

	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	limit  int
	offset int
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Список MCP серверов",
	Long: `Показывает список всех MCP серверов в проекте.

Команда отображает таблицу с информацией о всех MCP серверах:
• ID сервера
• Название сервера
• Текущий статус
• Дата создания
• Дата последнего обновления

Поддерживает постраничную навигацию с помощью флагов --limit и --offset.

Примеры использования:
  ai-agents-cli mcp-servers list
  ai-agents-cli mcp-servers list --limit 10
  ai-agents-cli mcp-servers list --offset 20 --limit 5`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Проверяем размер терминала
		if err := ui.CheckTerminalSize(); err != nil {
			errorHandler := errors.NewHandler()
			appErr := errorHandler.WrapUserError(err, "TERMINAL_SIZE_ERROR", "Ошибка размера терминала")
			fmt.Println(errorHandler.HandlePlain(appErr))
			return
		}

		// Показываем таблицу MCP серверов
		if err := ui.ShowMCPServersListFromAPI(ctx, limit, offset); err != nil {
			errorHandler := errors.NewHandler()
			appErr := errorHandler.WrapAPIError(err, "MCP_SERVERS_LIST_FAILED", "Ошибка получения списка MCP серверов")
			appErr = appErr.WithSuggestions(
				"Проверьте переменные окружения: IAM_KEY_ID, IAM_SECRET, IAM_ENDPOINT",
				"Убедитесь что вы авторизованы: ai-agents-cli auth login или в папке выполнения команды лежит .env файл с перемнными выше",
				"Проверьте доступность API: curl -I $IAM_ENDPOINT",
				"Обратитесь к администратору для получения учетных данных",
				"📚 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",
			)
			fmt.Println(errorHandler.HandlePlain(appErr))
			return
		}
	},
}

func init() {
	RootCMD.AddCommand(listCmd)

	listCmd.Flags().IntVarP(&limit, "limit", "l", 20, "Количество записей для отображения")
	listCmd.Flags().IntVarP(&offset, "offset", "o", 0, "Смещение для постраничной навигации")
}
