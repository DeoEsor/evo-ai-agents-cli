package mcp_server

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Генерация скрипта автодополнения для shell",
	Long: `Генерация скрипта автодополнения для указанной оболочки.

Примеры:
  # Bash
  ai-agents-cli mcp-servers completion bash > /etc/bash_completion.d/ai-agents-cli-mcp
  
  # Zsh
  ai-agents-cli mcp-servers completion zsh > "${fpath[1]}/_ai-agents-cli-mcp"`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			RootCMD.GenBashCompletion(os.Stdout)
		case "zsh":
			RootCMD.GenZshCompletion(os.Stdout)
		case "fish":
			RootCMD.GenFishCompletion(os.Stdout, true)
		case "powershell":
			RootCMD.GenPowerShellCompletion(os.Stdout)
		}
	},
}

func init() {
	RootCMD.AddCommand(completionCmd)
}
