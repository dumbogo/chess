package cmd

import (
	"log"

	"github.com/dumbogo/chess/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(joinCmd)

	joinCmd.Flags().StringVarP(&uuid, "uuid", "I", "", "uuid game")
	joinCmd.MarkFlagRequired("uuid")
}

var (
	uuid string
)

var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join game",
	Long:  "Join a started game",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := client.InitConn()
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		defer conn.Close()
		client.JoinGame(conn, uuid)
	},
}
