package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "master-alias",
	Short: "Manage custom aliases",
	Long:  "master-alias lets you create, run and list custom aliases with support for parameters.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
