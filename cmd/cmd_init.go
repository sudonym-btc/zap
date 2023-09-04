package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gookit/slog"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "cmdInit",
	Short: "Adds preexec to your shell",
	Long:  `This command will add preexec command to your shell, such that future package manager calls are registered and the tip command can run automatically.`,
	Run: func(cmd *cobra.Command, args []string) {
		shells := getShells()
		shells = append(shells, args...)
		slog.Debug("Initing cmd listener", shells)
		shells = funk.Uniq(shells).([]string)
		shells = funk.FilterString(shells, func(s string) bool {
			return s != ""
		})
		for _, shell := range shells {
			switch shell {
			case "-zsh":
			case "/bin/zsh":
				c := exec.Command("bash", "-c", `echo 'preexec () { 
	zap preexec $1
}' >> ~/.zshrc`)
				err := c.Run()
				if err != nil {
					panic(err)
				}
				fmt.Println("Added zap command listener to .zshrc")
				slog.Debug("Added zap command listener to .zshrc")

			case "/bin/bash":
				c := exec.Command("bash", "-c", `curl https://raw.githubusercontent.com/rcaloras/bash-preexec/master/bash-preexec.sh -o ~/.bash-preexec.sh && echo '[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh
			 preexec() { zap preexec $1; }' >> ~/.bashrc`)
				err := c.Run()
				if err != nil {
					panic(err)
				}
				fmt.Println("Added zap command listener to .bashrc")
				slog.Debug("Added zap command listener to .bashrc")

			default:
				fmt.Println("No shell detected for " + shell)
				slog.Debug("No shell detected for " + shell)
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func getShells() []string {
	return funk.Uniq([]string{os.Getenv("SHELL"), os.Getenv("0")}).([]string)
}
