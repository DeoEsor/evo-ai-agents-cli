package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	authCmd "github.com/cloud-ru/evo-ai-agents-cli/cmd/auth"
	"github.com/cloud-ru/evo-ai-agents-cli/cmd/create"
	"github.com/spf13/cobra"
)

var (
	isVerbose bool
)

// RootCMD represents the base command when called without any subcommands
var RootCMD = &cobra.Command{
	Use:   "ai-agents-cli",
	Short: "CLI инструмент для управления AI Agents в облачной платформе Cloud.ru",
	Long: `AI Agents CLI - это мощный инструмент командной строки для управления 
и развертывания AI агентов в облачной платформе Cloud.ru.

Основные возможности:
• Валидация конфигурационных файлов
• Управление MCP серверами
• Создание и настройка агентов
• Управление системами агентов
• Создание проектов из шаблонов
• Интеграция с CI/CD процессами

Для начала работы используйте команду 'validate' для проверки конфигурации
или '--help' для просмотра всех доступных команд.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Получаем значение флага verbose
		verbose, _ := cmd.Flags().GetBool("verbose")

		// Настройка логирования
		logger := log.New(os.Stderr)
		logger.SetReportTimestamp(true)
		logger.SetReportCaller(true)

		// Установка уровня логирования
		if verbose {
			logger.SetLevel(log.DebugLevel)
			logger.Info("Включен подробный режим логирования")
		} else {
			logger.SetLevel(log.InfoLevel)
		}

		log.SetDefault(logger)
		log.Debug("AI Agents CLI запущен", "version", "1.0.0", "verbose", verbose)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Показываем красивый help если нет аргументов
		if len(args) == 0 {
			log.Debug("Показ справки по командам")
			showBeautifulHelp()
		}
	},
	Args: cobra.ArbitraryArgs,
}

func init() {
	RootCMD.PersistentFlags().
		BoolVarP(&isVerbose, "verbose", "v", false, "Детализация процесса")

	// Custom help only for the root command; subcommands use Cobra default
	defaultHelp := func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	}
	RootCMD.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if cmd == RootCMD {
			showBeautifulHelp()
		} else {
			defaultHelp(cmd, args)
		}
	})

	// Add commands
	RootCMD.AddCommand(authCmd.RootCMD)
	RootCMD.AddCommand(create.RootCMD)
}

// showBeautifulHelp displays a beautifully formatted help message
func showBeautifulHelp() {
	// Define styles
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6B6B")).
		Margin(1, 0)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4ECDC4")).
		Bold(true).
		Margin(1, 0)

	descriptionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#95A5A6")).
		Margin(0, 0, 1, 0)

	commandStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#3498DB")).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7F8C8D")).
		Margin(0, 0, 0, 2)

	flagStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E74C3C")).
		Bold(true)

	// Header
	header := titleStyle.Render("🤖 AI Agents CLI")
	subtitle := subtitleStyle.Render("Мощный инструмент для управления AI агентами в Cloud.ru")
	description := descriptionStyle.Render("Создавайте, настраивайте и развертывайте AI агентов с помощью простых команд")

	// Commands section
	commandsTitle := subtitleStyle.Render("📋 Доступные команды:")

	commands := []struct {
		name        string
		description string
	}{
		{"create", "Создание проектов из шаблонов (agent, mcp)"},
		{"auth", "Управление аутентификацией (login, logout, status)"},
		{"agents", "Управление AI агентами"},
		{"mcp-servers", "Управление MCP серверами"},
		{"system", "Управление системами агентов"},
		{"ci", "CI/CD функции"},
		{"validate", "Валидация конфигурационных файлов"},
		{"completion", "Генерация скриптов автодополнения"},
	}

	var commandsText string
	for _, cmd := range commands {
		commandsText += commandStyle.Render("  "+cmd.name) + "\n" +
			descStyle.Render("    "+cmd.description) + "\n"
	}

	// Flags section
	flagsTitle := subtitleStyle.Render("🚩 Флаги:")
	flagsText := flagStyle.Render("  -v, --verbose") + "\n" +
		descStyle.Render("    Детализация процесса") + "\n" +
		flagStyle.Render("  -h, --help") + "\n" +
		descStyle.Render("    Показать справку")

	// Examples section
	examplesTitle := subtitleStyle.Render("💡 Примеры использования:")
	examplesText := commandStyle.Render("  ai-agents-cli create agent my-agent") + "\n" +
		descStyle.Render("    Создать новый AI агент") + "\n" +
		commandStyle.Render("  ai-agents-cli auth login") + "\n" +
		descStyle.Render("    Войти в систему") + "\n" +
		commandStyle.Render("  ai-agents-cli agents list") + "\n" +
		descStyle.Render("    Показать список агентов")

	// Documentation
	docsTitle := subtitleStyle.Render("📚 Документация:")
	docsText := descStyle.Render("  📖 Подробная документация: https://cloud.ru/docs/ai-agents/ug/index?source-platform=Evolution")

	// Combine all parts
	helpText := fmt.Sprintf("%s\n%s\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s\n%s\n\n%s\n%s",
		header,
		subtitle,
		description,
		commandsTitle,
		commandsText,
		flagsTitle,
		flagsText,
		examplesTitle,
		examplesText,
		docsTitle,
		docsText,
	)

	fmt.Println(helpText)
}
