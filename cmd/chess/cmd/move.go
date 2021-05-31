package cmd

import (
	"log"

	"github.com/dumbogo/chess/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(moveCmd)

	moveCmd.Flags().StringVarP(&from, "from", "f", "", "Square location from")
	moveCmd.Flags().StringVarP(&to, "to", "t", "", "Square location to")
	moveCmd.MarkFlagRequired("from")
	moveCmd.MarkFlagRequired("to")
}

var (
	from string
	to   string
)

var moveCmd = &cobra.Command{
	Use:   "move",
	Short: "Move piece",
	Long:  "Move piece from square locations",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := client.InitConn()
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		defer conn.Close()
		client.Move(conn, from, to)
	},
}
