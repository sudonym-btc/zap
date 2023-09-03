package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	connectView "github.com/sudonym-btc/zap/cmd/view/config"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect lightning wallet",
	Long:  `Connect Zap to your lightning wallet in order to be able to send payments.`,
	Run: func(cmd *cobra.Command, args []string) {
		tea.NewProgram(connectView.InitialConnectModel(), tea.WithAltScreen()).Run()

	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
