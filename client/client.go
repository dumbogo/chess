package client

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/dumbogo/chess/api"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

var (
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

	configName = "config"
	configType = "toml"
	configPath = "~/.chess"
)

func init() {
	initConfig()
}

func initConfig() {
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Config file not found, please add a config file on $HOME/.chess directory")
		} else {
			panic(err)
		}
	}
	APIServerURL = viper.GetString("API_SERVER_URL")
	ClientCertfile = viper.GetString("CLIENT_CERTFILE")
	ServerNameOverride = viper.GetString("SERVERNAME_OVERRIDE")
	AuthToken = viper.GetString("oauth2.github.token") // TODO: hardcoded to github, change it when implementing more providers
}

// InitConn initializes client connection to GRPC server
func InitConn() (*grpc.ClientConn, error) {
	// Set up the credentials for the connection.
	perRPC := oauth.NewOauthAccess(&oauth2.Token{
		AccessToken: AuthToken,
	})

	creds, err := credentials.NewClientTLSFromFile(ClientCertfile, ServerNameOverride)
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}
	opts := []grpc.DialOption{
		// In addition to the following grpc.DialOption, callers may also use
		// the grpc.CallOption grpc.PerRPCCredentials with the RPC invocation
		// itself.
		// See: https://godoc.org/google.golang.org/grpc#PerRPCCredentials
		grpc.WithPerRPCCredentials(perRPC),
		grpc.WithTransportCredentials(creds),
	}
	// opts = append(opts, grpc.WithBlock())
	return grpc.Dial(APIServerURL, opts...)
}

// StartGame creates a new Game
func StartGame(conn *grpc.ClientConn, name string, color pb.Color) {
	c := pb.NewChessServiceClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.StartGame(ctx, &pb.StartGameRequest{
		Name:  name,
		Color: color,
	})
	if err != nil {
		log.Fatalf("could not start game: %v", err)
	}
	fmt.Printf("UUID to connect: %s\n Please share this UUID to your fellow in order to play\n", r.GetUuid())
	storeGame(r.GetUuid(), name, color)
}

// Move calls to gprc API move
func Move(conn *grpc.ClientConn, fromSquare, toSquare string) {
	c := pb.NewChessServiceClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	color := pb.Color_value[(viper.GetString("game.color"))]
	r, err := c.Move(ctx, &pb.MoveRequest{
		Color:      pb.Color(color),
		Uuid:       viper.GetString("game.uuid"),
		FromSquare: fromSquare,
		ToSquare:   toSquare,
	})
	if err != nil {
		log.Fatalf("could not move piece: %v", err)
	}

	fmt.Printf("Board: \n:%s\n", r.GetBoard())
}

// JoinGame calls to gprc API JoinGame
func JoinGame(conn *grpc.ClientConn, uuid string) {
	c := pb.NewChessServiceClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.JoinGame(ctx, &pb.JoinGameRequest{Uuid: uuid})
	if err != nil {
		log.Fatalf("could not join game: %v", err)
	}

	fmt.Printf("Joined game uuid %s, name: %s color assigned: %s\n", r.GetUuid(), r.GetName(), r.GetColor())
	storeGame(r.GetUuid(), r.GetName(), r.GetColor())
}

func storeGame(uuid string, name string, color pb.Color) {
	viper.Set("game.uuid", uuid)
	viper.Set("game.name", name)
	viper.Set("game.color", color)
	if err := viper.WriteConfig(); err != nil {
		panic(err)
	}
}

// RegisterGithubToken records token on persisted config
func RegisterGithubToken(token string) {
	viper.Set("oauth2.github.token", token)
	if err := viper.WriteConfig(); err != nil {
		panic(err)
	}
}
