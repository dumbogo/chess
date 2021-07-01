package config

import (
	"log"
	"path/filepath"

	"github.com/dumbogo/chess/api"
	"github.com/spf13/viper"
)

var (
	configFileType = "toml"
)

// ServerConfig server configuration
type ServerConfig struct {
	ENV     string
	APIPort string // API.port

	// API config
	APIServerCert string // API.server_cert
	APIServerKey  string // API.server_key

	// Database config

	DBHost string // Database.host
	DBPort string // Database.port
	DBName string // Database.db_name

	// HTTP server
	HTTPServerScheme string // HTTP_server.Scheme
	HTTPServerHost   string // HTTP_server.Host
	HTTPServerPort   string // HTTP_server.Port

	// Sensitive config
	DBUser     string // CHESS_API_DATABASE_USERNAME env
	DBPassword string // CHESS_API_DATABASE_PASSWORD env

	// Github credentials
	GithubKey    string // CHESS_API_GITHUB_KEY
	GithubSecret string // CHESS_API_GITHUB_SECRET

	NATsURL string // NATS_URL

}

// LoadServerConfig ...
func LoadServerConfig(configFile string) (*ServerConfig, error) {
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
	c := &ServerConfig{}

	// Set default env to development if not set
	c.ENV = v.GetString("ENV")
	if !(c.ENV == api.EnvProduction || c.ENV == api.EnvTest || c.ENV == api.EnvDev) {
		c.ENV = api.EnvDev
	}

	c.APIPort = v.GetString("API.port")
	c.APIServerCert = v.GetString("API.server_cert")
	c.APIServerKey = v.GetString("API.server_key")

	c.DBHost = v.GetString("Database.host")
	c.DBPort = v.GetString("Database.port")
	c.DBName = v.GetString("Database.db_name")

	c.HTTPServerScheme = v.GetString("HTTP_server.Scheme")
	c.HTTPServerHost = v.GetString("HTTP_server.Host")
	c.HTTPServerPort = v.GetString("HTTP_server.Port")

	// TODO: Set ENVS as mandatory
	v.SetEnvPrefix("CHESS_API")
	v.AllowEmptyEnv(false) // This doesn't work as expected

	if err := v.BindEnv("DATABASE_USERNAME"); err != nil {
		log.Fatalf("Unexpected error %s", err.Error())
	}
	if !v.IsSet("DATABASE_USERNAME") {
		log.Fatalf("required env %s", "CHESS_API_DATABASE_USERNAME")
	}
	c.DBUser = v.GetString("database_username")

	if err := v.BindEnv("DATABASE_PASSWORD"); err != nil {
		log.Fatalf("Unexpected error %s", err.Error())
	}
	if !v.IsSet("DATABASE_PASSWORD") {
		log.Fatalf("required env %s", "CHESS_API_DATABASE_PASSWORD")
	}
	c.DBPassword = v.GetString("database_password")

	if err := v.BindEnv("GITHUB_KEY"); err != nil {
		log.Fatalf("Unexpected error %s", err.Error())
	}
	c.GithubKey = v.GetString("github_key")

	if err := v.BindEnv("GITHUB_SECRET"); err != nil {
		log.Fatalf("Unexpected error %s", err.Error())
	}
	c.GithubSecret = v.GetString("github_secret")

	if err := v.BindEnv("NATS_URL"); err != nil {
		log.Fatalf("Unexpected error %s", err.Error())
	}
	c.NATsURL = v.GetString("nats_url")
	return c, nil
}
