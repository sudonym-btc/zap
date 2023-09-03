package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	configView "github.com/sudonym-btc/zap/cmd/view/config"
	tipView "github.com/sudonym-btc/zap/cmd/view/tip"
	wallet "github.com/sudonym-btc/zap/service"
	"github.com/sudonym-btc/zap/service/config"
	"github.com/sudonym-btc/zap/service/email"
	"github.com/sudonym-btc/zap/service/tip"
	"github.com/thoas/go-funk"
)

var tipCmd = &cobra.Command{
	Use:   "tip",
	Short: "Tip current work dir dependencies, specific package, or github repos/users",
	Long: `tip
		tip npm angular
		tip npm angular@latest
		tip github sudonym-btc --amount 1000
		tip paco@walletofsatoshi.com --amount 1000 --comment "Thanks for the great work!"
		`,
	Run: func(cmd *cobra.Command, args []string) {
		tipModel := parseArgs(cmd, args)
		conf, _ := config.LoadConfig()

		fmt.Println("Checking wallet connection...")
		if conf.WalletConnect == "" {
			tea.NewProgram(configView.InitialConnectModel(), tea.WithAltScreen()).Run()
		}
		_, err2 := wallet.Parse_and_connect(conf.WalletConnect)
		if err2 != nil {
			tea.NewProgram(configView.InitialConnectModel(), tea.WithAltScreen()).Run()
		}

		fmt.Println("Checking email connection...")
		sendEmails, _ := cmd.Flags().GetBool("sendEmails")
		if sendEmails {
			_, err2 := email.Connect(conf.Smtp)
			if err2 != nil {
				tea.NewProgram(configView.InitialEmailModel(), tea.WithAltScreen()).Run()
			}
		}
		tea.NewProgram(tipModel, tea.WithAltScreen()).Run()
	},
}

func parseArgs(cmd *cobra.Command, args []string) tipView.TipModel {
	path, err := os.Getwd()
	if err != nil {
	}
	var commandResolver *tipView.CommandResolver
	var name string
	var packageVersion string
	if len(args) > 0 {
		// Address
		if len(tip.ExtractEmails(args[0])) > 0 {
			name = tip.ExtractEmails(args[0])[0]
			cr := funk.Find(tipView.CommandResolvers, func(cr tipView.CommandResolver) bool {
				return cr.Command == "address"
			}).(tipView.CommandResolver)
			commandResolver = &cr
		} else {
			// Specific command resolver
			var cr *tipView.CommandResolver
			index := funk.IndexOf(tipView.CommandResolvers, func(cr tipView.CommandResolver) bool {
				return cr.Command == args[0]
			})
			cr = &tipView.CommandResolvers[index]
			if cr != nil {
				commandResolver = cr
				if commandResolver.Command == "github-user" {
					if len(args) > 1 {
						name = args[1]
					} else {
						panic("Please provide a github username")
					}
				}
				if len(args) > 1 {
					name = args[1]
				}
			}
		}
	}
	amount, _ := cmd.Flags().GetInt("amount")
	comment, _ := cmd.Flags().GetString("comment")
	manual, _ := cmd.Flags().GetBool("manual")
	sendEmails, _ := cmd.Flags().GetBool("sendEmails")

	return tipView.InitialTipModel(tipView.TipModelParams{
		CommandResolver: commandResolver,
		Name:            &name,
		Version:         &packageVersion,
		Cwd:             &path,
		Amount:          &amount,
		Comment:         &comment,
		Manual:          &manual,
		SendEmails:      &sendEmails,
	})

}

func init() {
	rootCmd.AddCommand(tipCmd)
	tipCmd.PersistentFlags().Int("amount", 0, "Amount to tip in satoshis")
	tipCmd.PersistentFlags().String("comment", "", "Comment to attach to zaps")
	tipCmd.PersistentFlags().Bool("manual", false, "Should we step through each maintainer manually")
	tipCmd.PersistentFlags().Bool("sendEmails", false, "Should we attempt to send emails to maintainers if no lightning addresses are found")

}
