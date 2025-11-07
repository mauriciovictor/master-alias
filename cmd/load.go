package cmd

import (
	"fmt"
	"os"

	"master-alias.com/core"
	"master-alias.com/core/utils"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(loadCmd)
}

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "source commands",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		loadSourceAliases()
	},
}

func loadSourceAliases() {
	aliases, err := utils.ReadJSON(core.FILENAME)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No aliases found yet.")
			return
		}
		fmt.Println("Error reading alias.json:", err)
		return
	}

	if len(aliases) == 0 {
		fmt.Println("No aliases registered yet.")
		return
	}

	utils.CreateShellAliasFile()

	for _, a := range aliases {
		if err := utils.WriteAliasToFile(a.Name); err != nil {
		}
	}

	if err := utils.EnsureShellSource(); err != nil {
		fmt.Println("Error writing alias.json:", err)
		return
	}

	fmt.Println("Aliases atualizados! Execute o comando abaixo para carregar:")
	fmt.Println("# => source ~/.master-alias/master_aliases.sh")

}
