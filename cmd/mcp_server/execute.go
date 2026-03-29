package mcp_server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

var (
	executeToolName string
	executeArgs     string
)

var executeCmd = &cobra.Command{
	Use:   "execute <server-id>",
	Short: "Выполнить инструмент MCP сервера",
	Long: `Выполняет указанный инструмент на MCP сервере с переданными параметрами.

Примеры:
  ai-agents-cli mcp-servers execute <server-id> --tool my-tool --args '{"key": "value"}'`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		serverID := args[0]

		if executeToolName == "" {
			log.Fatal("Tool name is required. Use --tool flag")
		}

		var params map[string]interface{}
		if executeArgs != "" {
			if err := json.Unmarshal([]byte(executeArgs), &params); err != nil {
				log.Fatal("Failed to parse args JSON", "error", err)
			}
		} else {
			params = make(map[string]interface{})
		}

		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			log.Fatal("Failed to get API client", "error", err)
		}

		result, err := apiClient.MCPServers.ExecuteTool(ctx, serverID, executeToolName, params)
		if err != nil {
			log.Fatal("Failed to execute tool", "error", err, "server_id", serverID, "tool", executeToolName)
		}

		successStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("2")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		fmt.Println(successStyle.Render(fmt.Sprintf("✅ Инструмент '%s' выполнен", executeToolName)))
		fmt.Println()

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			log.Fatal("Failed to encode result", "error", err)
		}
	},
}

func init() {
	RootCMD.AddCommand(executeCmd)

	executeCmd.Flags().StringVarP(&executeToolName, "tool", "t", "", "Название инструмента для выполнения")
	executeCmd.Flags().StringVarP(&executeArgs, "args", "a", "", "Параметры инструмента в формате JSON")
}
