package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"master-alias.com/core"
	"master-alias.com/core/utils"

	"github.com/spf13/cobra"
)

var tagFlag string

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().StringVarP(&tagFlag, "tag", "t", "", "Filter aliases by tag")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all registered aliases",
	Long:  `Displays all aliases stored in the alias.json file with their commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		listAliases()
	},
}

func listAliases() {
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

	// Create an aligned table in the terminal
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(writer, "ID\tNAME\tCOMMAND\tTAG\tDESCRIPTION")
	fmt.Fprintln(writer, "------------------------------------\t----\t-------\t---\t-----------")

	for _, a := range aliases {
		if tagFlag != "" && a.Tag == tagFlag {
			fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n", a.Id, a.Name, a.Command, a.Tag, a.Description)
		} else if tagFlag == "" {
			fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s\n", a.Id, a.Name, a.Command, a.Tag, a.Description)
		}
	}

	writer.Flush()
}
