package agent

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history <agent-id>",
	Short: "История операций агента",
	Long:  "Показывает историю операций для указанного AI агента",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		agentID := args[0]

		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			log.Fatal("Failed to get API client", "error", err)
		}

		history, err := apiClient.Agents.GetHistory(ctx, agentID)
		if err != nil {
			log.Fatal("Failed to get agent history", "error", err, "agent_id", agentID)
		}

		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		statusStyle := lipgloss.NewStyle().
			Bold(true)

		shortID := agentID
		if len(shortID) > 8 {
			shortID = shortID[:8] + "..."
		}

		fmt.Println(headerStyle.Render(fmt.Sprintf("📜 История агента %s", shortID)))
		fmt.Println()

		if len(history.Data) == 0 {
			fmt.Println("🔍 История операций не найдена")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Время\tДействие\tСтатус\tСообщение")
		fmt.Fprintln(w, "-----\t--------\t------\t--------")

		for _, entry := range history.Data {
			status := entry.Status
			switch status {
			case "SUCCESS":
				status = statusStyle.Copy().Foreground(lipgloss.Color("2")).Render("✅ Успех")
			case "ERROR":
				status = statusStyle.Copy().Foreground(lipgloss.Color("1")).Render("❌ Ошибка")
			case "PENDING":
				status = statusStyle.Copy().Foreground(lipgloss.Color("3")).Render("⏳ В процессе")
			default:
				status = statusStyle.Copy().Foreground(lipgloss.Color("8")).Render("⚪ " + status)
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				entry.CreatedAt.Format("02.01.2006 15:04:05"),
				entry.Action,
				status,
				entry.Message,
			)
		}

		w.Flush()
	},
}

func init() {
	RootCMD.AddCommand(historyCmd)
}
