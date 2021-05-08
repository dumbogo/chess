package cmd

import (
	"log"

	"github.com/dumbogo/chess/api"
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
		db, err := api.InitDbConn(dbHost, dbPort, dbUser, dbPassword, dbName)
		if err != nil {
			log.Fatalf("failed to connect databse: %v", err)
		}
		if err := api.Migrate(db); err != nil {
			log.Fatalf("failed to run migrations: %v", err)
		}
		log.Printf("Successfully ran migrations")
	},
}
