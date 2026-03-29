package agent

import (
	"context"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

var suspendCmd = &cobra.Command{
	Use:   "suspend <agent-id>",
	Short: "Приостановить работу агента",
	Long:  "Приостанавливает работу активного AI агента",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		agentID := args[0]

		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			log.Fatal("Failed to get API client", "error", err)
		}

		if err := apiClient.Agents.Suspend(ctx, agentID); err != nil {
			log.Fatal("Failed to suspend agent", "error", err, "agent_id", agentID)
		}

		warningStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("3")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		fmt.Println(warningStyle.Render("⏸️  Агент приостановлен"))
		fmt.Printf("ID: %s\n", agentID)
	},
}

func init() {
	RootCMD.AddCommand(suspendCmd)
}
