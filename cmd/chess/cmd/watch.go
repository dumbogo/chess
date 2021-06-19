package cmd

import (
	"log"

	"github.com/dumbogo/chess/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(watchCmd)
	watchCmd.Flags().StringVarP(&uuid, "uuid", "I", "", "uuid game")
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "watch game",
	Long:  "watch movements live with a terminal session",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := client.InitConn()
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		defer conn.Close()
		client.Watch(conn, uuid)
	},
}
