package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"master-alias.com/core"
	"master-alias.com/core/utils"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:   "remove [ID|NAME]",
	Short: "Remove an alias by ID|NAME",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		removeAlias(args[0])
	},
}

func removeAlias(index string) {
	alias := utils.FindById(core.FILENAME, index)

	if alias.Id == "" {
		alias = utils.FindByName(core.FILENAME, index)
	}

	utils.RemoveItem(core.FILENAME, alias.Id)
	fmt.Println(alias.Name)
	utils.RemoveAliasFromFile(alias.Name)

	c := color.New(color.FgGreen)
	c.Println("  ✅ Alias removed successfully")
}
