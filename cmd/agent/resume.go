package agent

import (
	"context"
	"fmt"

	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume <agent-id>",
	Short: "Возобновить работу агента",
	Long:  "Возобновляет работу приостановленного AI агента",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		agentID := args[0]

		errorHandler := errors.NewHandler()
		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			appErr := errorHandler.WrapAPIError(err, "API_CLIENT_ERROR", "Ошибка получения API клиента")
			appErr = appErr.WithSuggestions(agentAPISuggestions()...)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		if err := apiClient.Agents.Resume(ctx, agentID); err != nil {
			appErr := errorHandler.WrapAPIError(err, "AGENT_RESUME_FAILED", "Ошибка возобновления агента")
			appErr = appErr.WithSuggestions(agentAPISuggestions()...)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		successStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		fmt.Println(successStyle.Render("✅ Агент возобновлен"))
		fmt.Printf("ID: %s\n", agentID)
	},
}

func init() {
	RootCMD.AddCommand(resumeCmd)
}
