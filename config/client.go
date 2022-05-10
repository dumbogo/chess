package config

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/pelletier/go-toml"
	"github.com/spf13/viper"
)

const (
	configName = "config"
	configType = "toml"
	// configPath loaded dinamically
	configPath = ".chess"
)

// ClientConfiguration configuration set by the client to interact with the game
type ClientConfiguration struct {
	// APIServerURL URL API to make calls
	APIServerURL string // API_SERVER_URL

	// ClientCertfile client certificate TLS location file
	ClientCertfile string // CLIENT_CERTFILE

	// ServerNameOverride is for testing only. If set to a non empty string,
	// it will override the virtual host name of authority (e.g. :authority header
	// field) in requests.
	ServerNameOverride string // SERVERNAME_OVERRIDE

	// AuthToken token authenticated to make API calls
	AuthToken string // oauth2.*.token

	// Game configuration current game
	Game *gameClientConfig
}

type gameClientConfig struct {
	// Name declared name by player initializer
	Name string
	// Color Color assigned to player
	Color string // game.color
	// Uuid game identifier
	UUID string // game.uuid
}

// LoadClientConfiguration loads config info from file configuration client
func LoadClientConfiguration() (*ClientConfiguration, error) {
	config := &ClientConfiguration{}

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	var (
		homeDir string
		err     error
	)
	if homeDir, err = os.UserHomeDir(); err != nil {
		return nil, fmt.Errorf("home directory not found %v", err)
	}
	configPath := path.Join(homeDir, configPath)
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("Config file not found, please add a config file on %s directory", configPath)
		}
		return nil, (err)
	}
	config.APIServerURL = viper.GetString("API_SERVER_URL")
	config.ClientCertfile = viper.GetString("CLIENT_CERTFILE")
	config.ServerNameOverride = viper.GetString("SERVERNAME_OVERRIDE")
	config.AuthToken = viper.GetString("oauth2.github.token") // TODO: hardcoded to github, change it when implementing more providers
	config.Game = &gameClientConfig{
		UUID:  viper.GetString("game.uuid"),
		Name:  viper.GetString("game.name"),
		Color: viper.GetString("game.color"),
	}
	return config, nil
}

// NewClientConfiguration returns new Client configuration
func NewClientConfiguration(opts ...ClientConfigurationOption) (*ClientConfiguration, error) {
	c := &ClientConfiguration{}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// ClientConfigurationOption option to configure client
type ClientConfigurationOption func(*ClientConfiguration)

// WithDefaultBaseClientConfiguration option to configure default base configuration
func WithDefaultBaseClientConfiguration() ClientConfigurationOption {
	return func(c *ClientConfiguration) {
		c.APIServerURL = "localhost"
		c.ClientCertfile = "$HOME/.chess/certs/x509/client.crt"
		c.ServerNameOverride = "www.fabrikam.com"
	}
}

// UpdateGame update Game Configuration and persist it
func (cc *ClientConfiguration) UpdateGame(uuid string, name string, color string) error {
	cc.Game.UUID = uuid
	cc.Game.Name = name
	cc.Game.Color = color

	viper.Set("game.uuid", cc.Game.UUID)
	viper.Set("game.name", cc.Game.Name)
	viper.Set("game.color", cc.Game.Color)

	if err := viper.WriteConfig(); err != nil {
		log.Fatalf("Unexpected Error: %s", err.Error())
		return err
	}
	return nil
}

// SetAuthToken sets auth token to configuration and persist
func (cc *ClientConfiguration) SetAuthToken(token string) error {
	viper.Set("oauth2.github.token", token)
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	cc.AuthToken = token
	return nil
}

// Marshal returns text representation of configuration
func Marshal() ([]byte, error) {
	c := viper.AllSettings()
	return toml.Marshal(c)
}

// Marshal returns text representation instance ClientConfiguration
func (cc *ClientConfiguration) Marshal() ([]byte, error) {
	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType(configType)

	v.Set("API_SERVER_URL", cc.APIServerURL)
	v.Set("CLIENT_CERTFILE", cc.ClientCertfile)
	v.Set("SERVERNAME_OVERRIDE", cc.ServerNameOverride)
	v.Set("oauth2.github.token", cc.AuthToken)

	if cc.Game != nil {
		v.Set("game.uuid", cc.Game.UUID)
		v.Set("game.name", cc.Game.Name)
		v.Set("game.color", cc.Game.Color)
	}
	c := v.AllSettings()
	return toml.Marshal(c)
}
