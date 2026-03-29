package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/api"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
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

		var req *api.AgentUpdateRequest

		if updateConfigFile != "" {
			data, err := os.ReadFile(updateConfigFile)
			if err != nil {
				log.Fatal("Failed to read config file", "error", err, "file", updateConfigFile)
			}

			var config struct {
				Name               string                 `json:"name"`
				Description        string                 `json:"description"`
				Options            map[string]interface{} `json:"options"`
				IntegrationOptions map[string]interface{} `json:"integration_options"`
			}

			if err := json.Unmarshal(data, &config); err != nil {
				log.Fatal("Failed to parse config file", "error", err)
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
			log.Fatal("Failed to get API client", "error", err)
		}

		agent, err := apiClient.Agents.Update(ctx, agentID, req)
		if err != nil {
			log.Fatal("Failed to update agent", "error", err, "agent_id", agentID)
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
