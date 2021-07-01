package cmd

import (
	"fmt"
	"log"

	"github.com/dumbogo/chess/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(defaultConfigCmd)
	configCmd.AddCommand(viewConfigCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Game configuration",
	Long:  "Set/Get game configuration",
}

var defaultConfigCmd = &cobra.Command{
	Use:   "default",
	Short: "Print default configuration",
	Long:  "Print default configuration with mandatory fields to play",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.NewClientConfiguration(config.WithDefaultBaseClientConfiguration())
		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}
		str, err := c.Marshal()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", str)
	},
}

var viewConfigCmd = &cobra.Command{
	Use:   "view",
	Short: "Show current configuration",
	Long:  "Print current configuration client chess game",
	Run: func(cmd *cobra.Command, args []string) {
		c, err := config.LoadClientConfiguration()
		if err != nil {
			log.Fatalf("Unexpected Error %s", err)
		}
		str, err := c.Marshal()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", str)
	},
}

// TODO: chess config credentials, sub command to configure credentias
// Must investigate about how to fix the problem of using CA credentials in the project
