package cmd

import (
	"log"
	"net"
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

	// Sensitive config
	dbUser     string // CHESS_API_DATABASE_USERNAME env
	dbPassword string // CHESS_API_DATABASE_PASSWORD env
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

	// TODO: Set ENVS as mandatory
	v.SetEnvPrefix("CHESS_API")
	v.AllowEmptyEnv(false) // This doesn't work as expected
	v.BindEnv("database_username")
	v.BindEnv("database_password")
	dbUser = v.GetString("database_username")
	dbPassword = v.GetString("database_password")
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
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}
