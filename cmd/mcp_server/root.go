package mcp_server

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var (
	isVerbose bool
)

// RootCMD represents the base command when called without any subcommands
var RootCMD = &cobra.Command{
	Use:   "mcp-servers",
	Short: "Управление MCP серверами",
	Long: `Управление MCP (Model Context Protocol) серверами в Cloud.ru.

Доступные операции:
• list          - Просмотр списка серверов
• get           - Получение информации о сервере
• create        - Создание нового сервера
• update        - Обновление существующего сервера
• delete        - Удаление сервера
• resume        - Возобновление работы сервера
• suspend       - Приостановка сервера
• history       - История операций сервера
• deploy        - Развертывание серверов из YAML
• tools         - Список инструментов сервера
• execute       - Вызов инструмента сервера
• completion    - Генерация скрипта автодополнения

Примеры:
  ai-agents-cli mcp-servers list
  ai-agents-cli mcp-servers get <server-id>
  ai-agents-cli mcp-servers tools <server-id>
  ai-agents-cli mcp-servers deploy mcp-servers.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Команда MCP серверов вызвана без подкоманды")
		// Показываем справку если нет подкоманд
		cmd.Help()
	},
	Args: cobra.ArbitraryArgs,
}

func init() {
	log.Debug("Инициализация MCP серверов команды")
}
