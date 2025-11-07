package cmd

import "github.com/spf13/cobra"

var gistCmd = &cobra.Command{
	Use:   "gist",
	Short: "Manage GitHub gists",
}

func init() {
	rootCmd.AddCommand(gistCmd)
}
