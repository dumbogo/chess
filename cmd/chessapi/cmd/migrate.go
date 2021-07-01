package cmd

import (
	"log"

	"github.com/dumbogo/chess/api"
	"github.com/dumbogo/chess/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().StringVarP(&configFile, "config", "c", "", "TOML configuration file to start API server")
	migrateCmd.MarkFlagRequired("config")
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "database migrations",
	Long:  "Run API database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			configuration *config.ServerConfig
			err           error
		)
		if configuration, err = config.LoadServerConfig(configFile); err != nil {
			panic(err)
		}
		_, err = api.InitDbConn(
			configuration.DBHost,
			configuration.DBPort,
			configuration.DBUser,
			configuration.DBPassword,
			configuration.DBName,
		)
		if err != nil {
			log.Fatalf("failed to connect databse: %v", err)
		}
		if err := api.Migrate(); err != nil {
			log.Fatalf("failed to run migrations: %v", err)
		}
		log.Printf("Successfully ran migrations")
	},
}
