package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/dumbogo/chess/api"
	"github.com/dumbogo/chess/config"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

const defaultTimeOutContext = time.Second

var clientConfig *config.ClientConfiguration

// InitConn loads configuration and initializes connection to server
func InitConn() (*grpc.ClientConn, error) {
	var err error
	clientConfig, err = config.LoadClientConfiguration()
	if err != nil {
		return nil, err
	}
	// Set up the credentials for the connection.
	perRPC := oauth.NewOauthAccess(&oauth2.Token{
		AccessToken: clientConfig.AuthToken,
	})
	pathCertFile, err := relPathtoFilePath(clientConfig.ClientCertfile)
	if err != nil {
		return nil, err
	}
	creds, err := credentials.NewClientTLSFromFile(
		pathCertFile,
		clientConfig.ServerNameOverride,
	)
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
	return grpc.Dial(clientConfig.APIServerURL, opts...)
}

// StartGame creates a new Game
func StartGame(conn *grpc.ClientConn, name string, color pb.Color) {
	c := pb.NewChessServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeOutContext)
	defer cancel()
	r, err := c.StartGame(ctx, &pb.StartGameRequest{Name: name, Color: color})
	if err != nil {
		log.Fatalf("could not start game: %v", err)
	}
	fmt.Printf("UUID to connect: %s\n Please share this UUID to your fellow in order to play\n", r.GetUuid())
	if err := clientConfig.UpdateGame(r.GetUuid(), name, color.String()); err != nil {
		panic(err)
	}
}

// Move call move piece server and print movement
func Move(conn *grpc.ClientConn, fromSquare, toSquare string) {
	c := pb.NewChessServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeOutContext)
	defer cancel()
	color := pb.Color_value[(clientConfig.Game.Color)]
	r, err := c.Move(ctx, &pb.MoveRequest{
		Color:      pb.Color(color),
		Uuid:       clientConfig.Game.UUID,
		FromSquare: fromSquare,
		ToSquare:   toSquare,
	})
	if err != nil {
		log.Fatalf("could not move piece: %v", err)
	}
	fmt.Printf("Board: \n:%s\n", r.GetBoard())
}

// JoinGame calls join game server
func JoinGame(conn *grpc.ClientConn, uuid string) {
	c := pb.NewChessServiceClient(conn)
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeOutContext)
	defer cancel()
	r, err := c.JoinGame(ctx, &pb.JoinGameRequest{Uuid: uuid})
	if err != nil {
		log.Fatalf("could not join game: %v", err)
	}
	fmt.Printf("Joined game uuid %s, name: %s color assigned: %s\n", r.GetUuid(), r.GetName(), r.GetColor())
	if err := clientConfig.UpdateGame(r.GetUuid(), r.GetName(), r.GetColor().String()); err != nil {
		panic(err)
	}
}

// Watch watches a live server game by uuid, if not provided, uses the configured by client.
// Outputs movements to STDOUT
func Watch(conn *grpc.ClientConn, uuid string) {
	c := pb.NewChessServiceClient(conn)
	if uuid == "" {
		uuid = clientConfig.Game.UUID
	}
	if uuid == "" {
		panic(errors.New("No current game, please either provide a game uuid or create/join one"))
	}
	stream, err := c.Watch(context.Background(), &pb.WatchRequest{Uuid: uuid})
	if err != nil {
		log.Fatalf("could not watch game: %v", err)
	}
	fmt.Printf("Watching game...")
	for {
		watchResponse, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Goodbye!")
			break
		}
		if err != nil {
			log.Fatalf("%v.Watch(_) = _, %v", c, err)
		}
		fmt.Printf("Turn: %s\n", watchResponse.GetTurn())
		fmt.Printf("Status: %s\n", watchResponse.GetStatus())
		fmt.Println(watchResponse.GetBoard())
	}
}

func relPathtoFilePath(path string) (string, error) {
	if !strings.Contains(path, "$HOME") {
		return path, nil
	}
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(path, "$HOME", userHomeDir), nil
}
