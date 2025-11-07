package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"master-alias.com/core"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Set up master-alias for the first time",
	Run: func(cmd *cobra.Command, args []string) {
		var shell string
		var githubToken string
		var reconfigure string

		if core.ConfigExists() {
			fmt.Println("✅ A configuration already exists. reconfigure?")

			shellPrompt := &survey.Select{
				Message: "Choose the shell:",
				Options: []string{"yes", "no"},
				Default: "no",
			}

			survey.AskOne(shellPrompt, &reconfigure)

			if reconfigure == "no" {
				return
			}

		}

		// Question 1 - shell selection
		shellPrompt := &survey.Select{
			Message: "Choose the shell:",
			Options: []string{"zsh", "bash"},
			Default: "zsh",
		}
		survey.AskOne(shellPrompt, &shell)

		//Github Token
		githubTokenPrompt := &survey.Password{
			Message: "Github Token:",
		}

		survey.AskOne(githubTokenPrompt, &githubToken)

		cfg := &core.Config{
			Shell:       shell,
			GithubToken: githubToken,
		}

		if err := core.SaveConfig(cfg); err != nil {
			fmt.Println("❌ Error saving configuration:", err)
			return
		}

		fmt.Println("🎉 Configuration completed successfully!")
		core.PrintConfig()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
