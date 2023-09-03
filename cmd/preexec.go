package cmd

import (
	"flag"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gookit/slog"
	"github.com/spf13/cobra"
	tipView "github.com/sudonym-btc/zap/cmd/view/tip"
	"github.com/thoas/go-funk"
)

// preexecCmd represents the preexec command
var preexecCmd = &cobra.Command{
	Use:   "preexec",
	Short: "Runs automatically when you run a package manager command",
	Long:  `Do not run manually`,
	Run: func(cmd *cobra.Command, args []string) {
		model := parsePreexecArgs(args)
		if model != nil {
			tea.NewProgram(tipView.InitialTipModel(*model), tea.WithAltScreen()).Run()
		}
	},
}

func init() {
	rootCmd.AddCommand(preexecCmd)
}

func parsePreexecArgs(args []string) *tipView.TipModelParams {
	path, err := os.Getwd()
	if err != nil {
		return nil
	}

	if len(args) > 0 {
		args = strings.Split(args[0], " ")

		cmdArgs := []string{}

		for _, arg := range args {
			if !strings.Contains(arg, "--") {
				cmdArgs = append(cmdArgs, arg)
			}
		}
		flag.CommandLine.Parse(cmdArgs)
		cmdArgs = flag.CommandLine.Args()
		resolver := funk.Find(tipView.PackageCommandResolvers, func(cr tipView.CommandResolver) bool {
			return cr.Command == cmdArgs[0]
		})
		if resolver != nil {

			do, pkgs, err := resolver.(tipView.CommandResolver).ParseArgs(cmdArgs[1:])

			if do {
				var name string = ""

				if pkgs != nil && len(*pkgs) > 0 {
					name = (*pkgs)[0].Name
				}
				r := resolver.(tipView.CommandResolver)
				return &tipView.TipModelParams{
					CommandResolver: &r,
					Cwd:             &path,
					Name:            &name,
				}

			}
			if err != nil {
				slog.Warn("Failed parsing args", err, cmdArgs)
				fmt.Println(err)
			}
		}
	}

	return nil
}
