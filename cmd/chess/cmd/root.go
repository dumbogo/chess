package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chess",
	Short: "Chess game",
	Long:  "Chess multi-player game on terminal",
}

// Execute executes root command
func Execute() error {
	return rootCmd.Execute()
}
