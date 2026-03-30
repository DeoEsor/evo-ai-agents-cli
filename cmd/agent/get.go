package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/errors"
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
		errorHandler := errors.NewHandler()

		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			appErr := errorHandler.WrapAPIError(err, "API_CLIENT_ERROR", "Ошибка получения API клиента")
			appErr = appErr.WithSuggestions(agentAPISuggestions()...)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		agent, err := apiClient.Agents.Get(ctx, agentID)
		if err != nil {
			appErr := errorHandler.WrapAPIError(err, "AGENT_GET_FAILED", "Ошибка получения агента")
			appErr = appErr.WithSuggestions(
				"Проверьте правильность ID агента: "+agentID,
				"Убедитесь что агент существует: ai-agents-cli agents list",
			)
			fmt.Println(errorHandler.HandlePlain(appErr))
			os.Exit(1)
		}

		if agentOutputFormat == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(agent); err != nil {
				appErr := errorHandler.WrapFileSystemError(err, "JSON_ENCODE_ERROR", "Ошибка кодирования JSON")
				fmt.Println(errorHandler.HandlePlain(appErr))
				os.Exit(1)
			}
			return
		}

		if isTerminal() {
			program := ui.NewAgentDetailViewModel(agent)
			if err := program.Start(); err != nil {
				appErr := errorHandler.WrapUserError(err, "UI_ERROR", "Ошибка отображения")
				fmt.Println(errorHandler.HandlePlain(appErr))
				os.Exit(1)
			}
		} else {
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
