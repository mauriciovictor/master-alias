package cmd

import (
	"github.com/spf13/cobra"
	"master-alias.com/core/integrations/github"
)

func init() {
	gistCmd.AddCommand(gistExportCmd)
}

var gistExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Create a new or update github gist",
	//Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		github.ExportGist()
	},
}
