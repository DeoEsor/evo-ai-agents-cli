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
	Long: `Управление MCP (Model Context Protocol) серверами.

MCP серверы предоставляют контекст и инструменты для AI агентов.
Эта команда позволяет создавать, настраивать и управлять MCP серверами
в вашем проекте.

Доступные операции:
• list - Просмотр списка серверов
• get - Получение информации о сервере
• create - Создание нового сервера
• update - Обновление существующего сервера
• delete - Удаление сервера
• resume - Возобновление работы сервера
• suspend - Приостановка сервера

Примеры использования:
  ai-agents-cli mcp-servers list
  ai-agents-cli mcp-servers get server-id
  ai-agents-cli mcp-servers create --name my-server`,
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
