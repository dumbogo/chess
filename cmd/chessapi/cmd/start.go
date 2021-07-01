package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/url"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/dumbogo/chess/api"
	pb "github.com/dumbogo/chess/api"
	"github.com/dumbogo/chess/config"
	"github.com/dumbogo/chess/messagebroker"
)

var (
	configFile string
)

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&configFile, "config", "c", "", "TOML configuration file to start API server")
	if err := startCmd.MarkFlagRequired("config"); err != nil {
		panic(err)
	}
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start API server",
	Long:  "Start API server",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			configuration *config.ServerConfig
			err           error
		)
		if configuration, err = config.LoadServerConfig(configFile); err != nil {
			panic(err)
		}
		lis, err := net.Listen("tcp", configuration.APIPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		cert, err := tls.LoadX509KeyPair(configuration.APIServerCert, configuration.APIServerKey)
		if err != nil {
			log.Fatalf("failed to load key pair: %s", err)
		}
		opts := []grpc.ServerOption{
			// The following grpc.ServerOption adds an interceptor for all unary
			// RPCs. To configure an interceptor for streaming RPCs, see:
			// https://godoc.org/google.golang.org/grpc#StreamInterceptor
			grpc.UnaryInterceptor(api.EnsureValidToken),

			// Enable TLS for all incoming connections.
			grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		}
		s := grpc.NewServer(opts...)
		db, err := api.InitDbConn(configuration.DBHost, configuration.DBPort, configuration.DBUser, configuration.DBPassword, configuration.DBName)
		if err != nil {
			log.Fatalf("failed to connect databse: %v", err)
		}

		// Load nats connection
		mb, err := messagebroker.New(messagebroker.Config{URL: configuration.NATsURL})
		if err != nil {
			log.Fatalf("Failed to initialize nats: %v", err)
		}
		pb.MessageBroker = mb
		pb.RegisterChessServiceServer(s, &pb.Server{
			Db: db,
		})

		// Load HTTP server
		go func() {
			s, err := api.NewHTTPServer(
				url.URL{
					Scheme: configuration.HTTPServerScheme,
					Host:   fmt.Sprintf("%s%s", configuration.HTTPServerHost, configuration.HTTPServerPort),
				},
				configuration.GithubKey,
				configuration.GithubSecret,
				// Ensure your key is sufficiently random - i.e. use Go's
				// crypto/rand or securecookie.GenerateRandomKey(32) and persist the result.
				"somerandomtext",
				configuration.ENV,
			)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			log.Printf("Listening HTTP server, on port %s\n", configuration.HTTPServerPort)
			if err := s.Listen(); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		}()

		// Load grpc server
		log.Printf("Listening gRPC server, on port %s\n", configuration.APIPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}
