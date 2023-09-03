package cmd

import (
	"fmt"
	"os"
	"os/exec"

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
		shells = funk.Uniq(shells).([]string)
		for _, shell := range shells {
			switch shell {
			case "/bin/zsh":
				c := exec.Command("bash", "-c", `echo 'preexec () { 
	zap preexec $1
}' >> ~/.zshrc`)
				err := c.Run()
				if err != nil {
					panic(err)
				}
			case "/bin/bash":
				c := exec.Command("bash", "-c", `curl https://raw.githubusercontent.com/rcaloras/bash-preexec/master/bash-preexec.sh -o ~/.bash-preexec.sh && echo '[[ -f ~/.bash-preexec.sh ]] && source ~/.bash-preexec.sh
			 preexec() { zap preexec $1; }' >> ~/.bashrc`)
				err := c.Run()
				if err != nil {
					panic(err)
				}
			default:
				fmt.Println("No shell detected")
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
