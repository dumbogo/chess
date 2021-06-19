package config

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/viper"
)

var (
	// ClientConfig client configuration
	ClientConfig struct {
		// AuthToken token authenticated to make API calls
		AuthToken string // oauth2.*.token
		// ClientCertfile client certificate TLS location file
		ClientCertfile string // CLIENT_CERTFILE

		// ServerNameOverride is for testing only. If set to a non empty string,
		// it will override the virtual host name of authority (e.g. :authority header
		// field) in requests.
		ServerNameOverride string // SERVERNAME_OVERRIDE

		// APIServerURL URL API to make calls
		APIServerURL string // API_SERVER_URL

		// Game configuration current game
		Game *GameClientConfig
	}

	configName = "config"
	configType = "toml"
	// configPath loaded dinamically
	configPath = ".chess"
)

// GameClientConfig configuration current game
type GameClientConfig struct {
	// Name declared name by player initializer
	Name string
	// Color Color assigned to player
	Color string // game.color
	// Uuid game identifier
	UUID string // game.uuid
}

// InitClientConfig initializes ClientConfig
func InitClientConfig() {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	var (
		homeDir string
		err     error
	)
	if homeDir, err = os.UserHomeDir(); err != nil {
		panic(fmt.Errorf("home directory not found %v", err))
	}
	configPath := path.Join(homeDir, configPath)
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Config file not found, please add a config file on %s directory", configPath)
		} else {
			panic(err)
		}
	}
	ClientConfig.APIServerURL = viper.GetString("API_SERVER_URL")
	ClientConfig.ClientCertfile = viper.GetString("CLIENT_CERTFILE")
	ClientConfig.ServerNameOverride = viper.GetString("SERVERNAME_OVERRIDE")
	ClientConfig.AuthToken = viper.GetString("oauth2.github.token") // TODO: hardcoded to github, change it when implementing more providers
	ClientConfig.Game = &GameClientConfig{
		UUID:  viper.GetString("game.uuid"),
		Name:  viper.GetString("game.name"),
		Color: viper.GetString("game.color"),
	}
}

// UpdateGame update Game Configuration and persist it
func UpdateGame(gc *GameClientConfig) error {
	viper.Set("game.uuid", gc.UUID)
	viper.Set("game.name", gc.Name)
	viper.Set("game.color", gc.Color)
	if err := viper.WriteConfig(); err != nil {
		log.Fatalf("Unexpected Error: %s", err.Error())
		return err
	}
	return nil
}

// SetClientAuthToken sets auth token to configuration and persist
func SetClientAuthToken(token string) error {
	viper.Set("oauth2.github.token", token)
	if err := viper.WriteConfig(); err != nil {
		log.Fatalf("Unexpected Error: %s", err.Error())
		return err
	}
	return nil
}
