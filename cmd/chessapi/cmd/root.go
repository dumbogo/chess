package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chessapi",
	Short: "Chess server API",
	Long:  "Chess server API",
}

// Execute executes root command
func Execute() error {
	return rootCmd.Execute()
}
