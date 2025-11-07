package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"master-alias.com/core"
	"master-alias.com/core/structs"
	"master-alias.com/core/utils"

	"github.com/spf13/cobra"
)

var name string
var command string
var tag string
var description string

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&name, "name", "n", "", "Name of the alias")
	addCmd.Flags().StringVarP(&command, "command", "c", "", "Command to be executed")
	addCmd.Flags().StringVarP(&tag, "tag", "t", "", "Tag to be added to the alias")
	addCmd.Flags().StringVarP(&description, "description", "d", "", "Description of the alias")
}

var addCmd = &cobra.Command{
	Use:   "add [NAME] [COMMAND]",
	Short: "Add a new alias",
	//Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		addAlias(args)
	},
}

func addAlias(args []string) {
	var aliasExist = utils.FindByName("alias.json", name)

	c := color.New(color.FgRed)
	if aliasExist.Name != "" {
		c.Println(" ❌  An alias with this name already exists")
		return
	}

	if err := utils.WriteAliasToFile(name); err != nil {
		utils.Debug("error writing aliases file:", err)
		return
	}

	aliases, err := utils.ReadJSON("alias.json")
	if err != nil {
		// if not exist, start with empty slice
		if os.IsNotExist(err) {
			aliases = []structs.Alias{}
		} else {
			utils.Debug("error reading alias.json:", err)
			return
		}
	}

	aliases = append(aliases, structs.Alias{
		Id:          uuid.NewString(),
		Name:        name,
		Command:     command,
		Tag:         tag,
		Description: description})

	if err := utils.WriteJSON(core.FILENAME, aliases); err != nil {
		utils.Debug("error writing alias.json:", err)
		return
	}

	utils.Debug("alias added:", name, command)

	c = color.New(color.FgGreen)
	c.Println("  ✅ Alias created successfully")

}
