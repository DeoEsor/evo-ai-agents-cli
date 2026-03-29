package mcp_server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/cloud-ru/evo-ai-agents-cli/internal/di"
	"github.com/spf13/cobra"
)

var (
	toolsOutputFormat string
)

var toolsCmd = &cobra.Command{
	Use:   "tools <server-id>",
	Short: "Список инструментов MCP сервера",
	Long:  "Показывает доступные инструменты для указанного MCP сервера",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		serverID := args[0]

		container := di.GetContainer()
		apiClient, err := container.GetAPI()
		if err != nil {
			log.Fatal("Failed to get API client", "error", err)
		}

		tools, err := apiClient.MCPServers.GetTools(ctx, serverID)
		if err != nil {
			log.Fatal("Failed to get MCP server tools", "error", err, "server_id", serverID)
		}

		if toolsOutputFormat == "json" {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(tools); err != nil {
				log.Fatal("Failed to encode JSON", "error", err)
			}
			return
		}

		headerStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

		shortID := serverID
		if len(shortID) > 8 {
			shortID = shortID[:8] + "..."
		}

		fmt.Println(headerStyle.Render(fmt.Sprintf("🔧 Инструменты MCP сервера %s", shortID)))
		fmt.Println()

		if len(tools) == 0 {
			fmt.Println("🔍 Инструменты не найдены")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Название\tОписание")
		fmt.Fprintln(w, "--------\t--------")

		for _, tool := range tools {
			fmt.Fprintf(w, "%s\t%s\n", tool.Name, tool.Description)
		}

		w.Flush()
	},
}

func init() {
	RootCMD.AddCommand(toolsCmd)

	toolsCmd.Flags().StringVarP(&toolsOutputFormat, "output", "o", "table", "Формат вывода (table, json)")
}
