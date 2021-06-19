package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(defaultConfigCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Game config",
	Long:  "Set/Get game configuration",
}

var defaultConfig string = `

`
var defaultConfigCmd = &cobra.Command{
	Use:   "default",
	Short: "Default configuration",
	Long:  "Set default configuration to play",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("give me the cash")
	},
}
