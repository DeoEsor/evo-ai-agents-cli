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
	updateName        string
	updateDescription string
	updateConfigFile  string
)

var updateCmd = &cobra.Command{
	Use:   "update <agent-id>",
	Short: "Обновить агента",
	Long:  "Обновляет существующего AI агента с новыми параметрами",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		agentID := args[0]
		errorHandler := errors.NewHandler()

		var req *api.AgentUpdateRequest

		if updateConfigFile != "" {
			data, err := os.ReadFile(updateConfigFile)
			if err != nil {
				appErr := errorHandler.WrapFileSystemError(err, "CONFIG_READ_ERROR", "Ошибка чтения файла конфигурации")
				fmt.Println(errorHandler.HandlePlain(appErr))
				os.Exit(1)
			}

			var config struct {
				Name               string                 `json:"name"`
				Description        string                 `json:"description"`
				Options            map[string]interface{} `json:"options"`
				IntegrationOptions map[string]interface{} `json:"integration_options"`
			}

			if err := json.Unmarshal(data, &config); err != nil {
				appErr := errorHandler.WrapValidationError(err, "CONFIG_PARSE_ERROR", "Ошибка парсинга файла конфигурации")
				fmt.Println(errorHandler.HandlePlain(appErr))
				os.Exit(1)
			}

			req = &api.AgentUpdateRequest{
				Name:               config.Name,
				Description:        config.Description,
				Options:            config.Options,
				IntegrationOptions: config.IntegrationOptions,
			}
		} else {
			req = &api.AgentUpdateRequest{}
			if updateName != "" {
				req.Name = updateName
			}
			if updateDescription != "" {
				req.Description = updateDescription
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

		agent, err := apiClient.Agents.Update(ctx, agentID, req)
		if err != nil {
			appErr := errorHandler.WrapAPIError(err, "AGENT_UPDATE_FAILED", "Ошибка обновления агента")
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

		fmt.Println(successStyle.Render("✅ Агент успешно обновлен"))
		fmt.Println()
		fmt.Printf("%s: %s\n", labelStyle.Render("ID"), valueStyle.Render(agent.ID))
		fmt.Printf("%s: %s\n", labelStyle.Render("Название"), valueStyle.Render(agent.Name))

		if agent.Description != "" {
			fmt.Printf("%s: %s\n", labelStyle.Render("Описание"), valueStyle.Render(agent.Description))
		}

		fmt.Printf("%s: %s\n", labelStyle.Render("Статус"), valueStyle.Render(agent.Status))
		fmt.Printf("%s: %s\n", labelStyle.Render("Обновлен"), valueStyle.Render(agent.UpdatedAt.Time.Format("02.01.2006 15:04:05")))
	},
}

func init() {
	RootCMD.AddCommand(updateCmd)

	updateCmd.Flags().StringVarP(&updateName, "name", "n", "", "Новое название агента")
	updateCmd.Flags().StringVarP(&updateDescription, "description", "d", "", "Новое описание агента")
	updateCmd.Flags().StringVarP(&updateConfigFile, "config", "c", "", "Путь к файлу конфигурации (JSON)")
}
