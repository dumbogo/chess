package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(moveCmd)
}

var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move piece",
	Long:  "Move piece from square locations",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: functionality to move piece
	},
}
