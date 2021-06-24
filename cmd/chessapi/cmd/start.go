package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/url"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/dumbogo/chess/api"
	pb "github.com/dumbogo/chess/api"
	"github.com/dumbogo/chess/messagebroker"
)

var (
	configFile     string
	configFileType = "toml"

	// Env environent, production, test, development
	Env string // ENV
	// API config
	apiPort       string // API.port
	apiServerCert string // API.server_cert
	apiServerKey  string // API.server_key

	// Database config

	dbHost string // Database.host
	dbPort string // Database.port
	dbName string // Database.db_name

	// HTTP server
	httpServerScheme string // HTTP_server.Scheme
	httpServerHost   string // HTTP_server.Host
	httpServerPort   string // HTTP_server.Port

	// Sensitive config
	dbUser     string // CHESS_API_DATABASE_USERNAME env
	dbPassword string // CHESS_API_DATABASE_PASSWORD env

	// Github credentials
	githubKey    string // CHESS_API_GITHUB_KEY
	githubSecret string // CHESS_API_GITHUB_SECRET

	natsURL string // NATS_URL
)

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&configFile, "config", "c", "", "TOML configuration file to start API server")
	if err := startCmd.MarkFlagRequired("config"); err != nil {
		panic(err)
	}
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// Load config file
	v := viper.New()
	// use filepath
	v.SetConfigName(filepath.Base(configFile))
	v.SetConfigType(configFileType)
	v.AddConfigPath(filepath.Dir(configFile))
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if configFile != "" {
				log.Fatalf("could not load file %s, please provide a valid one", configFile)
			}
		} else {
			panic(err)
		}
	}
	// Set default env to development if not set
	Env = v.GetString("ENV")
	if !(Env == api.EnvProduction || Env == api.EnvTest || Env == api.EnvDev) {
		Env = api.EnvDev
	}

	apiPort = v.GetString("API.port")
	apiServerCert = v.GetString("API.server_cert")
	apiServerKey = v.GetString("API.server_key")

	dbHost = v.GetString("Database.host")
	dbPort = v.GetString("Database.port")
	dbName = v.GetString("Database.db_name")

	httpServerScheme = v.GetString("HTTP_server.Scheme")
	httpServerHost = v.GetString("HTTP_server.Host")
	httpServerPort = v.GetString("HTTP_server.Port")

	// TODO: Set ENVS as mandatory
	v.SetEnvPrefix("CHESS_API")
	v.AllowEmptyEnv(false) // This doesn't work as expected

	if err := v.BindEnv("DATABASE_USERNAME"); err != nil {
		log.Fatalf("Unexpected error %s", err.Error())
	}
	if !v.IsSet("DATABASE_USERNAME") {
		log.Fatalf("required env %s", "CHESS_API_DATABASE_USERNAME")
	}
	dbUser = v.GetString("database_username")

	if err := v.BindEnv("DATABASE_PASSWORD"); err != nil {
		log.Fatalf("Unexpected error %s", err.Error())
	}
	if !v.IsSet("DATABASE_PASSWORD") {
		log.Fatalf("required env %s", "CHESS_API_DATABASE_PASSWORD")
	}
	dbPassword = v.GetString("database_password")

	if err := v.BindEnv("GITHUB_KEY"); err != nil {
		log.Fatalf("Unexpected error %s", err.Error())
	}
	githubKey = v.GetString("github_key")

	if err := v.BindEnv("GITHUB_SECRET"); err != nil {
		log.Fatalf("Unexpected error %s", err.Error())
	}
	githubSecret = v.GetString("github_secret")

	if err := v.BindEnv("NATS_URL"); err != nil {
		log.Fatalf("Unexpected error %s", err.Error())
	}
	natsURL = v.GetString("nats_url")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start API server",
	Long:  "Start API server",
	Run: func(cmd *cobra.Command, args []string) {
		lis, err := net.Listen("tcp", apiPort)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		cert, err := tls.LoadX509KeyPair(apiServerCert, apiServerKey)
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
		db, err := api.InitDbConn(dbHost, dbPort, dbUser, dbPassword, dbName)
		if err != nil {
			log.Fatalf("failed to connect databse: %v", err)
		}

		// Load nats connection
		mb, err := messagebroker.New(messagebroker.Config{URL: natsURL})
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
					Scheme: httpServerScheme,
					Host:   fmt.Sprintf("%s%s", httpServerHost, httpServerPort),
				},
				githubKey,
				githubSecret,
				// Ensure your key is sufficiently random - i.e. use Go's
				// crypto/rand or securecookie.GenerateRandomKey(32) and persist the result.
				"somerandomtext",
				Env,
			)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}
			log.Printf("Listening HTTP server, on port %s\n", httpServerPort)
			if err := s.Listen(); err != nil {
				log.Fatalf("failed to serve: %v", err)
			}
		}()

		// Load grpc server
		log.Printf("Listening gRPC server, on port %s\n", apiPort)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}
