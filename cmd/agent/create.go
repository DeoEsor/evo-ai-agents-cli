package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
	"github.com/spf13/cobra"
)

var (
	createName        string
	createDescription string
	createConfigFile  string
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Создать нового агента",
	Long:  "Создает нового AI агента с указанными параметрами",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		errorHandler := errors.NewHandler()

		var req *api.AgentCreateRequest

		if createConfigFile != "" {
			data, err := os.ReadFile(createConfigFile)
			if err != nil {
				appErr := errorHandler.WrapFileSystemError(err, "CONFIG_READ_ERROR", "Ошибка чтения файла конфигурации")
				fmt.Println(errorHandler.HandlePlain(appErr))
				os.Exit(1)
			}

			var config struct {
				Name               string                 `json:"name"`
				Description        string                 `json:"description"`
				InstanceTypeID     string                 `json:"instance_type_id"`
				Options            map[string]interface{} `json:"options"`
				MCPServers         []string               `json:"mcp_servers"`
				IntegrationOptions map[string]interface{} `json:"integration_options"`
			}

			if err := json.Unmarshal(data, &config); err != nil {
				appErr := errorHandler.WrapValidationError(err, "CONFIG_PARSE_ERROR", "Ошибка парсинга файла конфигурации")
				fmt.Println(errorHandler.HandlePlain(appErr))
				os.Exit(1)
			}

			req = &api.AgentCreateRequest{
				Name:               config.Name,
				Description:        config.Description,
				InstanceTypeID:     config.InstanceTypeID,
				Options:            config.Options,
				MCPServers:         config.MCPServers,
				IntegrationOptions: config.IntegrationOptions,
			}
		} else {
			if createName == "" {
				appErr := errors.ValidationError("MISSING_NAME", "Название обязательно").
					WithSuggestions("Используйте --name флаг или --config файл")
				fmt.Println(errorHandler.HandlePlain(appErr))
				os.Exit(1)
			}

			req = &api.AgentCreateRequest{
				Name:        createName,
				Description: createDescription,
			}
		}

		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			appErr := errorHandler.WrapAPIError(err, "API_CLIENT_ERROR", "Ошибка получения API клиента")
			appErr = appErr.WithSuggestions(agentAPISuggestions()...)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		agent, err := apiClient.Agents.Create(ctx, req)
		if err != nil {
			appErr := errorHandler.WrapAPIError(err, "AGENT_CREATE_FAILED", "Ошибка создания агента")
			appErr = appErr.WithSuggestions(agentAPISuggestions()...)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		successStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		labelStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99"))

		valueStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

		fmt.Println(successStyle.Render("✅ Агент успешно создан"))
		fmt.Println()
		fmt.Printf("%s: %s\n", labelStyle.Render("ID"), valueStyle.Render(agent.ID))
		fmt.Printf("%s: %s\n", labelStyle.Render("Название"), valueStyle.Render(agent.Name))

		if agent.Description != "" {
			fmt.Printf("%s: %s\n", labelStyle.Render("Описание"), valueStyle.Render(agent.Description))
		}

		fmt.Printf("%s: %s\n", labelStyle.Render("Статус"), valueStyle.Render(agent.Status))
		fmt.Printf("%s: %s\n", labelStyle.Render("Создан"), valueStyle.Render(agent.CreatedAt.Time.Format("02.01.2006 15:04:05")))
	},
}

func init() {
	RootCMD.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&createName, "name", "n", "", "Название агента")
	createCmd.Flags().StringVarP(&createDescription, "description", "d", "", "Описание агента")
	createCmd.Flags().StringVarP(&createConfigFile, "config", "c", "", "Путь к файлу конфигурации (JSON)")
}
