package cmd

import (
	"log"
	"net"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	pb "github.com/dumbogo/chess/api"
)

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&port, "port", "p", ":8000", "running port")
	startCmd.MarkFlagRequired("port")
}

var (
	port string
)

var startCmd = &cobra.Command{
	Use:   "run",
	Short: "Start API server",
	Long:  "Start API server",
	Run: func(cmd *cobra.Command, args []string) {
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterChessServiceServer(s, &pb.Server{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}
