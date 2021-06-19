package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "github.com/dumbogo/chess/api"
	"github.com/dumbogo/chess/config"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

const defaultTimeOutContext = time.Second

var clientConfig = &config.ClientConfig

func init() {
	config.InitClientConfig()
}

// InitConn initializes client connection to GRPC server
func InitConn() (*grpc.ClientConn, error) {
	// Set up the credentials for the connection.
	perRPC := oauth.NewOauthAccess(&oauth2.Token{
		AccessToken: clientConfig.AuthToken,
	})
	creds, err := credentials.NewClientTLSFromFile(
		clientConfig.ClientCertfile,
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
	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeOutContext)
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

// JoinGame calls to gprc API JoinGame
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
	storeGame(r.GetUuid(), r.GetName(), r.GetColor())
}

// Watch watches a live game, outputs movements to STDOUT
// If not uuid is provided, set the stored in configuration
func Watch(conn *grpc.ClientConn, uuid string) {
	c := pb.NewChessServiceClient(conn)
	if uuid == "" {
		fmt.Printf("clientCOnfig: %+v\n", clientConfig)
		uuid = clientConfig.Game.UUID
	}
	// Contact the server and print out its response.
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

func storeGame(uuid string, name string, color pb.Color) {
	if err := config.UpdateGame(&config.GameClientConfig{
		Name:  name,
		UUID:  uuid,
		Color: color.String(),
	}); err != nil {
		panic(err)
	}
}

// RegisterGithubToken records token on persisted config
func RegisterGithubToken(token string) {
	if err := config.SetClientAuthToken(token); err != nil {
		panic(err)
	}
}
