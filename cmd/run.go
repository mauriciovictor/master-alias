package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"master-alias.com/core"
	"master-alias.com/core/utils"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run [NAME]",
	Short: "Execute a command by alias name",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runAlias(args)
	},
}

func runAlias(args []string) {
	if len(args) < 1 {
		utils.Debug("usage: master-alias run [NAME]")
		return
	}

	aliasName := args[0]
	params := args[1:]
	utils.Debug("Running alias:", aliasName)
	utils.Debug("Params:", params)

	alias := utils.FindByName("alias.json", aliasName)
	if alias.Name == "" {
		utils.Debug("Alias not found.")
		return
	}

	command := alias.Command

	// Replace $1, $2, $@ manually
	for i, p := range params {
		placeholder := fmt.Sprintf("$%d", i+1)
		command = strings.ReplaceAll(command, placeholder, p)
	}
	command = strings.ReplaceAll(command, "$@", strings.Join(params, " "))

	utils.Debug("Executing:", command)

	config, err := core.LoadConfig()
	if err != nil || config.Shell == "" {
		fmt.Println("Error identifying terminal")
		return
	}

	cmd := exec.Command(config.Shell, "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		utils.Debug("Error executing command:", err)
	}
}
