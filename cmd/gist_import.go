package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"master-alias.com/core/integrations/github"
)

var gistIdFlag string

func init() {
	gistImportCmd.Flags().StringVarP(&gistIdFlag, "gist_id", "i", "", "Id of the gist to import")
	gistCmd.AddCommand(gistImportCmd)
}

var gistImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import aliases from a github gist",
	//Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if gistIdFlag == "" {
			fmt.Println("Required flag --gist_id not set.")
			return
		}

		github.ImportGist(gistIdFlag)

	},
}
