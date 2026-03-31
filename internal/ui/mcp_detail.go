package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
)

// MCPDetailModel представляет детальный просмотр MCP сервера с табами
type MCPDetailModel struct {
	MCPServer *api.MCPServer
	Tabs      *TabModel
}

// NewMCPDetailModel создает новую модель детального просмотра MCP сервера
func NewMCPDetailModel(mcpServer *api.MCPServer) *MCPDetailModel {
	tabs := []string{
		"📋 Общая информация",
		"🐳 Образ",
		"🔌 Порты",
		"📈 Масштабирование",
		"🌿 Окружение",
		"🔗 Интеграция",
		"⚙️ Опции",
		"📊 Инструменты",
	}

	content := []string{
		formatMCPGeneralInfo(mcpServer),
		formatMCPImageSource(mcpServer),
		formatMCPPorts(mcpServer),
		formatMCPScaling(mcpServer),
		formatMCPEnvironment(mcpServer),
		formatMCPIntegration(mcpServer),
		formatMCPOptions(mcpServer),
		formatMCPTools(mcpServer),
	}

	return &MCPDetailModel{
		MCPServer: mcpServer,
		Tabs:      NewTabModel(tabs, content),
	}
}

// Render отображает детальную информацию о MCP сервере
func (m *MCPDetailModel) Render() string {
	header := fmt.Sprintf("🔌 Детальная информация: %s", m.MCPServer.Name)
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		MarginBottom(1)

	return headerStyle.Render(header) + "\n" + m.Tabs.Render()
}

// formatMCPGeneralInfo форматирует общую информацию MCP сервера
func formatMCPGeneralInfo(mcp *api.MCPServer) string {
	// Стили для разных секций
	leftStyle := lipgloss.NewStyle().
		Width(50).
		Align(lipgloss.Left)

	rightStyle := lipgloss.NewStyle().
		Width(50).
		Align(lipgloss.Right)

	centerStyle := lipgloss.NewStyle().
		Width(100).
		Align(lipgloss.Center).
		MarginBottom(1)

	// Левая колонка - основная информация
	leftInfo := strings.Builder{}
	leftInfo.WriteString(fmt.Sprintf("🆔 ID: %s\n", mcp.ID))
	leftInfo.WriteString(fmt.Sprintf("📝 Название: %s\n", mcp.Name))
	leftInfo.WriteString(fmt.Sprintf("📄 Описание: %s\n", getDescription(mcp.Description)))
	leftInfo.WriteString(fmt.Sprintf("📊 Статус: %s\n", FormatMCPServerStatus(mcp.Status)))

	// Project ID не доступен в структуре MCPServer

	// Добавляем информацию о типе инстанса
	if mcp.InstanceType.ID != "" {
		leftInfo.WriteString(fmt.Sprintf("💻 Инстанс: %s (%s)\n", mcp.InstanceType.Name, mcp.InstanceType.SKUCode))
	}

	// Правая колонка - временные метки и пользователи
	rightInfo := strings.Builder{}
	rightInfo.WriteString(fmt.Sprintf("📅 Создан: %s\n", mcp.CreatedAt.Time.Format("02.01.2006 15:04:05")))
	rightInfo.WriteString(fmt.Sprintf("🔄 Обновлен: %s\n", mcp.UpdatedAt.Time.Format("02.01.2006 15:04:05")))

	if mcp.CreatedBy != "" {
		rightInfo.WriteString(fmt.Sprintf("👤 Создал: %s\n", mcp.CreatedBy))
	}
	if mcp.UpdatedBy != "" {
		rightInfo.WriteString(fmt.Sprintf("✏️ Изменил: %s\n", mcp.UpdatedBy))
	}

	// Центральная секция - URL и дополнительная информация
	centerInfo := strings.Builder{}
	if mcp.PublicURL != "" {
		centerInfo.WriteString(fmt.Sprintf("🌐 Публичный URL: %s\n", mcp.PublicURL))
	}

	// Объединяем все секции
	var result strings.Builder
	result.WriteString(centerStyle.Render(centerInfo.String()))

	// Создаем две колонки
	leftContent := leftStyle.Render(leftInfo.String())
	rightContent := rightStyle.Render(rightInfo.String())

	result.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, leftContent, rightContent))

	return result.String()
}

// formatMCPImageSource форматирует информацию об образе
func formatMCPImageSource(mcp *api.MCPServer) string {
	var info strings.Builder

	info.WriteString("🐳 Источник образа:\n\n")

	if mcp.ImageSource != nil {
		if registry, ok := mcp.ImageSource["registry"].(string); ok {
			info.WriteString(fmt.Sprintf("📦 Registry: %s\n", registry))
		}
		if repository, ok := mcp.ImageSource["repository"].(string); ok {
			info.WriteString(fmt.Sprintf("📁 Repository: %s\n", repository))
		}
		if tag, ok := mcp.ImageSource["tag"].(string); ok {
			info.WriteString(fmt.Sprintf("🏷️ Tag: %s\n", tag))
		}
		if digest, ok := mcp.ImageSource["digest"].(string); ok {
			info.WriteString(fmt.Sprintf("🔐 Digest: %s\n", digest))
		}
	}

	return info.String()
}

// formatMCPPorts форматирует информацию о портах
func formatMCPPorts(mcp *api.MCPServer) string {
	var info strings.Builder

	info.WriteString("🔌 Порты:\n\n")

	if len(mcp.ExposedPorts) > 0 {
		for i, port := range mcp.ExposedPorts {
			info.WriteString(fmt.Sprintf("  %d. Порт %d\n", i+1, port))
		}
	} else {
		info.WriteString("Нет открытых портов\n")
	}

	return info.String()
}

// formatMCPScaling форматирует настройки масштабирования
func formatMCPScaling(mcp *api.MCPServer) string {
	var info strings.Builder

	info.WriteString("📈 Настройки масштабирования:\n\n")

	if mcp.Scaling != nil && len(mcp.Scaling) > 0 {
		for key, value := range mcp.Scaling {
			info.WriteString(fmt.Sprintf("  🔧 %s: %v\n", key, value))
		}
	} else {
		info.WriteString("Настройки масштабирования не заданы\n")
	}

	return info.String()
}

// formatMCPEnvironment форматирует переменные окружения
func formatMCPEnvironment(mcp *api.MCPServer) string {
	var info strings.Builder

	info.WriteString("🌿 Переменные окружения:\n\n")

	if mcp.EnvironmentOptions != nil && len(mcp.EnvironmentOptions) > 0 {
		for key, value := range mcp.EnvironmentOptions {
			info.WriteString(fmt.Sprintf("  🔧 %s: %v\n", key, value))
		}
	} else {
		info.WriteString("Переменные окружения не заданы\n")
	}

	return info.String()
}

// formatMCPIntegration форматирует настройки интеграции
func formatMCPIntegration(mcp *api.MCPServer) string {
	var info strings.Builder

	info.WriteString("🔗 Настройки интеграции:\n\n")

	if mcp.IntegrationOptions != nil && len(mcp.IntegrationOptions) > 0 {
		for key, value := range mcp.IntegrationOptions {
			info.WriteString(fmt.Sprintf("  🔧 %s: %v\n", key, value))
		}
	} else {
		info.WriteString("Настройки интеграции не заданы\n")
	}

	return info.String()
}

// formatMCPOptions форматирует все опции
func formatMCPOptions(mcp *api.MCPServer) string {
	var info strings.Builder

	info.WriteString("⚙️ Все опции:\n\n")

	if mcp.Options != nil {
		for key, value := range mcp.Options {
			info.WriteString(fmt.Sprintf("🔧 %s: %v\n", key, value))
		}
	} else {
		info.WriteString("Нет дополнительных опций\n")
	}

	return info.String()
}

// formatMCPTools форматирует инструменты MCP сервера
func formatMCPTools(mcp *api.MCPServer) string {
	var info strings.Builder

	info.WriteString("📊 Инструменты MCP сервера:\n\n")

	if len(mcp.Tools) > 0 {
		for i, tool := range mcp.Tools {
			info.WriteString(fmt.Sprintf("🔧 %d. %s\n", i+1, tool.Name))
			if tool.Description != "" {
				info.WriteString(fmt.Sprintf("   📝 %s\n", tool.Description))
			}
			if tool.InputSchema != nil {
				info.WriteString("   📋 Схема входных данных:\n")
				for key, value := range tool.InputSchema {
					info.WriteString(fmt.Sprintf("     • %s: %v\n", key, value))
				}
			}
			info.WriteString("\n")
		}
	} else {
		info.WriteString("Нет доступных инструментов\n")
	}

	return info.String()
}

// FormatMCPServerStatus форматирует статус MCP сервера
func FormatMCPServerStatus(status string) string {
	switch status {
	case "ACTIVE":
		return "🟢 Активен"
	case "SUSPENDED":
		return "🟡 Приостановлен"
	case "ERROR":
		return "🔴 Ошибка"
	case "PENDING":
		return "⏳ Ожидает"
	case "CREATING":
		return "🔨 Создается"
	case "DELETING":
		return "🗑️ Удаляется"
	default:
		return "⚪ " + status
	}
}
