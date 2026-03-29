package agent

import (
	"context"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
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

		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			log.Fatal("Failed to get API client", "error", err)
		}

		if err := apiClient.Agents.Resume(ctx, agentID); err != nil {
			log.Fatal("Failed to resume agent", "error", err, "agent_id", agentID)
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
