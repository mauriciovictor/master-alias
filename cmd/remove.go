package cmd

import (
	"github.com/fatih/color"
	"master-alias.com/core"
	"master-alias.com/core/utils"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove [ID]",
	Short: "Remove an alias by ID",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		removeAlias(args[0])
	},
}

func removeAlias(id string) {
	alias := utils.FindById(core.FILENAME, id)

	utils.RemoveItem(core.FILENAME, id)
	utils.RemoveAliasFromFile(alias.Name)

	c := color.New(color.FgGreen)
	c.Println("  ✅ Alias removed successfully")
}
