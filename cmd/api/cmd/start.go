package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"github.com/dumbogo/chess/api"
	pb "github.com/dumbogo/chess/api"
)

var (
	configFile     string
	configFileType = "toml"

	// API config
	apiPort string // API.port

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
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&configFile, "config", "c", "", "TOML configuration file to start API server")
	startCmd.MarkFlagRequired("config")

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
			log.Fatalf("could not load file %s, please provide a valid one", configFile)
		} else {
			panic(err)
		}
	}
	apiPort = v.GetString("API.port")

	dbHost = v.GetString("Database.host")
	dbPort = v.GetString("Database.port")
	dbName = v.GetString("Database.db_name")

	httpServerScheme = v.GetString("HTTP_server.Scheme")
	httpServerHost = v.GetString("HTTP_server.Host")
	httpServerPort = v.GetString("HTTP_server.Port")

	// TODO: Set ENVS as mandatory
	v.SetEnvPrefix("CHESS_API")
	v.AllowEmptyEnv(false) // This doesn't work as expected
	v.BindEnv("database_username")
	v.BindEnv("database_password")
	v.BindEnv("github_key")
	v.BindEnv("github_secret")
	dbUser = v.GetString("database_username")
	dbPassword = v.GetString("database_password")
	githubKey = v.GetString("github_key")
	githubSecret = v.GetString("github_secret")
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
		s := grpc.NewServer()
		db, err := api.InitDbConn(dbHost, dbPort, dbUser, dbPassword, dbName)
		if err != nil {
			log.Fatalf("failed to connect databse: %v", err)
		}
		pb.RegisterChessServiceServer(s, &pb.Server{
			Db: db,
		})

		// Load HTTP server
		go func() {
			api.InitHTTPRouter(url.URL{Scheme: httpServerScheme, Host: httpServerHost}, githubKey, githubSecret)
			log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost%s", httpServerPort), api.HTTPRouter))
		}()

		// Load grpc server
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}
