package client

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/dumbogo/chess/api"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	// APIServerURL URL API to make calls
	APIServerURL string
	// PublicKeyFile Local file location public key to authenticate
	PublicKeyFile string

	configName = "config"
	configType = "toml"
	configPath = "$HOME/.chess"
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
	PublicKeyFile = viper.GetString("PUBLIC_KEY_FILE")
}

// InitConn initializes client connection to GRPC server
func InitConn() (*grpc.ClientConn, error) {
	return grpc.Dial(APIServerURL, grpc.WithInsecure())
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

	fmt.Printf("Moved piece, response :%v\n", r)
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
	viper.WriteConfig()
}
