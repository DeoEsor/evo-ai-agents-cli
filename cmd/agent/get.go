package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/ui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	agentOutputFormat string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <agent-id>",
	Short: "Получить информацию об агенте",
	Long:  "Показывает подробную информацию о конкретном агенте",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		agentID := args[0]

		// Получаем API клиент из DI контейнера
		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			log.Fatal("Failed to get API client", "error", err)
		}

		// Получаем информацию об агенте
		agent, err := apiClient.Agents.Get(ctx, agentID)
		if err != nil {
			log.Fatal("Failed to get agent", "error", err, "agent_id", agentID)
		}

		if agentOutputFormat == "json" {
			// Выводим в JSON формате
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(agent); err != nil {
				log.Fatal("Failed to encode JSON", "error", err)
			}
			return
		}

		// Показываем детальную информацию с табами
		if isTerminal() {
			// Интерактивная версия с табами
			program := ui.NewAgentDetailViewModel(agent)
			if err := program.Start(); err != nil {
				log.Fatal("Failed to start detail view", "error", err)
			}
		} else {
			// Простая версия для не-терминала
			result := ui.RenderAgentDetails(agent, ctx, container)
			fmt.Println(result)
		}
	},
}

// isTerminal проверяет, является ли терминал терминалом
func isTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

func init() {
	RootCMD.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&agentOutputFormat, "output", "o", "table", "Формат вывода (table, json)")
}
