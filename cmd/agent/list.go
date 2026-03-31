package agent

import (
	"context"
	"fmt"

	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/ui"
	"github.com/spf13/cobra"
)

var (
	listLimit  int
	listOffset int
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Список AI агентов",
	Long: `Показывает список всех AI агентов в проекте.

Команда отображает таблицу с информацией о всех агентах:
• ID агента
• Название агента
• Описание агента
• Текущий статус
• Тип агента
• Дата создания
• Дата последнего обновления

Поддерживает постраничную навигацию с помощью флагов --limit и --offset.

Примеры использования:
  ai-agents-cli agents list
  ai-agents-cli agents list --limit 10
  ai-agents-cli agents list --offset 20 --limit 5`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		// Показываем таблицу агентов
		if err := ui.ShowAgentsListFromAPI(ctx, listLimit, listOffset); err != nil {
			// Создаем обработчик ошибок
			errorHandler := errors.NewHandler()
			appErr := errorHandler.WrapAPIError(err, "AGENTS_LIST_FAILED", "Ошибка получения списка агентов")

			// Добавляем подсказки для решения проблемы
			appErr = appErr.WithSuggestions(
				"Проверьте переменные окружения: IAM_KEY_ID, IAM_SECRET, IAM_ENDPOINT",
				"Убедитесь что вы авторизованы: ai-agents-cli auth login или в папке выполнения команды лежит .env файл с перемнными выше",
				"Проверьте доступность API: curl -I $IAM_ENDPOINT",
				"Обратитесь к администратору для получения учетных данных",
				"📚 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution",
			)

			// Отображаем простую ошибку без стилей и дублирования
			fmt.Println(errorHandler.HandlePlain(appErr))
			return
		}
	},
}

func init() {
	listCmd.Flags().IntVarP(&listLimit, "limit", "l", 20, "Количество записей для отображения")
	listCmd.Flags().IntVarP(&listOffset, "offset", "o", 0, "Смещение для постраничной навигации")
}
