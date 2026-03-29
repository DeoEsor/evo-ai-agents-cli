package agent

import (
	"context"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

var (
	forceDelete bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete <agent-id>",
	Short: "Удалить агента",
	Long:  "Удаляет существующего AI агента",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		agentID := args[0]

		if !forceDelete {
			fmt.Printf("Вы уверены, что хотите удалить агента %s? (y/N): ", agentID)
			var response string
			fmt.Scanln(&response)

			if response != "y" && response != "Y" && response != "yes" && response != "Yes" {
				fmt.Println("❌ Операция отменена")
				return
			}
		}

		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			log.Fatal("Failed to get API client", "error", err)
		}

		if err := apiClient.Agents.Delete(ctx, agentID); err != nil {
			log.Fatal("Failed to delete agent", "error", err, "agent_id", agentID)
		}

		successStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		fmt.Println(successStyle.Render("✅ Агент успешно удален"))
		fmt.Printf("ID: %s\n", agentID)
	},
}

func init() {
	RootCMD.AddCommand(deleteCmd)

	deleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Принудительное удаление без подтверждения")
}
