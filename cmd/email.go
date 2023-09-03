package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	connectView "github.com/sudonym-btc/zap/cmd/view/config"
)

// connectCmd represents the connect command
var emailCmd = &cobra.Command{
	Use:   "email",
	Short: "Configure email connectivity",
	Long:  `Connect your email account with zap, so that you can send lightning gifts if no lightning address available`,
	Run: func(cmd *cobra.Command, args []string) {
		tea.NewProgram(connectView.InitialEmailModel(), tea.WithAltScreen()).Run()
	},
}

func init() {
	rootCmd.AddCommand(emailCmd)
}
