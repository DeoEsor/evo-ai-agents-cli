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
	Long: `Управление AI агентами в облачной платформе Cloud.ru.

Доступные операции:
• list          - Просмотр списка агентов
• get           - Получение информации об агенте
• create        - Создание нового агента
• update        - Обновление существующего агента
• delete        - Удаление агента
• resume        - Возобновление работы агента
• suspend       - Приостановка агента
• history       - История операций агента
• deploy        - Развертывание агентов из YAML
• marketplace   - Поиск агентов в маркетплейсе
• completion    - Генерация скрипта автодополнения

Примеры:
  ai-agents-cli agents list
  ai-agents-cli agents get <agent-id>
  ai-agents-cli agents create --name my-agent
  ai-agents-cli agents deploy agents.yaml`,
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
