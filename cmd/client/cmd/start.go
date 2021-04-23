package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start game",
	Long:  "Start a new game",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: functionality to start new game
	},
}
