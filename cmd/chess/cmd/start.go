package cmd

import (
	"log"

	pb "github.com/dumbogo/chess/api"
	"github.com/dumbogo/chess/client"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().StringVarP(&name, "name", "n", "", "Game name")
	startCmd.Flags().StringVarP(&color, "color", "c", "white", "Color to chose")
	startCmd.MarkFlagRequired("name")
}

var (
	name  string
	color string
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start game",
	Long:  "Start a new game",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := client.InitConn()
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		defer conn.Close()

		inputColor := pb.Color_WHITE
		switch color {
		case "white":
		case "black":
			inputColor = pb.Color_BLACK
		default:
			log.Fatalf("Must define either \"white\" or \"black\" color")
		}
		client.StartGame(conn, name, inputColor)
	},
}
