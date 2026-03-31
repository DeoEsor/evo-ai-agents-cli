package agent

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	isVerbose bool
)

// RootCMD represents the base command when called without any subcommands
var RootCMD = &cobra.Command{
	Use:   "agents",
	Short: "Управление AI агентами",
	Long: `Управление AI агентами в облачной платформе.

AI агенты - это интеллектуальные помощники, которые могут выполнять
различные задачи с использованием языковых моделей и MCP серверов.

Доступные операции:
• list - Просмотр списка агентов
• get - Получение информации об агенте
• create - Создание нового агента
• update - Обновление существующего агента
• delete - Удаление агента
• resume - Возобновление работы агента
• suspend - Приостановка агента

Примеры использования:
  ai-agents-cli agents list
  ai-agents-cli agents get agent-id
  ai-agents-cli agents create --name my-agent`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Команда агентов вызвана без подкоманды")
		// Показываем справку если нет подкоманд
		cmd.Help()
	},
	Args: cobra.ArbitraryArgs,
}

func init() {
	log.Debug("Инициализация команды агентов")
}
